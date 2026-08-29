package core_service_jwt

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretKey string `envconfig:"SECRET_KEY"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("JWT", &cfg); err != nil {
		return Config{}, fmt.Errorf("process read config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}

	return cfg
}
