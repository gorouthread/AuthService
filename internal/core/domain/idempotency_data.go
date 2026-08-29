package domain

type IdempotencyData struct {
	Method     string `json:"method"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
}
