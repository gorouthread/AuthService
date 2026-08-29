package auth_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/romreign/AuthService/internal/core/domain"
	core_errors "github.com/romreign/AuthService/internal/core/errors"
	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_request "github.com/romreign/AuthService/internal/core/transport/http/request"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

type RefreshRequest struct {
	SessionRequest
}

type RefreshResponse struct {
	SessionResponse
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
		var resp RefreshResponse
		if err := json.Unmarshal(cachedResponse.Body, &resp); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate response")
			return
		}

		responseHandler.JSONResponse(resp, http.StatusOK)
		return
	}

	session := DomainFromSessionDTO(request.SessionRequest)

	session, err = h.authService.Refresh(ctx, session)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	resp := RefreshDTOFromDomain(session)
	marshBody, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed to marshal idempotency response", "error", err)
		responseHandler.JSONResponse(resp, http.StatusCreated)
		return
	}

	idempData := &domain.IdempotencyData{
		StatusCode: http.StatusCreated,
		Body:       marshBody,
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

	responseHandler.JSONResponse(resp, http.StatusCreated)
}

func RefreshDTOFromDomain(session domain.Session) RefreshResponse {
	return RefreshResponse{
		SessionResponse: SessionResponse{
			AccessToken:      session.AccessToken,
			AccessExpiresAt:  session.AccessExpiresAt,
			RefreshToken:     session.RefreshToken,
			RefreshExpiresAt: session.RefreshExpiresAt,
		},
	}
}
