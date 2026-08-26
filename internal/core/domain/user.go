package domain

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

func NewUser(
	id uuid.UUID,
	login string,
	passwordHash string,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func NewUserUninitialized(
	login string,
	passwordHash string,
) User {
	return NewUser(
		UninitializedUUID,
		login,
		passwordHash,
		UninitializedRole,
		UninitializedTime,
		UninitializedTime,
	)
}
