package auth_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (r *AuthRepositoryRedis) SaveIdempotencyKey(ctx context.Context, idemKey string, resp IdempotencyResponse, ttl time.Duration) error {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal idempotency resp: %w", err)
	}

	key := fmt.Sprintf(idempotencyPrefix, idemKey)
	return r.pool.Set(ctx, key, jsonData, ttl).Err()
}
