package system

import (
	"context"
	"database/sql"
	amqp "github.com/rabbitmq/amqp091-go"
	"zenport/internal/rb"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"zenport/internal/config"
	"zenport/internal/waiter"
)

type Service interface {
	Config() config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	QM() *amqp.Connection
	RBSession() chan chan rb.Session
	CH() *amqp.Channel
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() zerolog.Logger
}

type Module interface {
	Startup(context.Context, Service) error
}
