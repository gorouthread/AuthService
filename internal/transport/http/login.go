package auth_transport_http

import (
	"net/http"
)

type LoginRequest struct {
	AuthRequest
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {

}
