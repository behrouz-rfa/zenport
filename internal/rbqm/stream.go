package rbqm

import (
	"context"
	"crypto/sha1"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
	"zenport/internal/rb"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"zenport/internal/am"
)

const maxRetries = 5

type Stream struct {
	session chan chan rb.Session
}

func NewStream(session chan chan rb.Session) *Stream {
	return &Stream{session}
}

type message []byte

var _ am.MessageStream[am.RawMessage, am.RawMessage] = (*Stream)(nil)

// exchange binds the publishers to the subscribers
const exchange = "pubsub"

func (s *Stream) Publish(ctx context.Context, topicName string, rawMsg am.RawMessage) (err error) {
	var data []byte

	data, err = proto.Marshal(&StreamMessageRb{
		Id:   rawMsg.ID(),
		Name: rawMsg.MessageName(),
		Data: rawMsg.Data(),
	})
	msg := make(chan message)
	defer close(msg)
	go func() {
		publish(s.session, msg)
	}()

	msg <- data

	return nil
}

// identity returns the same host/process unique string for the lifetime of
// this process so that subscriber reconnections reuse the same queue name.
func identity() string {
	hostname, err := os.Hostname()
	h := sha1.New()
	fmt.Fprint(h, hostname)
	fmt.Fprint(h, err)
	fmt.Fprint(h, os.Getpid())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler[am.RawMessage], options ...am.SubscriberOption) error {

	msg := make(chan message)

	go func() {
		subscribe(s.session, msg)
	}()

	go func() {
		for ms := range msg {
			m := &StreamMessageRb{}
			err := proto.Unmarshal(ms, m)
			if err != nil {
				// TODO Nak? ... logging?
				return
			}

			msg := &rawMessage{
				id:    m.GetId(),
				name:  m.GetName(),
				data:  m.GetData(),
				acked: false,
			}

			errc := make(chan error)

			go func() {
				errc <- handler.HandleMessage(context.Background(), msg)
			}()

			select {
			case err = <-errc:
				return

			}

		}
	}()

	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler[am.RawMessage]) func(*nats.Msg) {
	return func(natsMsg *nats.Msg) {
		var err error

		m := &StreamMessageRb{}
		err = proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			// TODO Nak? ... logging?
			return
		}

		msg := &rawMessage{
			id:       m.GetId(),
			name:     m.GetName(),
			data:     m.GetData(),
			acked:    false,
			ackFn:    func() error { return natsMsg.Ack() },
			nackFn:   func() error { return natsMsg.Nak() },
			extendFn: func() error { return natsMsg.InProgress() },
			killFn:   func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				// TODO logging?
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					// TODO logging?
				}
				return
			}
			if nakErr := msg.NAck(); nakErr != nil {
				// TODO logging?
			}
		case <-wCtx.Done():
			// TODO logging?
			return
		}
	}
}

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func publish(sessions chan chan rb.Session, messages <-chan message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan message, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Printf("publishing...")

	Publish:
		for {
			var body message
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
				err := pub.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					pub.Close()
					break Publish
				}

			case body, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

// subscribe consumes deliveries from an exclusive queue from a fanout exchange and sends to the application specific messages chan.
func subscribe(sessions chan chan rb.Session, messages chan<- message) {
	queue := identity()

	for session := range sessions {
		sub := <-session

		if _, err := sub.QueueDeclare(queue, false, true, true, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", queue, err)
			return
		}

		routingKey := "application specific routing key for fancy topologies"
		if err := sub.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %q, %v", exchange, err)
			return
		}

		deliveries, err := sub.Consume(queue, "", false, true, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", queue, err)
			return
		}

		log.Printf("subscribed...")

		for msg := range deliveries {
			messages <- msg.Body
			sub.Ack(msg.DeliveryTag, false)
		}
	}
}
