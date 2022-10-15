package rest

import "net/http"

type response struct {
	code    int
	body    []byte
	headers http.Header
}

func (r *response) Body() ([]byte, error) {
	return r.body, nil
}

func (r *response) Headers() http.Header {
	return r.headers
}

func (r *response) StatusCode() int {
	return r.code
}

func newResponse(code int, body []byte, h http.Header) *response {
	return &response{
		code:    code,
		body:    body,
		headers: h,
	}
}

func NotFound() HttpResponse {
	// TODO: Default headers
	return newResponse(http.StatusNotFound, []byte{}, http.Header{})
}

// TODO: BadRequest, Forbidden, InternalServerError
