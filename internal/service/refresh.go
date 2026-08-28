package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Refresh(ctx context.Context, session domain.Session) (domain.Session, error) {
	hashToken, err := utils.HashString(session.RefreshToken)
	if err != nil {
		return domain.Session{}, fmt.Errorf("hash refresh JWT: %w", err)
	}

	session.RefreshToken = hashToken

	sessionFromDB, err := s.authRepositoryPostgres.GetSessionByRefreshToken(ctx, session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get session from database: %w: %w", err, core_transport_http_response.ErrNotFound)
	}

	if sessionFromDB.RefreshExpiresAt.Before(time.Now()) {
		return domain.Session{}, fmt.Errorf("token is expires: %w", core_transport_http_response.ErrInvalidArgument)
	}

	if sessionFromDB.IsRevoked == true {
		return domain.Session{}, fmt.Errorf("token is revoked: %w", core_transport_http_response.ErrInvalidArgument)
	}

	user, err := s.authRepositoryPostgres.GetUserByID(ctx, sessionFromDB.UserID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get user from database: %w: %w", err, core_transport_http_response.ErrNotFound)
	}

	accessToken, accessClaims, err := s.jwtManager.CreateToken(user.ID, user.Role, 15*time.Minute)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create access jwt: %w", err)
	}

	sessionFromDB.SetAccessParam(
		accessToken,
		accessClaims.IssuedAt.Time,
		accessClaims.ExpiresAt.Time,
	)

	return sessionFromDB, nil
}
