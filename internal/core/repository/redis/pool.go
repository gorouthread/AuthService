package core_repository_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Pool interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd

	Pipeline() redis.Pipeliner

	Close() error

	Client() *redis.Client

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	client    *redis.Client
	opTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, cfg Config) (*ConnectionPool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &ConnectionPool{
		client:    client,
		opTimeout: cfg.Timeout,
	}, nil
}

func (p *ConnectionPool) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) *redis.StatusCmd {
	return p.client.Set(ctx, key, value, expiration)
}

func (p *ConnectionPool) SetNX(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) *redis.BoolCmd {
	return p.client.SetNX(ctx, key, value, expiration)
}

func (p *ConnectionPool) Get(ctx context.Context, key string) *redis.StringCmd {
	return p.client.Get(ctx, key)
}

func (p *ConnectionPool) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return p.client.Del(ctx, keys...)
}

func (p *ConnectionPool) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return p.client.Exists(ctx, keys...)
}

func (p *ConnectionPool) Expire(
	ctx context.Context,
	key string,
	expiration time.Duration,
) *redis.BoolCmd {
	return p.client.Expire(ctx, key, expiration)
}

func (p *ConnectionPool) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return p.client.TTL(ctx, key)
}

func (p *ConnectionPool) Pipeline() redis.Pipeliner {
	return p.client.Pipeline()
}

func (p *ConnectionPool) Close() error {
	return p.client.Close()
}

func (p *ConnectionPool) Client() *redis.Client {
	return p.client
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
