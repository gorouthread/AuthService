package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Login(ctx context.Context, user domain.User) (domain.Session, error) {
	if err := user.Validate(); err != nil {
		return domain.Session{}, fmt.Errorf("validate user domain: %w", err)
	}

	userFromDb, err := s.authRepositoryPostgres.GetUserByLogin(ctx, user)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get user from database: %w: %w", err, core_transport_http_response.ErrNotFound)
	}

	accessToken, accessClaims, err := s.jwtManager.CreateToken(userFromDb.ID, user.Role, 15*time.Minute)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create access jwt: %w", err)
	}

	refreshToken, refreshClaims, err := s.jwtManager.CreateToken(userFromDb.ID, user.Role, 24*time.Hour)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create refresh jwt: %w", err)
	}

	hashRefreshToken, err := utils.HashString(refreshToken)
	if err != nil {
		return domain.Session{}, fmt.Errorf("hashing refresh JWT: %w", err)
	}

	session := domain.NewSessionUninitialized(
		userFromDb.ID,
		hashRefreshToken,
		refreshClaims.IssuedAt.Time,
		refreshClaims.ExpiresAt.Time,
	)

	session, err = s.authRepositoryPostgres.CreateSession(ctx, session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("save session to database: %w", err)
	}

	session.SetAccessParam(
		accessToken,
		accessClaims.IssuedAt.Time,
		accessClaims.ExpiresAt.Time,
	)
	return session, nil
}
