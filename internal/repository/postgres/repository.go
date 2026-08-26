package auth_repository_postgres

import (
	core_repository_postgres "github.com/romreign/AuthService/internal/core/repository/postgres"
	core_repository_redis "github.com/romreign/AuthService/internal/core/repository/redis"
)

type AuthRepositoryPsql struct {
	pgPool core_repository_postgres.Pool
	rdPool core_repository_redis.Pool
}

func NewAuthRepositoryPsql(
	postgresPool core_repository_postgres.Pool,
	redisPool core_repository_redis.Pool,
) AuthRepositoryPsql {
	return AuthRepositoryPsql{
		pgPool: postgresPool,
		rdPool: redisPool,
	}
}
