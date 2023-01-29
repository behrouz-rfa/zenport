package monolith

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"zenport/internal/config"
	"zenport/internal/rb"
	"zenport/internal/waiter"
)

type Monolith interface {
	Config() config.AppConfig
	DB() *sql.DB
	Logger() zerolog.Logger
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	QM() *amqp.Connection
	CH() *amqp.Channel
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	RBSession() chan chan rb.Session
}

type Module interface {
	Startup(context.Context, Monolith) error
}
