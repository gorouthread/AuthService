package auth_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (r *AuthRepositoryRedis) SaveAccessToken(ctx context.Context, token string, data AuthData, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal auth data: %w", err)
	}

	key := fmt.Sprintf(tokenPrefix, token)
	return r.pool.Set(ctx, key, jsonData, ttl).Err()
}
