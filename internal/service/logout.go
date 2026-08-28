package auth_service

import (
	"context"
	"fmt"

	"github.com/romreign/AuthService/internal/core/domain"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Logout(ctx context.Context, session domain.Session) error {
	hashToken, err := utils.HashString(session.RefreshToken)
	if err != nil {
		return fmt.Errorf("hash refresh JWT: %w", err)
	}

	session.RefreshToken = hashToken

	sessionFromDB, err := s.authRepositoryPostgres.GetSessionByRefreshToken(ctx, session)
	if err != nil {
		return fmt.Errorf("get session from database: %w: %w", err, core_transport_http_response.ErrNotFound)
	}

	_, err = s.authRepositoryPostgres.UpdateSession(ctx, sessionFromDB)
	if err != nil {
		return fmt.Errorf("update session from database: %w", err)
	}

	return nil
}
