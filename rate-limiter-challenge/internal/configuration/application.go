package configuration

import (
	"reflect"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Application struct {
	MaxRequestIP    int    `envconfig:"MAX_REQUESTS_IP" required:"true"`
	MaxRequestToken int    `envconfig:"MAX_REQUESTS_TOKEN" required:"true"`
	BlockTimeIP     int    `envconfig:"BLOCK_TIME_IP" required:"true"`
	BlockTimeToken  int    `envconfig:"BLOCK_TIME_TOKEN" required:"true"`
	RedisHost       string `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
	RedisPort       string `envconfig:"REDIS_PORT" required:"true" default:"6379"`
}

func (a *Application) Load() error {
	err := godotenv.Load()
	if err != nil {
		return errors.Wrap(err, "could not load local environment file")
	}
	err = envconfig.Process("", a)
	if err != nil {
		return errors.Wrap(err, "could not process values from environment file")
	}
	return nil
}

func (a *Application) GetValue(key string) string {
	r := reflect.ValueOf(a)
	f := reflect.Indirect(r).FieldByName(key)
	return f.String()
}

func NewApplication() *Application {
	return &Application{}
}
