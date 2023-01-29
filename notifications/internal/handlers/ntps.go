package handlers

import (
	"context"
	"zenport/internal/am"
	"zenport/internal/ddd"
	"zenport/ntps/ntpspb"
)

func RegisterNtpHandlers(ntpHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return ntpHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(ntpspb.NtpsAggregateChannel, evtMsgHandler, am.MessageFilter{
		ntpspb.TimeCreatedEvent,
	}, am.GroupName("notification-ntp"))
}
