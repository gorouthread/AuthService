package auth_service

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
)

func (s *AuthService) Register(ctx context.Context, user domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validate user domain: %w", err)
	}

	_, err := s.authRepositoryPostgres.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("save user to database: %w", err)
	}

	return nil
}
