package application

import (
	"context"
	"zenport/ntps/internal/domain"
	"zenport/ntps/ntpspb"

	"zenport/internal/am"
	"zenport/internal/ddd"
)

type IntegrationEventHandlers[T ddd.AggregateEvent] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*IntegrationEventHandlers[ddd.AggregateEvent])(nil)

func NewIntegrationEventHandlers(publisher am.MessagePublisher[ddd.Event]) *IntegrationEventHandlers[ddd.AggregateEvent] {
	return &IntegrationEventHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func (h IntegrationEventHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.TimeCreatedEvent:
		return h.onTimeCreated(ctx, event)
	}
	return nil
}

func (h IntegrationEventHandlers[T]) onTimeCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.TimeCreated)
	return h.publisher.Publish(ctx, ntpspb.NtpsAggregateChannel,
		ddd.NewEvent(ntpspb.TimeCreatedEvent, &ntpspb.TimeCreated{
			Id:   event.AggregateID(),
			Time: payload.Time,
		}),
	)
}
