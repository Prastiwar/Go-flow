package http

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Prastiwar/Go-flow/rest"
)

type request struct {
	r         *http.Request
	writer    http.ResponseWriter
	timestamp time.Time
}

func (r *request) Body() io.ReadCloser {
	return r.r.Body
}

func (r *request) ContentLength() int64 {
	return r.r.ContentLength
}

func (r *request) Context() context.Context {
	return r.r.Context()
}

func (r *request) Headers() http.Header {
	return r.r.Header
}

func (r *request) Method() string {
	return r.r.Method
}

func (r *request) Query() url.Values {
	return r.r.URL.Query()
}

func (r *request) Timestamp() time.Time {
	return r.timestamp
}

func (r *request) Url() string {
	return r.r.URL.Path
}

func newRequest(writer http.ResponseWriter, r *http.Request) rest.HttpRequest {
	return &request{
		r:         r,
		writer:    writer,
		timestamp: time.Now(),
	}
}
