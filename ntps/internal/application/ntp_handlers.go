package application

import (
	"context"
	"zenport/internal/ddd"
	"zenport/ntps/internal/domain"
)

type NtpHandlers[T ddd.AggregateEvent] struct {
	mall domain.NtpRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*NtpHandlers[ddd.AggregateEvent])(nil)

func NewMallHandlers(mall domain.NtpRepository) *NtpHandlers[ddd.AggregateEvent] {
	return &NtpHandlers[ddd.AggregateEvent]{
		mall: mall,
	}
}
func (h NtpHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.TimeCreatedEvent:
		return h.onTimeCreated(ctx, event)
	}
	return nil
}

func (h NtpHandlers[T]) onTimeCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.TimeCreated)
	return h.mall.AddTime(ctx, event.AggregateID(), payload.Time)
}
