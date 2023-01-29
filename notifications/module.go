package notifications

import (
	"context"
	"zenport/internal/am"
	"zenport/internal/ddd"
	"zenport/internal/jetstream"
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
	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	//eventStream := am.NewEventStream(reg, rbqm.NewStream(mono.RBSession()))

	// setup application
	customers := postgres.NewNtpCacheRepository("notifications.customers_cache", mono.DB())
	var app application.App
	app = application.New(customers)
	app = logging.LogApplicationAccess(app, mono.Logger())

	customerHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewCustomerHandlers(customers),
		"Customer", mono.Logger(),
	)

	// setup Driver adapters
	if err := handlers.RegisterNtpHandlers(customerHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
