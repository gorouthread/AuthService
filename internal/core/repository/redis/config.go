package core_repository_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT"`
	Password string `envconfig:"PASSWORD"`
	DB       int    `envconfig:"DB"`

	PoolSize     int `envconfig:"POOL_SIZE"`
	MinIdleConns int `envconfig:"MIN_IDLE_CONNS"`
	MaxRetries   int `envconfig:"MAX_RETRIES"`

	DialTimeout  time.Duration `envconfig:"DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT"`
	PoolTimeout  time.Duration `envconfig:"POOL_TIMEOUT"`

	Timeout time.Duration `envconfig:"TIMEOUT"`
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
