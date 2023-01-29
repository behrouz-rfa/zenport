package application

import (
	"context"
	"zenport/internal/ddd"
	"zenport/ntps/ntpspb"
)

type CustomerHandlers[T ddd.Event] struct {
	cache NtpCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*CustomerHandlers[ddd.Event])(nil)

func NewCustomerHandlers(cache NtpCacheRepository) CustomerHandlers[ddd.Event] {
	return CustomerHandlers[ddd.Event]{
		cache: cache,
	}
}

func (h CustomerHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case ntpspb.TimeCreatedEvent:
		return h.onCustomerRegistered(ctx, event)
	}

	return nil
}

func (h CustomerHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := event.Payload().(*ntpspb.TimeCreated)
	return h.cache.ShowRequest(ctx, payload.GetId(), payload.GetTime())
}
