package auth_transport_http

import (
	"net/http"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type LogoutRequest struct {
	*SessionRequest
}

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(w, log)

	log.Debug("invoke Logout handler")

	var request LogoutRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	idempotencyKey := r.Header.Get("X-Request-ID")
	session, err := DomainFromSessionDTO(*request.SessionRequest, idempotencyKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to logout account")
		return
	}

	err = h.authService.Logout(ctx, session)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to logout account")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusCreated)
}
