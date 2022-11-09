package mocks

import (
	"net/http"
)

type ResponseWriter struct {
	OnHeader      func() http.Header
	OnWrite       func([]byte) (int, error)
	OnWriteHeader func(int)
	OnResponse    func(int, interface{}) error
}

func (w *ResponseWriter) Header() http.Header {
	return w.OnHeader()
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.OnWrite(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.OnWriteHeader(statusCode)
}

func (w *ResponseWriter) Response(code int, data interface{}) error {
	return w.OnResponse(code, data)
}
