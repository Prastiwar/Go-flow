package httpf

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter

	Response(code int, data interface{}) error
}

type jsonWriterDecorator struct {
	http.ResponseWriter
}

func (d *jsonWriterDecorator) Response(code int, data interface{}) error {
	return Json(d, code, data)
}

type Request interface {
	// TODO: implement
	PathParam(key string)
}

type Handler interface {
	ServeHTTP(w ResponseWriter, r *http.Request) error
}

type HandlerFunc func(w ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request) error {
	return h(w, r)
}
