package auth_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash []byte
	IsRevoked        bool
	CreatedAt        time.Time
	ExpiresAt        time.Time
}
