package auth_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type LogoutRequest struct {
	*SessionRequest
}

// Logout        godoc
// @Summary      Logout user
// @Description  Revokes the current user session using the provided refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Idempotency-Key header string true "Unique idempotency key (UUID)"
// @Param        request body LogoutRequest true "Logout request"
// @Success      200 "User successfully logged out"
// @Failure      400 {object} core_transport_http_response.ErrorResponse
// @Failure      401 {object} core_transport_http_response.ErrorResponse
// @Failure      500 {object} core_transport_http_response.ErrorResponse
// @Router       /auth/logout [post]
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

	idempotencyKeyString := r.Header.Get(idempotencyKeyHeader)
	idempotencyKey, err := uuid.Parse(idempotencyKeyString)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf(
			"invalid idempotency key: %w",
			core_errors.ErrInvalidArgument,
		), "failed to parse idempotency key")
		return
	}

	cachedResponse, err := h.idempotencyService.GetIdempotencyKey(
		ctx,
		idempotencyKey,
		r.Method,
		r.URL.Path,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get cached response")
		return
	}

	if cachedResponse != nil {
		responseHandler.JSONResponse(nil, http.StatusOK)
		return
	}

	session := DomainFromSessionDTO(*request.SessionRequest)

	err = h.authService.Logout(ctx, session)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to logout account")
		return
	}

	idempData := &domain.IdempotencyData{
		StatusCode: http.StatusCreated,
		Body:       nil,
		Method:     r.Method,
		URL:        r.URL.Path,
	}
	if err := h.idempotencyService.SaveIdempotencyKey(
		ctx,
		idempotencyKey,
		idempData,
	); err != nil {
		log.Error("failed to save idempotency response", "error", err)
	}

	responseHandler.JSONResponse(nil, http.StatusOK)
}
