package rest

import (
	"io"
	"net/http"
)

var (
	defaultHeaders = http.Header{}
)

type HttpResponse interface {
	Body() io.ReadCloser
	StatusCode() int
	Headers() http.Header
}

type ResponseOptions struct {
	Body    io.ReadCloser
	Headers http.Header
}

type ResponseOption func(o *ResponseOptions)

func NewResponseOptions(options ...ResponseOption) ResponseOptions {
	opts := &ResponseOptions{}
	for _, v := range options {
		v(opts)
	}
	if opts.Body == nil {
		opts.Body = http.NoBody
	}
	if opts.Headers == nil {
		opts.Headers = defaultHeaders
	}
	return *opts
}

func WithBody(body io.ReadCloser) ResponseOption {
	return func(o *ResponseOptions) {
		o.Body = body
	}
}

func WithHeader(key, value string) ResponseOption {
	return func(o *ResponseOptions) {
		o.Headers.Add(key, value)
	}
}

func WithHeaders(headers http.Header) ResponseOption {
	return func(o *ResponseOptions) {
		o.Headers = headers
	}
}

type response struct {
	code    int
	body    io.ReadCloser
	headers http.Header
}

func (r *response) Body() io.ReadCloser {
	return r.body
}

func (r *response) Headers() http.Header {
	return r.headers
}

func (r *response) StatusCode() int {
	return r.code
}

func NewResponse(code int, options ...ResponseOption) *response {
	opts := NewResponseOptions(options...)
	return &response{
		code:    code,
		body:    opts.Body,
		headers: opts.Headers,
	}
}

func Ok(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusOK, options...)
}

func Created(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusCreated, options...)
}

func NoContent(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusNoContent, options...)
}

func BadRequest(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusBadRequest, options...)
}

func Unprocessable(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusUnprocessableEntity, options...)
}

func NotFound(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusNotFound, options...)
}

func Forbidden(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusForbidden, options...)
}

func Unathorized(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusUnauthorized, options...)
}

func InternalServerError(options ...ResponseOption) HttpResponse {
	return NewResponse(http.StatusInternalServerError, options...)
}
