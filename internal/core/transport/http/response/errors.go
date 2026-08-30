package core_transport_http_response

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
