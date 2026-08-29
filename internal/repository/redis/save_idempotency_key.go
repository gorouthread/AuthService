package auth_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
)

func (r *AuthRepositoryRedis) SaveIdempotencyKey(ctx context.Context, key string, data *domain.IdempotencyData, ttl time.Duration) error {
	req := IdempotencyResponse{
		Method:     data.Method,
		URL:        data.URL,
		StatusCode: data.StatusCode,
		Body:       data.Body,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal idempotency resp: %w", err)
	}

	return r.pool.SetNX(ctx, key, jsonData, ttl).Err()
}
