package rb

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Session composes an amqp.Connection with an amqp.Channel
type Session struct {
	*amqp.Connection
	*amqp.Channel
}

func (s Session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

func Redial(ctx context.Context, url string, exchange string) chan chan Session {
	Sessions := make(chan chan Session)
	go func() {
		sess := make(chan Session)
		defer close(Sessions)
		for {
			select {
			case Sessions <- sess:
			case <-ctx.Done():
				log.Println("Shuting down Session")
				return
			}
			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, url)
			}
			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}
			if err := ch.ExchangeDeclare(exchange, "fanout", false, true, false, false, nil); err != nil {
				log.Fatalf("Connot declare fanout exchange : %v", err)
			}

			select {
			case sess <- Session{conn, ch}:
			case <-ctx.Done():
				log.Println("shut down new Session")
				return
			}
		}
	}()

	return Sessions
}
