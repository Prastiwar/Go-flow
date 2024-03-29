package httpf_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/httpf"
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
		client    func(t *testing.T) httpf.Client
		assertion assert.ResultErrorFunc[*http.Response]
	}{
		{
			name: "success-response",
			ctx:  context.TODO(),
			req:  &http.Request{URL: &url.URL{}},
			client: func(t *testing.T) httpf.Client {
				roundTripper := &mocks.RoundTripper{
					OnRoundTrip: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "pass",
						}, nil
					},
				}
				return httpf.NewClient(httpf.WithTransport(roundTripper))
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
			client: func(t *testing.T) httpf.Client {
				return httpf.NewClient(httpf.WithTimeout(time.Millisecond))
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
			client: func(t *testing.T) httpf.Client {
				setCookiesCounter := assert.Count(t, 1, "SetCookies was expected to be called").AtLeast()
				cookiesCounter := assert.Count(t, 1, "Cookies was expected to be called").AtLeast()
				cookies := &mocks.CookiesJar{
					OnSetCookies: func(u *url.URL, cookies []*http.Cookie) {
						setCookiesCounter.Inc()
					},
					OnCookies: func(u *url.URL) []*http.Cookie {
						cookiesCounter.Inc()
						return nil
					},
				}
				return httpf.NewClient(httpf.WithCookies(cookies))
			},
			assertion: func(t *testing.T, result *http.Response, err error) {
			},
		},
		{
			name: "success-redirect",
			ctx:  context.TODO(),
			req:  &http.Request{URL: googleUrl},
			client: func(t *testing.T) httpf.Client {
				return httpf.NewClient(httpf.WithRedirectHandler(func(req *http.Request, via []*http.Request) error {
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
		run       func(t *testing.T, c httpf.Client, url string) (*http.Response, error)
		assertion func(t *testing.T, req *http.Request)
	}{
		{
			name: http.MethodGet,
			run: func(t *testing.T, c httpf.Client, url string) (*http.Response, error) {
				return c.Get(context.TODO(), url)
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)
			},
		},
		{
			name: http.MethodPost,
			run: func(t *testing.T, c httpf.Client, url string) (*http.Response, error) {
				body := bytes.NewBufferString("test-body")
				return c.Post(context.TODO(), url, body)
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)

				data, err := io.ReadAll(req.Body)
				assert.NilError(t, err)
				assert.Equal(t, "test-body", string(data))
			},
		},
		{
			name: http.MethodPost,
			run: func(t *testing.T, c httpf.Client, urlPath string) (*http.Response, error) {
				form := url.Values{
					"test": []string{"form"},
				}
				return c.PostForm(context.TODO(), urlPath, form)
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)

				data, err := io.ReadAll(req.Body)
				assert.NilError(t, err)
				assert.Equal(t, "test=form", string(data))
			},
		},
		{
			name: http.MethodPut,
			run: func(t *testing.T, c httpf.Client, url string) (*http.Response, error) {
				body := bytes.NewBufferString("test-body")
				return c.Put(context.TODO(), url, body)
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)

				data, err := io.ReadAll(req.Body)
				assert.NilError(t, err)
				assert.Equal(t, "test-body", string(data))
			},
		},
		{
			name: http.MethodDelete,
			run: func(t *testing.T, c httpf.Client, url string) (*http.Response, error) {
				return c.Delete(context.TODO(), url)
			},
			assertion: func(t *testing.T, req *http.Request) {
				assert.NotNil(t, req)
			},
		},
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
			c := httpf.NewClient(httpf.WithTransport(rt))

			_, _ = tt.run(t, c, "test")
		})
	}

	for _, tt := range tests {
		t.Run(tt.name+"-error", func(t *testing.T) {
			c := httpf.NewClient()

			_, _ = tt.run(t, c, "!@#$%^&*()")
		})
	}
}
