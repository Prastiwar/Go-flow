package httpf

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
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
	timeNow := time.Now()
	timeNowString := strconv.FormatInt(int64(timeNow.Unix()), 10)

	tests := []struct {
		name      string
		limiter   rate.Limiter
		assertion func(t *testing.T) (HandlerFunc, func(headers http.Header, err error))
	}{
		{
			name: "success-no-exceed",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func() uint64 { return 9 },
				OnTake: func() rate.Token {
					return mocks.TokenMock{
						OnUse:      func() error { return nil },
						OnResetsAt: func() time.Time { return timeNow },
					}
				},
			},
			assertion: func(t *testing.T) (HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 1)
				handler := func(w ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "9", timeNowString)
					assert.NilError(t, err)
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-rate-exceed",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func() uint64 { return 0 },
				OnTake: func() rate.Token {
					return mocks.TokenMock{
						OnUse:      func() error { return rate.ErrRateLimitExceeded },
						OnResetsAt: func() time.Time { return timeNow },
					}
				},
			},
			assertion: func(t *testing.T) (HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "0", timeNowString)
					assert.ErrorIs(t, err, rate.ErrRateLimitExceeded)
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-unknown-error",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func() uint64 { return 10 },
				OnTake: func() rate.Token {
					return mocks.TokenMock{
						OnUse:      func() error { return errors.New("unknown-error") },
						OnResetsAt: func() time.Time { return timeNow },
					}
				},
			},
			assertion: func(t *testing.T) (HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "10", timeNowString)
					assert.ErrorWith(t, err, "unknown-error")
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-handler-error",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func() uint64 { return 9 },
				OnTake: func() rate.Token {
					return mocks.TokenMock{
						OnUse:      func() error { return nil },
						OnResetsAt: func() time.Time { return timeNow },
					}
				},
			},
			assertion: func(t *testing.T) (HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 1)
				handler := func(w ResponseWriter, r *http.Request) error {
					counter.Inc()
					return errors.New("handler-error")
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "9", timeNowString)
					assert.ErrorWith(t, err, "handler-error")
					counter.Assert(t)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header, 3)
			handler, assertion := tt.assertion(t)
			request := &http.Request{RemoteAddr: "0.0.0.0"}

			store := mocks.LimiterStoreMock{
				OnLimit: func(key string) rate.Limiter {
					return tt.limiter
				},
			}

			writer := &jsonWriterDecorator{
				&mocks.ResponseWriter{
					OnHeader: func() http.Header {
						return headers
					},
				},
			}

			got := RateLimitMiddleware(handler, store, IPRateKey())
			err := got.ServeHTTP(writer, request)

			assertion(headers, err)
		})
	}
}

func assertRateLimitHeaders(t *testing.T, header http.Header, limit, remaining, reset string) {
	const prefix = "incorrect value for header"
	assert.Equal(t, limit, header.Get(RateLimitLimitHeader), prefix, RateLimitLimitHeader)
	assert.Equal(t, remaining, header.Get(RateLimitRemainingHeader), prefix, RateLimitLimitHeader)
	assert.Equal(t, reset, header.Get(RateLimitResetHeader), prefix, RateLimitLimitHeader)
}

func urlFromPath(path string) *url.URL {
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return url
}
