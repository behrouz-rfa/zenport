package application

import (
	"context"
	"zenport/internal/ddd"
	"zenport/ntps/ntpspb"
)

type NtpHandlers[T ddd.Event] struct {
	cache NtpCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*NtpHandlers[ddd.Event])(nil)

func NewNtpHandlers(cache NtpCacheRepository) NtpHandlers[ddd.Event] {
	return NtpHandlers[ddd.Event]{
		cache: cache,
	}
}

func (h NtpHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case ntpspb.TimeCreatedEvent:
		return h.onTimeCreated(ctx, event)
	}

	return nil
}

func (h NtpHandlers[T]) onTimeCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*ntpspb.TimeCreated)
	return h.cache.ShowRequest(ctx, payload.GetId(), payload.GetTime())
}
