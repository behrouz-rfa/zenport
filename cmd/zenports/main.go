package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net/http"
	"os"
	"zenport/gates"
	"zenport/internal/config"
	"zenport/internal/logger"
	"zenport/internal/monolith"
	"zenport/internal/rb"
	"zenport/internal/rpc"
	"zenport/internal/waiter"
	"zenport/internal/web"
	"zenport/notifications"
	"zenport/ntps"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {

	var cfg config.AppConfig
	// parse config/env/...
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	//create app
	m := app{cfg: cfg}

	// init infrastructure...
	m.db, err = sql.Open("pgx", cfg.PG.Conn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.db)

	// init nats & jetstream
	m.nc, err = nats.Connect(cfg.Nats.URL)
	if err != nil {

		return err
	}
	defer m.nc.Close()
	m.js, err = initJetStream(cfg.Nats, m.nc)
	if err != nil {
		return err
	}

	//init rabbitqm session
	m.rb_session = rb.Redial(context.Background(), cfg.RB.URL, cfg.RB.Exchange)

	//set logger
	m.logger = initLogger(cfg)
	//set rpc config
	m.rpc = initRpc(cfg.Rpc)
	//set  chi mux
	m.mux = initMux(cfg.Web)
	//create waiter : waiter for catch signal for interrupt disconnect rpc rabitq nats
	m.waiter = waiter.New(waiter.CatchSignals())

	// init modules
	m.modules = []monolith.Module{
		&gates.Module{},
		&ntps.Module{},
		&notifications.Module{},
	}

	//start all moduls as monolith
	if err = m.startupModules(); err != nil {
		return err
	}

	// Mount general web resources
	m.mux.Mount("/", http.FileServer(http.FS(web.WebUI)))

	fmt.Println("started zenport application")
	defer fmt.Println("stopped zenport application")

	m.waiter.Add(
		m.waitForWeb,
		m.waitForRPC,
		m.waitForStream,
	)

	return m.waiter.Wait()
}

func initRpc(_ rpc.RpcConfig) *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	return server
}

func initMux(_ web.WebConfig) *chi.Mux {
	return chi.NewMux()
}

func initJetStream(cfg config.NatsConfig, nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", cfg.Stream)},
	})

	return js, err
}
func initLogger(cfg config.AppConfig) zerolog.Logger {
	return logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    logger.Level(cfg.LogLevel),
	})
}
