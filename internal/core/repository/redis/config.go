package core_repository_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `env:"HOST" default:"localhost"`
	Port     string `env:"PORT" default:"6379"`
	Password string `env:"PASSWORD" default:""`
	DB       int    `env:"DB" default:"0"`

	PoolSize     int `env:"POOL_SIZE" default:"10"`
	MinIdleConns int `env:"MIN_IDLE_CONNS" default:"5"`
	MaxRetries   int `env:"MAX_RETRIES" default:"3"`

	DialTimeout  time.Duration `env:"DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"3s"`
	PoolTimeout  time.Duration `env:"POOL_TIMEOUT" default:"4s"`

	Timeout time.Duration `env:"TIMEOUT" default:"5s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("REDIS", &cfg); err != nil {
		return Config{}, fmt.Errorf("process read config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get redis pool config: %w", err)
		panic(err)
	}

	return cfg
}
