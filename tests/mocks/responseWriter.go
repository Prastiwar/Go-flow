package mocks

import (
	"net/http"
)

type ResponseWriter struct {
	OnHeader      func() http.Header
	OnWrite       func([]byte) (int, error)
	OnWriteHeader func(int)
}

func NewResponseWriter(header func() http.Header, write func([]byte) (int, error), writeHeader func(int)) *ResponseWriter {
	return &ResponseWriter{
		OnHeader:      header,
		OnWrite:       write,
		OnWriteHeader: writeHeader,
	}
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
