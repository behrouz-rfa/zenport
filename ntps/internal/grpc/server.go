package grpc

import (
	"context"
	"github.com/google/uuid"
	"zenport/ntps/internal/application"
	"zenport/ntps/internal/application/commands"
	"zenport/ntps/ntpspb"

	"google.golang.org/grpc"
)

type server struct {
	app application.App
	ntpspb.UnimplementedTimeServiceServer
}

func (s server) GetTime(ctx context.Context, request *ntpspb.GetTimeRequest) (*ntpspb.GetTimeResponse, error) {
	t := request.Time
	time, err := s.app.CreateTime(ctx, commands.CreateTime{ID: uuid.New().String(), Time: t})
	if err != nil {
		return nil, err
	}
	return &ntpspb.GetTimeResponse{Time: time.Time}, nil
}

var _ ntpspb.TimeServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	ntpspb.RegisterTimeServiceServer(registrar, server{app: app})
	return nil
}
