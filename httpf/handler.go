package httpf

import "net/http"

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response. It extends http.ResponseWriter with
// Response function which should be used to share common response format.
type ResponseWriter interface {
	http.ResponseWriter

	Response(code int, data interface{}) error
}

// A jsonWriterDecorator implements ResponseWriter interface and provides
// json writing for Response.
type jsonWriterDecorator struct {
	http.ResponseWriter
}

// Response calls httpf.Json(d, code, data).
func (d *jsonWriterDecorator) Response(code int, data interface{}) error {
	return Json(d, code, data)
}

// A Handler responds to an HTTP request
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return any occurring error. The error should be handler by
// Router which should finish request process.
type Handler interface {
	ServeHTTP(w ResponseWriter, r *http.Request) error
}

// The HandlerFunc type is an adapter to allow the use of ordinary
// functions as HTTP handlers. If h is a function with
// the appropriate signature, HandlerFunc(h) is a Handler that calls h.
type HandlerFunc func(w ResponseWriter, r *http.Request) error

// ServeHTTP calls h(w, r).
func (h HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request) error {
	return h(w, r)
}
