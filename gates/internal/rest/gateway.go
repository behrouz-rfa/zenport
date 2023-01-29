package rest

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"zenport/gates/gatespb"

	"github.com/go-chi/chi/v5"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/times"

	gateway := runtime.NewServeMux()
	err := gatespb.RegisterGatesServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the GRPC gates
	mux.Mount(apiRoot, gateway)

	return nil
}
