package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
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
	password string,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		ID:           id,
		Login:        login,
		PasswordHash: password,
		Role:         role,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func NewUserUninitialized(
	login string,
	passwordHash string,
) (User, error) {
	user := NewUser(
		UninitializedUUID,
		login,
		passwordHash,
		UninitializedRole,
		UninitializedTime,
		UninitializedTime,
	)

	return user, nil
}

func (u User) Validate() error {
	if len(u.Login) < 3 {
		return fmt.Errorf(
			"Login length < 3: %v: %w",
			u.Login,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
