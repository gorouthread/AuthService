package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/pkg/utils"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	AccessToken      string
	RefreshToken     string
	IsRevoked        bool
	AccessCreatedAt  time.Time
	AccessExpiresAt  time.Time
	RefreshCreatedAt time.Time
	RefreshExpiresAt time.Time
}

func NewSession(
	id uuid.UUID,
	userID uuid.UUID,
	refreshToken string,
	isRevoked bool,
	refreshCreatedAt time.Time,
	refreshExpiresAt time.Time,
) Session {
	return Session{
		ID:               id,
		UserID:           userID,
		RefreshToken:     refreshToken,
		IsRevoked:        isRevoked,
		RefreshCreatedAt: refreshCreatedAt,
		RefreshExpiresAt: refreshExpiresAt,
	}
}

func NewSessionUninitialized(
	userID uuid.UUID,
	refreshToken string,
	refreshCreatedAt time.Time,
	refreshExpiresAt time.Time,
) Session {
	return Session{
		ID:               UninitializedUUID,
		UserID:           userID,
		RefreshToken:     refreshToken,
		IsRevoked:        false,
		RefreshCreatedAt: refreshCreatedAt,
		RefreshExpiresAt: refreshExpiresAt,
	}
}

func NewSessionByRefreshToken(refreshToken string) Session {
	refreshTokenHash := utils.HashToken(refreshToken)

	return NewSessionUninitialized(
		UninitializedUUID,
		refreshTokenHash,
		UninitializedTime,
		UninitializedTime,
	)
}

func (s *Session) SetAccessParam(
	accessToken string,
	accessCreatedAt time.Time,
	accessExpiresAt time.Time,
) {
	s.AccessToken = accessToken
	s.AccessCreatedAt = accessCreatedAt
	s.AccessExpiresAt = accessExpiresAt
}
