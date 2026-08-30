package auth_repository_postgres

import (
	core_repository_postgres "github.com/romreign/AuthService/internal/core/repository/postgres"
)

type AuthRepositoryPostgres struct {
	pool core_repository_postgres.Pool
}

func NewAuthRepositoryPostgres(
	pool core_repository_postgres.Pool,
) AuthRepositoryPostgres {
	return AuthRepositoryPostgres{
		pool: pool,
	}
}
