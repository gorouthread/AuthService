package auth_transport_http

import (
	"net/http"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type RegisterRequest struct {
	*AuthRequest
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(w, log)

	log.Debug("invoke Register handler")

	var request RegisterRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	idempotencyKey := r.Header.Get("X-Request-ID")

	user, err := DomainFromAuthDTO(*request.AuthRequest, idempotencyKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create payment")
		return
	}

	err = h.authService.Register(ctx, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create payment")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusCreated)
}
