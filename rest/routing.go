package rest

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpRequest interface {
	Method() string
	Url() string
	Headers() http.Header
	Query() url.Values
	Body() io.ReadCloser
	ContentLength() int64
	Timestamp() time.Time
}

type HttpResponse interface {
	Body() ([]byte, error)
	StatusCode() int
	Headers() http.Header
}

type httpHandlerFunc func(req HttpRequest) HttpResponse

func (f httpHandlerFunc) Handle(req HttpRequest) HttpResponse {
	return f(req)
}

func HttpHandlerFunc(h func(req HttpRequest) HttpResponse) HttpHandler {
	return httpHandlerFunc(h)
}

type HttpHandler interface {
	Handle(req HttpRequest) HttpResponse
}

type Route interface {
	Pattern() string
	QueryParams() string
}

type HttpRouter interface {
	Register(pattern string, h HttpHandler)
	RegisterFunc(pattern string, h func(req HttpRequest) HttpResponse)
}
