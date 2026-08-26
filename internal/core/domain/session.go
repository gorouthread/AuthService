package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash []byte
	IsRevoked        bool
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

func NewSession(
	id uuid.UUID,
	userID uuid.UUID,
	refreshTokenHash []byte,
	isRevoked bool,
	createdAt time.Time,
	expiresAt time.Time,
) Session {
	return Session{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: refreshTokenHash,
		IsRevoked:        isRevoked,
		CreatedAt:        createdAt,
		ExpiresAt:        expiresAt,
	}
}
