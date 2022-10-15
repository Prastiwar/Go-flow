package net

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type request struct {
	req *http.Request
}

func (r *request) Body() io.ReadCloser {
	return r.req.Body
}

func (r *request) ContentLength() int64 {
	return r.req.ContentLength
}

func (r *request) Headers() http.Header {
	return r.req.Header
}

func (r *request) Method() string {
	return r.req.Method
}

func (r *request) Query() url.Values {
	return r.req.URL.Query()
}

func (r *request) Timestamp() time.Time {
	return time.Now() // TODO: how to retrieve real timestamp
}

func (r *request) Url() string {
	return r.req.RequestURI
}
