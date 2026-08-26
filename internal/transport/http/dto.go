package auth_transport_http

import "github.com/romreign/AuthService/internal/core/domain"

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func DomainFromDTO(authRequest AuthRequest) (domain.User, error) {
	return domain.User{}, nil
}
