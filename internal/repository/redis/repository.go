package auth_repository_redis

import core_repository_redis "github.com/romreign/AuthService/internal/core/repository/redis"

type AuthRepositoryRedis struct {
	pool core_repository_redis.Pool
}

func NewAuthRepositoryRedis(
	pool core_repository_redis.Pool,
) *AuthRepositoryRedis {
	return &AuthRepositoryRedis{
		pool: pool,
	}
}
