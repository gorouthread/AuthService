package core_transport_http_response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}

func (w *ResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
		w.written = true
	}
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		if w.statusCode == 0 {
			w.WriteHeader(http.StatusOK)
		}
		w.written = true
	}

	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
