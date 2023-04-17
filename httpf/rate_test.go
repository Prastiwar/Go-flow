package httpf_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/httpf"
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
			headerNames: []string{httpf.XForwardedForHeader},
			request: &http.Request{
				RemoteAddr: "0.0.0.0",
				Header: http.Header{
					httpf.XForwardedForHeader: []string{"1.1.1.1"},
				},
			},
			key: "1.1.1.1",
		},
		{
			name:        "success-with-header-not-in-request",
			headerNames: []string{httpf.XForwardedForHeader},
			request:     &http.Request{RemoteAddr: "0.0.0.0"},
			key:         "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := httpf.IPRateKey(tt.headerNames...)
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
			factory := httpf.PathRateKey()

			key := factory(tt.request)

			assert.Equal(t, tt.key, key)
		})
	}
}

func TestComposeRateKeyFactories(t *testing.T) {
	t.Run("success-composed", func(t *testing.T) {
		f := httpf.ComposeRateKeyFactories(httpf.IPRateKey(), httpf.PathRateKey())

		key := f(&http.Request{
			Method:     http.MethodGet,
			URL:        urlFromPath("https://test.com/api/resource"),
			RemoteAddr: "0.0.0.0"})

		assert.Equal(t, "0.0.0.0 GET:/api/resource", key)
	})

	t.Run("success-single", func(t *testing.T) {
		f := httpf.ComposeRateKeyFactories(httpf.IPRateKey())
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
		_ = httpf.ComposeRateKeyFactories()
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	timeNow := time.Now()
	timeResetAt := time.Now().Add(time.Second)

	tests := []struct {
		name      string
		store     rate.LimiterStore
		limiter   rate.Limiter
		assertion func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error))
	}{
		{
			name: "success-no-exceed",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func(ctx context.Context) (uint64, error) { return 9, nil },
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse:      func() error { return nil },
						OnResetsAt: func() time.Time { return timeNow },
					}, nil
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 1)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "9", &timeNow, err)
					assert.NilError(t, err)
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-rate-exceed",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func(ctx context.Context) (uint64, error) { return 0, nil },
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse:      func() error { return rate.ErrRateLimitExceeded },
						OnResetsAt: func() time.Time { return timeResetAt },
					}, nil
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "0", &timeResetAt, err)
					assert.ErrorIs(t, err, rate.ErrRateLimitExceeded)
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-unknown-error",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func(ctx context.Context) (uint64, error) { return 10, nil },
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse:      func() error { return errors.New("unknown-error") },
						OnResetsAt: func() time.Time { return timeNow },
					}, nil
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "", "", nil, err)
					assert.ErrorWith(t, err, "unknown-error")
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-handler-error",
			limiter: mocks.LimiterMock{
				OnLimit:  func() uint64 { return 10 },
				OnTokens: func(ctx context.Context) (uint64, error) { return 9, nil },
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return mocks.TokenMock{
						OnUse:      func() error { return nil },
						OnResetsAt: func() time.Time { return timeNow },
					}, nil
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 1)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return errors.New("handler-error")
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "10", "9", &timeNow, err)
					assert.ErrorWith(t, err, "handler-error")
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-limiter-take-error",
			limiter: mocks.LimiterMock{
				OnTake: func(ctx context.Context) (rate.Token, error) {
					return nil, errors.New("limiter-take-error")
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "", "", nil, err)
					assert.ErrorWith(t, err, "limiter-take-error")
					counter.Assert(t)
				}
			},
		},
		{
			name: "failure-store-limit-error",
			store: mocks.LimiterStoreMock{
				OnLimit: func(ctx context.Context, key string) (rate.Limiter, error) {
					return nil, errors.New("store-error")
				},
			},
			assertion: func(t *testing.T) (httpf.HandlerFunc, func(headers http.Header, err error)) {
				counter := assert.Count(t, 0)
				handler := func(w httpf.ResponseWriter, r *http.Request) error {
					counter.Inc()
					return nil
				}
				return handler, func(headers http.Header, err error) {
					assertRateLimitHeaders(t, headers, "", "", nil, err)
					assert.ErrorWith(t, err, "store-error")
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

			store := tt.store
			if store == nil {
				store = mocks.LimiterStoreMock{
					OnLimit: func(ctx context.Context, key string) (rate.Limiter, error) {
						return tt.limiter, nil
					},
				}
			}

			writer := mocks.HttpfResponseWriterMock{
				OnHeader: func() http.Header {
					return headers
				},
			}

			got := httpf.RateLimitMiddleware(handler, store, httpf.IPRateKey())
			err := got.ServeHTTP(writer, request)

			assertion(headers, err)
		})
	}
}

func assertRateLimitHeaders(t *testing.T, header http.Header, limit, remaining string, resetTime *time.Time, err error) {
	const prefix = "incorrect value for header"

	assert.Equal(t, limit, header.Get(httpf.RateLimitLimitHeader), prefix, httpf.RateLimitLimitHeader)
	assert.Equal(t, remaining, header.Get(httpf.RateLimitRemainingHeader), prefix, httpf.RateLimitLimitHeader)

	if resetTime == nil {
		assert.Equal(t, "", header.Get(httpf.RateLimitResetHeader), prefix, httpf.RateLimitLimitHeader)
		assert.Equal(t, "", header.Get(httpf.RetryAfterHeader), prefix, httpf.RetryAfterHeader)
		return
	}

	reset := strconv.FormatInt(int64(resetTime.Unix()), 10)
	assert.Equal(t, reset, header.Get(httpf.RateLimitResetHeader), prefix, httpf.RateLimitLimitHeader)

	if errors.Is(err, rate.ErrRateLimitExceeded) {
		delta := resetTime.Sub(time.Now()).Seconds()
		retryAfter := strconv.FormatInt(int64(delta), 10)
		assert.Equal(t, retryAfter, header.Get(httpf.RetryAfterHeader), prefix, httpf.RetryAfterHeader)
	} else {
		assert.Equal(t, "", header.Get(httpf.RetryAfterHeader), prefix, httpf.RetryAfterHeader)
	}
}

func urlFromPath(path string) *url.URL {
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return url
}
