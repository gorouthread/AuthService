package auth_transport_http

import "github.com/romreign/AuthService/internal/core/domain"

type AuthRequest struct {
	Login    string `json:"login"    validate:"required" example:"username"`
	Password string `json:"password" validate:"required" example:"superPass"`
}

type SessionRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
}

func DomainFromAuthDTO(authRequest AuthRequest, idempotencyKey string) (domain.User, error) {
	return domain.NewUserUninitialized(
		authRequest.Login,
		authRequest.Password,
		idempotencyKey,
	)
}

func DomainFromSessionDTO(sessionRequest SessionRequest, idempotencyKey string) (domain.Session, error) {
	return domain.NewSessionByRefreshToken(
		sessionRequest.RefreshToken,
		idempotencyKey,
	)
}
