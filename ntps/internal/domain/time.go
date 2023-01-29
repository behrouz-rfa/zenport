package domain

import (
	"github.com/stackus/errors"
	"time"
	"zenport/internal/ddd"

	"zenport/internal/es"
)

type Time struct {
	es.Aggregate
	Time string
}

const NtpAggregate = "npts.Time"

var (
	ErrStoreNameIsBlank = errors.Wrap(errors.ErrBadRequest, "it's wrong question ask : What time is it?")
)
var _ interface {
	es.EventApplier
} = (*Time)(nil)

func NewNtp(id string) *Time {
	return &Time{
		Aggregate: es.NewAggregate(id, NtpAggregate),
	}
}
func CreateTime(id string, input string) (*Time, error) {
	if input != "What time is it?" {
		return nil, ErrStoreNameIsBlank
	}
	ntp := NewNtp(id)
	ntp.AddEvent(TimeCreatedEvent, &TimeCreated{
		Time: time.Now().Format(time.RFC3339),
	})
	return ntp, nil
}

// Key implements registry.Registerable
func (Time) Key() string { return NtpAggregate }

// ApplyEvent implements es.EventApplier
func (s *Time) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *TimeCreated:
		s.Time = payload.Time

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}
