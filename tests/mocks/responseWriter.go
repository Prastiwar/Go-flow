package mocks

import (
	"net/http"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ http.ResponseWriter = &ResponseWriter{}

type ResponseWriter struct {
	OnHeader      func() http.Header
	OnWrite       func([]byte) (int, error)
	OnWriteHeader func(int)
	OnResponse    func(int, interface{}) error
}

func (m *ResponseWriter) Header() http.Header {
	assert.ExpectCall(m.OnHeader)
	return m.OnHeader()
}

func (m *ResponseWriter) Write(b []byte) (int, error) {
	assert.ExpectCall(m.OnWrite)
	return m.OnWrite(b)
}

func (m *ResponseWriter) WriteHeader(statusCode int) {
	assert.ExpectCall(m.OnWriteHeader)
	m.OnWriteHeader(statusCode)
}

func (m *ResponseWriter) Response(code int, data interface{}) error {
	assert.ExpectCall(m.OnResponse)
	return m.OnResponse(code, data)
}
