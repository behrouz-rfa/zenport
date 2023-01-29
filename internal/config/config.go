package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
	"zenport/internal/rpc"

	"github.com/stackus/dotenv"
	"os"
	"zenport/internal/web"
)

type (
	PGConfig struct {
		Conn string `required:"true"`
	}
	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"zenports"`
	}
	//		URL      string `default:"amqp://guest:guest@localhost:5675/"`
	RBQMConfig struct {
		URL      string `default:"amqp://rabbitmq:rabbitmq@rabbit1:5672/"`
		Channel  string `default:"zenports"`
		Exchange string `default:"pubsub"`
	}

	AppConfig struct {
		Environment     string
		Rpc             rpc.RpcConfig
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		Nats            NatsConfig
		RB              RBQMConfig
		PG              PGConfig
		Web             web.WebConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg AppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIROMENT"))); err != nil {
		return
	}
	err = envconfig.Process("", &cfg)

	return
}
