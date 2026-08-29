package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
	"github.com/romreign/AuthService/pkg/utils"
)

func (s *AuthService) Login(ctx context.Context, user domain.User) (domain.Session, error) {
	if err := user.Validate(); err != nil {
		return domain.Session{}, fmt.Errorf("validate user domain: %w", err)
	}

	userFromDb, err := s.authRepositoryPostgres.GetUserByLogin(ctx, user)
	if err != nil {
		return domain.Session{}, fmt.Errorf("get user from database: %w: %w", err, core_errors.ErrNotFound)
	}

	if err := utils.CheckPassword(user.PasswordHash, userFromDb.PasswordHash); err != nil {
		return domain.Session{}, fmt.Errorf("compare password: %w: %w", err, core_errors.ErrNotFound)
	}

	accessToken, accessClaims, err := s.jwtManager.CreateToken(userFromDb.ID, userFromDb.Role, 15*time.Minute)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create access jwt: %w", err)
	}

	refreshToken, refreshClaims, err := s.jwtManager.CreateToken(userFromDb.ID, userFromDb.Role, 24*time.Hour)
	if err != nil {
		return domain.Session{}, fmt.Errorf("create refresh jwt: %w", err)
	}

	hashRefreshToken := utils.HashToken(refreshToken)

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

	session.RefreshToken = refreshToken
	return session, nil
}
