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

func NewResponse(code int, body io.ReadCloser, headers http.Header) *response {
	return &response{
		code:    code,
		body:    body,
		headers: headers,
	}
}

// TODO:
// func WithBody(body io.ReadCloser) Option {
// 	return func(resp *response) {
// 		resp.body = body
// 	}
// }

func Ok() HttpResponse {
	return NewResponse(http.StatusOK, http.NoBody, defaultHeaders)
}

func Created() HttpResponse {
	return NewResponse(http.StatusCreated, http.NoBody, defaultHeaders)
}

func NoContent() HttpResponse {
	return NewResponse(http.StatusNoContent, http.NoBody, defaultHeaders)
}

func BadRequest() HttpResponse {
	return NewResponse(http.StatusBadRequest, http.NoBody, defaultHeaders)
}

func Unprocessable() HttpResponse {
	return NewResponse(http.StatusUnprocessableEntity, http.NoBody, defaultHeaders)
}

func NotFound() HttpResponse {
	return NewResponse(http.StatusNotFound, http.NoBody, defaultHeaders)
}

func Forbidden() HttpResponse {
	return NewResponse(http.StatusForbidden, http.NoBody, defaultHeaders)
}

func Unathorized() HttpResponse {
	return NewResponse(http.StatusUnauthorized, http.NoBody, defaultHeaders)
}

func InternalServerError() HttpResponse {
	return NewResponse(http.StatusInternalServerError, http.NoBody, defaultHeaders)
}
