package auth_repository_redis

import (
	"time"

	"github.com/google/uuid"
)

const (
	tokenPrefix       = "auth:token:%s"
	idempotencyPrefix = "idempotency:%s"
)

type AuthData struct {
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"roles"`
	ExpiresAt time.Time `json:"created_at"`
}

type IdempotencyResponse struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
}
