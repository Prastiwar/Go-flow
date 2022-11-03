package httpf

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A Client is an HTTP client containing convenient API to send request with common
// HTTP methods. Send function is the fundamental implementation for the Client which
// provides a way to send request over HTTP and receive response
type Client interface {
	Send(ctx context.Context, req *http.Request) (*http.Response, error)

	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	PostForm(ctx context.Context, url string, form url.Values) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	Delete(ctx context.Context, url string) (*http.Response, error)

	Close()
}

// ClientOptions defines http.Client constructor parameters which can be set on NewClient
type ClientOptions struct {
	Transport     http.RoundTripper
	CheckRedirect func(req *http.Request, via []*http.Request) error
	Jar           http.CookieJar
	Timeout       time.Duration
}

// ClientOption defines single function to mutate options
type ClientOption func(*ClientOptions)

// NewClientOptions returns a new instance of ClientOptions with is result of merged ClientOption slice
func NewClientOptions(opts ...ClientOption) ClientOptions {
	o := &ClientOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return *o
}

// WithTransport sets option which specifies the mechanism by which individual HTTP requests are made
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(o *ClientOptions) {
		o.Transport = transport
	}
}

// WithRedirectHandler sets option which specifies the policy for handling redirects
func WithRedirectHandler(handler func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(o *ClientOptions) {
		o.CheckRedirect = handler
	}
}

// WithCookies sets option which specifies cookie jar used to insert relevant cookies
// into every outbound Request and is updated with the cookie values of every inbound Response
func WithCookies(cookieJar http.CookieJar) ClientOption {
	return func(o *ClientOptions) {
		o.Jar = cookieJar
	}
}

// WithTimeout sets option which specifies a time limit for requests made by Client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// A client is Client adapter for http.Client
type client struct {
	c *http.Client
}

// NewClient returns a new instace of Client which is adapter for http.Client. Provided
// options can be set optionally to pass the values in http.Client construction
func NewClient(opts ...ClientOption) *client {
	o := NewClientOptions(opts...)

	return &client{
		c: &http.Client{
			Transport:     o.Transport,
			CheckRedirect: o.CheckRedirect,
			Jar:           o.Jar,
			Timeout:       o.Timeout,
		},
	}
}

// Send calls http.Client Do function using request with given context
func (c *client) Send(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.c.Do(req.WithContext(ctx))
}

// Get sends GET request
func (c *client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req)
}

// Post sends POST request using application/json Content-Type as default value. To use different type
// use Send with request containing appropiate Content-Type header
func (c *client) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ContentTypeHeader, ApplicationJsonType)
	return c.Send(ctx, req)
}

// PostForm sends POST request using application/x-www-form-urlencoded Content-Type and encoded form values as body
func (c *client) PostForm(ctx context.Context, url string, form url.Values) (*http.Response, error) {
	body := io.NopCloser(strings.NewReader(form.Encode()))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ContentTypeHeader, ApplicationFormEncodedType)
	return c.Send(ctx, req)
}

// Put sends PUT request using application/json Content-Type as default value. To use different type
// use Send with request containing appropiate Content-Type header
func (c *client) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ContentTypeHeader, ApplicationJsonType)
	return c.Send(ctx, req)
}

// Delete sends DELETE request
func (c *client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req)
}

// Close calls http.Client CloseIdleConnections function
func (c *client) Close() {
	c.c.CloseIdleConnections()
}
