package httpf

import (
	"net/http"
)

type ErrorHandler interface {
	Handle(w http.ResponseWriter, r *http.Request, err error)
}

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

func (h ErrorHandlerFunc) Handle(w http.ResponseWriter, r *http.Request, err error) {
	h(w, r, err)
}

type Router interface {
	http.Handler

	Handle(pattern string, handler http.Handler)
}

type RouteBuilder interface {
	Get(pattern string, handler Handler) RouteBuilder
	Post(pattern string, handler Handler) RouteBuilder
	Put(pattern string, handler Handler) RouteBuilder
	Delete(pattern string, handler Handler) RouteBuilder
	Patch(pattern string, handler Handler) RouteBuilder
	Options(pattern string, handler Handler) RouteBuilder

	WithErrorHandler(handler ErrorHandler) RouteBuilder
	WithWriterDecorator(decorator func(http.ResponseWriter) ResponseWriter) RouteBuilder

	Build() Router
}
