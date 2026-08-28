package auth_repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *AuthRepositoryRedis) GetAccessToken(ctx context.Context, token string) (*AuthData, error) {
	key := fmt.Sprintf(tokenPrefix, token)

	val, err := r.pool.Get(ctx, key).Result()
	if err != nil {
		return nil, err // Возвращает redis.Nil, если не найдено
	}

	var data AuthData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("unmarshal auth data: %w", err)
	}

	return &data, nil
}
