package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type User struct {
	ID             uuid.UUID
	IdempotencyKey uuid.UUID
	Login          string
	PasswordHash   string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	idempotencyKey string,
) (User, error) {
	user := NewUser(
		UninitializedUUID,
		login,
		passwordHash,
		UninitializedRole,
		UninitializedTime,
		UninitializedTime,
	)

	idempKey, err := uuid.Parse(idempotencyKey)
	if err != nil {
		return User{}, fmt.Errorf(
			"parse idempotency key string to uuid: %w",
			core_transport_http_response.ErrInvalidArgument,
		)
	}

	user.IdempotencyKey = idempKey

	return user, nil
}

func (u User) Validate() error {
	if len(u.Login) < 3 {
		return fmt.Errorf(
			"Login length < 3: %v: %w",
			u.Login,
			core_transport_http_response.ErrInvalidArgument,
		)
	}

	return nil
}
