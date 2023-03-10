package ntps

import (
	"context"
	"zenport/internal/am"
	"zenport/internal/ddd"
	"zenport/internal/es"
	"zenport/internal/jetstream"
	pg "zenport/internal/postgres"
	"zenport/internal/rbqm"
	"zenport/internal/registry"
	"zenport/internal/registry/serdes"
	"zenport/internal/system"
	"zenport/ntps/internal/domain"
	"zenport/ntps/internal/grpc"
	"zenport/ntps/internal/handlers"
	"zenport/ntps/internal/postgres"
	"zenport/ntps/ntpspb"

	"zenport/internal/monolith"
	"zenport/ntps/internal/application"
	"zenport/ntps/internal/logging"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {

	return Root(ctx, mono)
}
func Root(ctx context.Context, mono system.Service) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = registrations(reg); err != nil {
		return err
	}
	if err = ntpspb.Registrations(reg); err != nil {
		return
	}
	var eventStream am.EventStream
	if mono.Config().RABBITMQC.IsEnable {
		eventStream = am.NewEventStream(reg, rbqm.NewStream(mono.RBSession()))
	} else {
		eventStream = am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	}
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateNtp := es.AggregateNtpWithMiddleware(
		pg.NewEventStore("ntps.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
	)
	ntps := es.NewAggregateRepository[*domain.Time](domain.NtpAggregate, reg, aggregateNtp)
	ntp := postgres.NewNtpRepository("ntps.ntps", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(ntps),
		mono.Logger(),
	)
	mallHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewMallHandlers(ntp),
		"Ntp", mono.Logger(),
	)
	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewIntegrationEventHandlers(eventStream),
		"IntegrationEvents", mono.Logger(),
	)

	// setup Driver adapters
	handlers.RegisterNtpHandlers(mallHandlers, domainDispatcher)
	handlers.RegisterIntegrationEventHandlers(integrationEventHandlers, domainDispatcher)

	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	return nil
}

// for decoding and encoding event for
// sending message to message brokers
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Time
	if err = serde.Register(domain.Time{}, func(v any) error {
		store := v.(*domain.Time)
		store.Aggregate = es.NewAggregate("", domain.NtpAggregate)
		return nil
	}); err != nil {
		return
	}

	// Time events
	if err = serde.Register(domain.TimeCreated{}); err != nil {
		return
	}
	return
}
