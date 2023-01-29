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
	RbqConfig struct {
		URL      string `required:"true"`
		Channel  string `default:"zenports"`
		Exchange string `default:"pubsub"`
	}

	RABBITMQConfig struct {
		IsEnable bool `required:"true"`
	}

	AppConfig struct {
		Environment     string
		Rpc             rpc.RpcConfig
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		Nats            NatsConfig
		PG              PGConfig
		RB              RbqConfig
		RABBITMQC       RABBITMQConfig
		Web             web.WebConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg AppConfig, err error) {

	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}
	err = envconfig.Process("", &cfg)

	return
}
