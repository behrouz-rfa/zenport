package notifications

import (
	"context"
	"zenport/internal/am"
	"zenport/internal/ddd"
	"zenport/internal/jetstream"
	"zenport/internal/rbqm"
	"zenport/internal/registry"
	"zenport/internal/system"
	"zenport/notifications/internal/postgres"
	"zenport/ntps/ntpspb"

	"zenport/notifications/internal/handlers"

	"zenport/internal/monolith"
	"zenport/notifications/internal/application"

	"zenport/notifications/internal/logging"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	return Root(ctx, mono)
}
func Root(ctx context.Context, mono system.Service) (err error) {

	// setup Driven adapters
	reg := registry.New()
	if err := ntpspb.Registrations(reg); err != nil {
		return err
	}
	//work with nats or rabbitq
	var eventStream am.EventStream
	if mono.Config().RABBITMQC.IsEnable {
		eventStream = am.NewEventStream(reg, rbqm.NewStream(mono.RBSession()))
	} else {
		eventStream = am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	}
	// setup application
	ntpCacheRepo := postgres.NewNtpCacheRepository("notifications.ntp_cache", mono.DB())
	var app application.App
	app = application.New(ntpCacheRepo)
	app = logging.LogApplicationAccess(app, mono.Logger())

	ntpHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewNtpHandlers(ntpCacheRepo),
		"Notification", mono.Logger(),
	)

	// setup Driver adapters
	if err := handlers.RegisterNtpHandlers(ntpHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
