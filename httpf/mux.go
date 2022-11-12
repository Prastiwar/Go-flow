package httpf

import (
	"net/http"
	"sync"
)

type routeHandler struct {
	pattern         string
	handlers        map[string]Handler
	writerDecorator func(http.ResponseWriter) ResponseWriter
	paramsParser    ParamsParser
	errorHandler    ErrorHandler
}

func (r *routeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, ok := r.handlers[req.Method]
	if !ok {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	writer := r.writerDecorator(w)

	if r.paramsParser != nil {
		pathParams := r.paramsParser.ParseParams(req)
		req = WithParams(req, pathParams)
	}

	if err := h.ServeHTTP(writer, req); err != nil {
		r.errorHandler.Handle(w, req, err)
	}
}

// serveMuxBuilder implements RouterBuilder interface with
// building the http.ServeMux router.
type serveMuxBuilder struct {
	mu sync.RWMutex

	routes          map[string]map[string]Handler
	errorHandler    ErrorHandler
	writerDecorator func(http.ResponseWriter) ResponseWriter
	paramsParser    ParamsParser
}

// NewServeMuxBuilder returns RouterBuilder which build results in adapting
// http.ServeMux implementation to handle errors, decorate http.ResponseWriter or use ParamsParser.
// Note http.ServeMux does not support defining parameters in pattern.
// For default behaviour of corresponding With.. option can be found in option func comment.
func NewServeMuxBuilder() *serveMuxBuilder {
	return &serveMuxBuilder{
		routes: make(map[string]map[string]Handler),
	}
}

// Get registers handler to pattern using GET method
func (b *serveMuxBuilder) Get(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodGet, pattern, handler)
}

// Post registers handler to pattern using POST method
func (b *serveMuxBuilder) Post(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPost, pattern, handler)
}

// Put registers handler to pattern using PUT method
func (b *serveMuxBuilder) Put(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPut, pattern, handler)
}

// Delete registers handler to pattern using DELETE method
func (b *serveMuxBuilder) Delete(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodDelete, pattern, handler)
}

// Patch registers handler to pattern using PATCH method
func (b *serveMuxBuilder) Patch(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPatch, pattern, handler)
}

// Options registers handler to pattern using OPTIONS method
func (b *serveMuxBuilder) Options(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodOptions, pattern, handler)
}

// WithErrorHandler sets ErrorHandler used in Build. If will not be provided Router will
// write response using http.Error with http.StatusInternalServerError
func (b *serveMuxBuilder) WithErrorHandler(handler ErrorHandler) RouteBuilder {
	b.errorHandler = handler
	return b
}

// WithWriterDecorator sets function which should decorate http.ResponseWriter coming from handler. If will not be provided
// Router will use json writer decorator
func (b *serveMuxBuilder) WithWriterDecorator(decorator func(http.ResponseWriter) ResponseWriter) RouteBuilder {
	b.writerDecorator = decorator
	return b
}

// WithParamsParser sets parser which which should inject parsed path parameters to http request. If will not be provided
// httpf.Params will always return empty map without error
func (b *serveMuxBuilder) WithParamsParser(parser ParamsParser) RouteBuilder {
	b.paramsParser = parser
	return b
}

// Build registers the registered handlers in builder to http.ServeMux using mux.HandleFunc
// which matches accurate HTTP method or returns MethodNotAllowed status. It also wraps handler with
// proper error handling and decorating incoming http.ResponseWriter.
// If ResponseWriter decorator was not set jsonWriterDecorator is used instead.
// If ErrorHandler was not set just http.Error is called with Internal Server status.
func (b *serveMuxBuilder) Build() Router {
	if b.writerDecorator == nil {
		b.writerDecorator = func(w http.ResponseWriter) ResponseWriter { return &jsonWriterDecorator{w} }
	}

	if b.errorHandler == nil {
		b.errorHandler = ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		})
	}

	mux := http.NewServeMux()
	for route, handlers := range b.routes {
		r := &routeHandler{
			pattern:         route,
			handlers:        handlers,
			writerDecorator: b.writerDecorator,
			paramsParser:    b.paramsParser,
			errorHandler:    b.errorHandler,
		}
		mux.Handle(route, r)
	}
	return mux
}

// handle registers handler to given pattern with method
func (b *serveMuxBuilder) handle(method string, pattern string, h Handler) RouteBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	p := b.routes[pattern]
	if p == nil {
		p = make(map[string]Handler)
	}
	p[method] = h
	b.routes[pattern] = p
	return b
}
