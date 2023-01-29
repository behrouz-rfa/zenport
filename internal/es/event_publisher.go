package es

import (
	"context"

	"zenport/internal/ddd"
)

type EventPublisher struct {
	AggregateNtp
	publisher ddd.EventPublisher[ddd.AggregateEvent]
}

var _ AggregateNtp = (*EventPublisher)(nil)

func NewEventPublisher(publisher ddd.EventPublisher[ddd.AggregateEvent]) AggregateNtpMiddleware {
	eventPublisher := EventPublisher{
		publisher: publisher,
	}

	return func(store AggregateNtp) AggregateNtp {
		eventPublisher.AggregateNtp = store
		return eventPublisher
	}
}

func (p EventPublisher) Save(ctx context.Context, aggregate EventSourcedAggregate) error {
	if err := p.AggregateNtp.Save(ctx, aggregate); err != nil {
		return err
	}
	return p.publisher.Publish(ctx, aggregate.Events()...)
}
