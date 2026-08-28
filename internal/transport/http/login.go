package auth_transport_http

import (
	"net/http"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type LoginRequest struct {
	*AuthRequest
}

type LoginResponse struct {
	AccessToken      string    `json:"access_token"       validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
	AccessExpiresAt  time.Time `json:"access_expired_at"  validate:"required"`
	RefreshToken     string    `json:"refresh_token"      validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
	RefreshExpiresAt time.Time `json:"refresh_expired_at" validate:"required"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(w, log)

	log.Debug("invoke Login handler")

	var request LoginRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	idempotencyKey := r.Header.Get("X-Request-ID")

	user, err := DomainFromAuthDTO(*request.AuthRequest, idempotencyKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login account")
		return
	}

	session, err := h.authService.Login(ctx, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login account")
		return
	}

	response := LoginDTOFromDomain(session)
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func LoginDTOFromDomain(session domain.Session) LoginResponse {
	return LoginResponse{
		AccessToken:      session.AccessToken,
		AccessExpiresAt:  session.AccessExpiresAt,
		RefreshToken:     session.RefreshToken,
		RefreshExpiresAt: session.RefreshExpiresAt,
	}
}
