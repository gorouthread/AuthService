package core_transport_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
)

type HTTPResponseHandler struct {
	w   http.ResponseWriter
	log *core_logger.Logger
}

func NewHTTPResponseHandler(w http.ResponseWriter, log *core_logger.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		w:   w,
		log: log,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)
	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write http response body", "error", err)
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Debug

	case errors.Is(err, ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	h.errorResponse(statusCode, err, msg, logFunc)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.errorResponse(statusCode, err, msg, h.log.Error)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
	logFunc func(string, ...any),
) {
	logFunc(msg, "error", err)

	response := ErrorResponse{
		Err: err,
		Msg: msg,
	}

	h.JSONResponse(response, statusCode)
}
