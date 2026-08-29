package auth_repository_postgres

import (
	"context"
	"fmt"

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

func (a *AuthRepositoryPostgres) Info() {
	row := a.pool.QueryRow(context.Background(), `
	SELECT
		current_database(),
		current_schema(),
		current_user,
		inet_server_addr(),
		inet_server_port()
	`)

	var (
		dbName     string
		schema     string
		user       string
		serverIP   *string
		serverPort int
	)

	err := row.Scan(
		&dbName,
		&schema,
		&user,
		&serverIP,
		&serverPort,
	)
	if err != nil {
		fmt.Println("check postgres connection: %w", err)
	}

	fmt.Println("postgres connection info",
		"database", dbName,
		"schema", schema,
		"user", user,
		"server_ip", serverIP,
		"server_port", serverPort,
	)
}
