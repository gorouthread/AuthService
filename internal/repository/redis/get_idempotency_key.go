package auth_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *AuthRepositoryRedis) GetIdempotencyKey(ctx context.Context, idemKey string) (*IdempotencyResponse, error) {
	key := fmt.Sprintf(idempotencyPrefix, idemKey)

	val, err := r.pool.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var resp IdempotencyResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal idempotency resp: %w", err)
	}

	return &resp, nil
}
