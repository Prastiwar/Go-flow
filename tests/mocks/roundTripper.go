package mocks

import (
	"net/http"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ http.RoundTripper = &RoundTripper{}

type RoundTripper struct {
	OnRoundTrip func(*http.Request) (*http.Response, error)
}

func (m *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	assert.ExpectCall(m.OnRoundTrip)
	return m.OnRoundTrip(r)
}
