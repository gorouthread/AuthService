package auth_transport_http

import (
	"net/http"
	"time"

	"github.com/romreign/AuthService/internal/core/domain"
	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type RefreshRequest struct {
	*SessionRequest
}

type RefreshResponse struct {
	AccessToken     string    `json:"access_token"       validate:"required" example:"4db118d8-a4df-438a-8128-24a9886393e5"`
	AccessExpiresAt time.Time `json:"access_expired_at"  validate:"required"`
}

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(w, log)

	log.Debug("invoke Refresh handler")

	var request RefreshRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	idempotencyKey := r.Header.Get("X-Request-ID")
	session, err := DomainFromSessionDTO(*request.SessionRequest, idempotencyKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	session, err = h.authService.Refresh(ctx, session)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	resp := RefreshDTOFromDomain(session)
	responseHandler.JSONResponse(resp, http.StatusCreated)
}

func RefreshDTOFromDomain(session domain.Session) RefreshResponse {
	return RefreshResponse{
		AccessToken:     session.AccessToken,
		AccessExpiresAt: session.AccessExpiresAt,
	}
}
