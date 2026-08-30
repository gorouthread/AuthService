package auth_transport_http

import (
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
)

type AuthRequest struct {
	Login    string `json:"login"    validate:"required,min=3" example:"username"`
	Password string `json:"password" validate:"required"       example:"superPass"`
}

type SessionRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
}

type SessionResponse struct {
	AccessToken      string    `json:"access_token"       validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
	AccessExpiresAt  time.Time `json:"access_expired_at"  validate:"required"`
	RefreshToken     string    `json:"refresh_token"      validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
	RefreshExpiresAt time.Time `json:"refresh_expired_at" validate:"required"`
}

func DomainFromAuthDTO(authRequest AuthRequest) (domain.User, error) {
	return domain.NewUserUninitialized(
		authRequest.Login,
		authRequest.Password,
	)
}

func DomainFromSessionDTO(sessionRequest SessionRequest) domain.Session {
	return domain.NewSessionByRefreshToken(
		sessionRequest.RefreshToken,
	)
}
