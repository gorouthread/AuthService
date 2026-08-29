package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env         string `envconfig:"ENV"`
	ServiceName string `envconfig:"SERVICE_NAME"`
	Version     string `envconfig:"VERSION"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("APP", &cfg); err != nil {
		return Config{}, fmt.Errorf("process read config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return cfg
}
