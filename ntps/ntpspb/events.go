package ntpspb

import (
	"zenport/internal/registry"
	"zenport/internal/registry/serdes"
)

const (
	NtpsAggregateChannel = "zenports.ntps.events.Ntp"

	TimeCreatedEvent = "ntpsapi.TimeCreated"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Time events
	if err := serde.Register(&TimeCreated{}); err != nil {
		return err
	}

	return nil
}
func (*TimeCreated) Key() string { return TimeCreatedEvent }
