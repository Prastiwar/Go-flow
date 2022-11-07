package mocks

import (
	"net/http"
)

type RoundTripper struct {
	OnRoundTrip func(*http.Request) (*http.Response, error)
}

func (rt *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.OnRoundTrip(r)
}
