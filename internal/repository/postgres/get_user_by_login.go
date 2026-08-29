package auth_repository_postgres

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
)

func (r *AuthRepositoryPostgres) GetUserByLogin(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			id,
			login,
			password_hash,
			role,
			created_at,
			updated_at
		FROM auth_msginx.users
		WHERE login = $1;
	`

	row := r.pool.QueryRow(ctx, query, user.Login)

	var userModel User
	if err := row.Scan(
		&userModel.ID,
		&userModel.Login,
		&userModel.PasswordHash,
		&userModel.Role,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	); err != nil {
		return domain.User{}, fmt.Errorf("failed to scan user row: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Login,
		userModel.PasswordHash,
		userModel.Role,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return userDomain, nil
}
