package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Refresh(ctx context.Context, session domain.Session) (domain.Session, error) {
	sessionFromDB, err := s.authRepositoryPostgres.GetSessionByRefreshToken(ctx, session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get session from database: %w: %w", err, core_errors.ErrNotFound)
	}

	if sessionFromDB.RefreshExpiresAt.Before(time.Now()) {
		return domain.Session{}, fmt.Errorf("token is expires: %w", core_errors.ErrUnauthorized)
	}

	if sessionFromDB.IsRevoked == true {
		return domain.Session{}, fmt.Errorf("token is revoked: %w", core_errors.ErrUnauthorized)
	}

	sessionFromDB.IsRevoked = true
	sessionFromDB, err = s.authRepositoryPostgres.UpdateSession(ctx, sessionFromDB)
	if err != nil {
		return domain.Session{}, fmt.Errorf("update session from database: %w", err)
	}

	user, err := s.authRepositoryPostgres.GetUserByID(ctx, sessionFromDB.UserID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get user from database: %w: %w", err, core_errors.ErrNotFound)
	}

	accessToken, accessClaims, err := s.jwtManager.CreateToken(user.ID, user.Role, 15*time.Minute)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create access jwt: %w", err)
	}

	remaining := time.Until(sessionFromDB.RefreshExpiresAt)
	refreshToken, refreshClaims, err := s.jwtManager.CreateToken(user.ID, user.Role, remaining)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create refresh jwt: %w", err)
	}

	hashRefreshToken := utils.HashToken(refreshToken)

	sessionFromDB = domain.NewSessionUninitialized(
		user.ID,
		hashRefreshToken,
		refreshClaims.IssuedAt.Time,
		refreshClaims.ExpiresAt.Time,
	)

	session, err = s.authRepositoryPostgres.CreateSession(ctx, sessionFromDB)
	if err != nil {
		return domain.Session{}, fmt.Errorf("save session to database: %w", err)
	}

	session.SetAccessParam(
		accessToken,
		accessClaims.IssuedAt.Time,
		accessClaims.ExpiresAt.Time,
	)

	session.RefreshToken = refreshToken
	return session, nil
}
