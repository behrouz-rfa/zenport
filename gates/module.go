package gates

import (
	"context"
	"zenport/gates/internal/application"
	"zenport/gates/internal/grpc"
	"zenport/gates/internal/logging"
	postgres "zenport/gates/internal/postgress"
	"zenport/gates/internal/rest"
	"zenport/internal/monolith"
	"zenport/internal/system"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, mono system.Service) (err error) {
	//
	timeRepo := postgres.NewTimeRepository("ntps.time", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	ntps := grpc.NewNtpRepository(conn)
	var app application.App
	app = application.NewApplication(timeRepo, ntps)
	app = logging.LogApplicationAccess(app, mono.Logger())

	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	return nil
}
