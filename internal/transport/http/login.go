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

type LoginRequest struct {
	AuthRequest
}

type LoginResponse struct {
	SessionResponse
}

// Login         godoc
// @Summary      Login user
// @Description  Authenticates a user and creates a new session with access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Idempotency-Key  header    string         true  "Unique idempotency key (UUID)"
// @Param        request          body      LoginRequest   true  "Login credentials"
// @Success      200              {object}  LoginResponse  "User successfully authenticated"
// @Failure      400              {object}  core_transport_http_response.ErrorResponse
// @Failure      401              {object}  core_transport_http_response.ErrorResponse
// @Failure      500              {object}  core_transport_http_response.ErrorResponse
// @Router       /auth/login [post]
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
		var resp LoginResponse
		if err := json.Unmarshal(cachedResponse.Body, &resp); err != nil {
			responseHandler.ErrorResponse(err, "failed to decode and validate response")
			return
		}

		responseHandler.JSONResponse(resp, http.StatusOK)
		return
	}

	user, err := DomainFromAuthDTO(request.AuthRequest)
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
	marshBody, err := json.Marshal(response)
	if err != nil {
		log.Error("failed to marshal idempotency response", "error", err)
		responseHandler.JSONResponse(response, http.StatusCreated)
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

	responseHandler.JSONResponse(response, http.StatusOK)
}

func LoginDTOFromDomain(session domain.Session) LoginResponse {
	return LoginResponse{
		SessionResponse: SessionResponse{
			AccessToken:      session.AccessToken,
			AccessExpiresAt:  session.AccessExpiresAt,
			RefreshToken:     session.RefreshToken,
			RefreshExpiresAt: session.RefreshExpiresAt,
		},
	}
}
