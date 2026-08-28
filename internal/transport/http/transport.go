package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/romreign/AuthService/internal/core/domain"
	core_transport_http_server "github.com/romreign/AuthService/internal/core/transport/http/server"
)

type AuthService interface {
	Register(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.Session, error)
	Logout(ctx context.Context, session domain.Session) error
	Refresh(ctx context.Context, session domain.Session) (domain.Session, error)
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
	}
}
