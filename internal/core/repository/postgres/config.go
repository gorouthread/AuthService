package core_repository_postgres

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `env:"HOST"`
	Port     string        `env:"PORT"`
	User     string        `env:"USER"`
	Password string        `env:"PASSWORD"`
	Database string        `env:"DATABASE"`
	Timeout  time.Duration `env:"TIMEOUT"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return Config{}, fmt.Errorf("process read config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get postgres pool config: %w", err)
		panic(err)
	}

	return cfg
}
