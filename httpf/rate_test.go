package httpf

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestIPRateKey(t *testing.T) {
	tests := []struct {
		name        string
		headerNames []string
		request     *http.Request
		key         string
	}{
		{
			name:        "success-no-headers",
			headerNames: make([]string, 0),
			request:     &http.Request{RemoteAddr: "0.0.0.0"},
			key:         "0.0.0.0",
		},
		{
			name:        "success-with-header",
			headerNames: []string{XForwardedForHeader},
			request: &http.Request{
				RemoteAddr: "0.0.0.0",
				Header: http.Header{
					XForwardedForHeader: []string{"1.1.1.1"},
				},
			},
			key: "1.1.1.1",
		},
		{
			name:        "success-with-header-not-in-request",
			headerNames: []string{XForwardedForHeader},
			request:     &http.Request{RemoteAddr: "0.0.0.0"},
			key:         "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := IPRateKey(tt.headerNames...)
			key := factory(tt.request)
			assert.Equal(t, tt.key, key)
		})
	}
}

func TestPathRateKey(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
		key     string
	}{
		{
			name: "success-get",
			request: &http.Request{
				Method: http.MethodGet,
				URL:    urlFromPath("https://test.com/api/smth"),
			},
			key: "GET:/api/smth",
		},
		{
			name: "success-post",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    urlFromPath("https://test.com/api/smth/"),
			},
			key: "POST:/api/smth/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := PathRateKey()

			key := factory(tt.request)

			assert.Equal(t, tt.key, key)
		})
	}
}

func TestComposeRateKeyFactories(t *testing.T) {
	t.Run("success-composed", func(t *testing.T) {
		f := ComposeRateKeyFactories(IPRateKey(), PathRateKey())

		key := f(&http.Request{
			Method:     http.MethodGet,
			URL:        urlFromPath("https://test.com/api/resource"),
			RemoteAddr: "0.0.0.0"})

		assert.Equal(t, "0.0.0.0 GET:/api/resource", key)
	})

	t.Run("success-single", func(t *testing.T) {
		f := ComposeRateKeyFactories(IPRateKey())
		key := f(&http.Request{
			Method:     http.MethodGet,
			URL:        urlFromPath("https://test.com/api/resource"),
			RemoteAddr: "0.0.0.0"})
		assert.Equal(t, "0.0.0.0", key)
	})

	t.Run("failure-no-factories", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = ComposeRateKeyFactories()
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name string
		h    HandlerFunc
		r    *http.Request
		err  error
	}{
		{
			name: "success-no-exceed",
			// TODO: check if handler is called
		},
		{
			name: "failure-rate-exceed",
			// TODO: check if handler is not callded
			// TODO: check if headers are written
		},
		{
			name: "failure-unknown-error",
			// TODO: check if handler is not callded
			// TODO: check if headers are not written
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := rate.LimiterStoreMock{
				OnLimit: func(key string) rate.Limiter {
					return rate.LimiterMock{
						OnUse: func() rate.Token {
							return rate.TokenMock{
								OnUse: func() error {
									return tt.err
								},
							}
						},
					}
				},
			}
			writer := &jsonWriterDecorator{}

			got := RateLimitMiddleware(tt.h, store, IPRateKey())
			err := got.ServeHTTP(writer, tt.r)

			assert.Equal(t, tt.err, err)
		})
	}
}

func urlFromPath(path string) *url.URL {
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return url
}
