package auth_transport_http

import (
	"net/http"
)

type RegisterRequest struct {
	AuthRequest
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {

}
