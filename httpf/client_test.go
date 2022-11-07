package httpf

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

func TestClientSend(t *testing.T) {
	googleUrl, err := url.Parse("https://google.com/")
	assert.NilError(t, err)

	tests := []struct {
		name      string
		ctx       context.Context
		req       *http.Request
		client    func(t *testing.T) Client
		assertion assert.ResultErrorFunc[*http.Response]
	}{
		{
			name: "success-response",
			ctx:  context.TODO(),
			req:  &http.Request{URL: &url.URL{}},
			client: func(t *testing.T) Client {
				roundTripper := &mocks.RoundTripper{
					OnRoundTrip: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "pass",
						}, nil
					},
				}
				return NewClient(WithTransport(roundTripper))
			},
			assertion: func(t *testing.T, result *http.Response, err error) {
				assert.NilError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "pass", result.Status)
			},
		},
		{
			name: "success-timeout",
			ctx:  context.TODO(),
			req:  &http.Request{URL: googleUrl},
			client: func(t *testing.T) Client {
				return NewClient(WithTimeout(200 * time.Millisecond))
			},
			assertion: func(t *testing.T, result *http.Response, err error) {
				assert.ErrorWith(t, err, context.DeadlineExceeded.Error())
				assert.Equal(t, nil, result)
			},
		},
		{
			name: "success-cookies",
			ctx:  context.TODO(),
			req:  &http.Request{URL: googleUrl},
			client: func(t *testing.T) Client {
				counter1 := assert.Count(t, 2, "SetCookies was expected to be called")
				counter2 := assert.Count(t, 2, "Cookies was expected to be called")
				cookies := &mocks.CookiesJar{
					OnSetCookies: func(u *url.URL, cookies []*http.Cookie) {
						counter1.Inc()
					},
					OnCookies: func(u *url.URL) []*http.Cookie {
						counter2.Inc()
						return nil
					},
				}
				return NewClient(WithCookies(cookies))
			},
			assertion: func(t *testing.T, result *http.Response, err error) {
			},
		},
		{
			name: "success-redirect",
			ctx:  context.TODO(),
			req:  &http.Request{URL: googleUrl},
			client: func(t *testing.T) Client {
				return NewClient(WithRedirectHandler(func(req *http.Request, via []*http.Request) error {
					return errors.New("test-redirect")
				}))
			},
			assertion: func(t *testing.T, result *http.Response, err error) {
				assert.ErrorWith(t, err, "test-redirect")
				assert.NotNil(t, result)
				assert.Equal(t, 301, result.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client(t)

			got, err := c.Send(tt.ctx, tt.req)

			c.Close()
			tt.assertion(t, got, err)
		})
	}
}

func TestClientConenientMethods(t *testing.T) {
	tests := []struct {
		name      string
		run       func(t *testing.T, c Client) (*http.Response, error)
		assertion func(t *testing.T, req *http.Request)
	}{
		{
			name: http.MethodGet,
			run: func(t *testing.T, c Client) (*http.Response, error) {
				return c.Get(context.TODO(), "test")
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)
			},
		},
		// TODO: implement
		// {
		// 	name: http.MethodPost,
		// 	// Post(ctx context.Context, url string, body io.Reader) (*http.Response, error)
		// },
		// {
		// 	name: http.MethodPost,
		// 	// PostForm(ctx context.Context, url string, form url.Values) (*http.Response, error)
		// },
		// {
		// 	name: http.MethodPut,
		// 	// Put(ctx context.Context, url string, body io.Reader) (*http.Response, error)
		// },
		// {
		// 	name: http.MethodDelete,
		// 	// Delete(ctx context.Context, url string) (*http.Response, error)
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &mocks.RoundTripper{
				OnRoundTrip: func(r *http.Request) (*http.Response, error) {
					tt.assertion(t, r)
					assert.Equal(t, tt.name, r.Method, "method was invalid")
					return &http.Response{}, nil
				},
			}
			c := NewClient(WithTransport(rt))

			_, _ = tt.run(t, c)
		})
	}
}
