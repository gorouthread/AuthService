package auth_repository_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
)

func (r *AuthRepositoryPostgres) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
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
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, userID)

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
