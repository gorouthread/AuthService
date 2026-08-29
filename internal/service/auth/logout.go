package auth_service

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
)

func (s *AuthService) Logout(ctx context.Context, session domain.Session) error {
	sessionFromDB, err := s.authRepositoryPostgres.GetSessionByRefreshToken(ctx, session)
	if err != nil {
		return fmt.Errorf("get session from database: %w: %w", err, core_errors.ErrNotFound)
	}

	sessionFromDB.IsRevoked = true

	_, err = s.authRepositoryPostgres.UpdateSession(ctx, sessionFromDB)
	if err != nil {
		return fmt.Errorf("update session from database: %w", err)
	}

	return nil
}
