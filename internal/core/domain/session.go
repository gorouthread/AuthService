package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
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
	IdempotencyKey   uuid.UUID
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

func NewSessionByRefreshToken(refreshToken string, idempotencyKey string) (Session, error) {
	refreshTokenHash, err := utils.HashString(refreshToken)
	if err != nil {
		return Session{}, fmt.Errorf("hash refresh token: %w", err)
	}

	session := NewSessionUninitialized(
		UninitializedUUID,
		refreshTokenHash,
		UninitializedTime,
		UninitializedTime,
	)

	idempKey, err := uuid.Parse(idempotencyKey)
	if err != nil {
		return Session{}, fmt.Errorf(
			"parse idempotency key string to uuid: %w",
			core_transport_http_response.ErrInvalidArgument,
		)
	}

	session.IdempotencyKey = idempKey
	return session, nil
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
