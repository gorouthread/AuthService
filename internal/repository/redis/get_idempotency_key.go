package auth_repository_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/romreign/AuthService/internal/core/domain"
)

func (r *AuthRepositoryRedis) GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyData, error) {
	val, err := r.pool.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("get idempotency key from redis: %w", err)
	}

	var resp IdempotencyResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal idempotency resp: %w", err)
	}

	data := &domain.IdempotencyData{
		Method:     resp.Method,
		URL:        resp.URL,
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
	}

	return data, nil
}
