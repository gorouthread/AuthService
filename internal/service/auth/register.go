package auth_service

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Register(ctx context.Context, user domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validate user domain: %w", err)
	}

	ph, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = ph

	_, err = s.authRepositoryPostgres.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("save user to database: %w", err)
	}

	return nil
}
