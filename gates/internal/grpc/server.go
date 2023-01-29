package grpc

import (
	"context"
	"google.golang.org/grpc"
	"zenport/gates/gatespb"
	"zenport/gates/internal/application"
)

type server struct {
	app application.App
	gatespb.UnimplementedGatesServiceServer
}

func (s server) GetTime(ctx context.Context, request *gatespb.GetTimeRequest) (*gatespb.GetTimeResponse, error) {
	t := request.Ask
	time, err := s.app.GetTime(ctx, application.TimeRequest{Ask: t})
	if err != nil {
		return nil, err
	}
	return &gatespb.GetTimeResponse{Time: time.Time}, nil
}

var _ gatespb.GatesServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	gatespb.RegisterGatesServiceServer(registrar, server{app: app})
	return nil
}
