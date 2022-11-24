package httpf

import (
	"net/http"
)

// A ErrorHandler handles error returned from Handler
//
// Handle should write response to the ResponseWriter in common
// used format with proper mapped error.
type ErrorHandler interface {
	Handle(w http.ResponseWriter, r *http.Request, err error)
}

// The ErrorHandlerFunc type is an adapter to allow the use of ordinary
// functions as Error handlers. If h is a function with
// the appropriate signature, ErrorHandlerFunc(h) is a ErrorHandler that calls h.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// Handle calls h(w, r, err)
func (h ErrorHandlerFunc) Handle(w http.ResponseWriter, r *http.Request, err error) {
	h(w, r, err)
}

// A Router is an HTTP request multiplexer. It should match the URL of each incoming
// request against a list of registered patterns and call the handler
// for the pattern that most closely matches the URL.
// Router also should take care of sanitizing the URL request path and the Host
// header, stripping the port number and redirecting any request containing . or
// .. elements or repeated slashes to an equivalent, cleaner URL.
type Router interface {
	http.Handler

	Handle(pattern string, handler http.Handler)
}

// A RouteBuilder is convenient builder for routing registration. It defines
// function for each HTTP Method. Pattern should be able to be registered with
// any method. It's also responsible to use ErrorHandler and WriterDecorator in
// mapping from Handler to http.Handler so errors can be handled gracefully and
// http.ResponseWriter would be decorated with Response function.
type RouteBuilder interface {
	Get(pattern string, handler Handler) RouteBuilder
	Post(pattern string, handler Handler) RouteBuilder
	Put(pattern string, handler Handler) RouteBuilder
	Delete(pattern string, handler Handler) RouteBuilder
	Patch(pattern string, handler Handler) RouteBuilder
	Options(pattern string, handler Handler) RouteBuilder

	WithErrorHandler(handler ErrorHandler) RouteBuilder
	WithWriterDecorator(decorator func(http.ResponseWriter) ResponseWriter) RouteBuilder
	WithParamsParser(parser ParamsParser) RouteBuilder

	Build() Router
}
