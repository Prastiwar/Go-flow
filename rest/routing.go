package rest

import "net/http"

type httpHandlerFunc func(req HttpRequest) HttpResponse

func (f httpHandlerFunc) Handle(req HttpRequest) HttpResponse {
	return f(req)
}

func HttpHandlerFunc(h func(req HttpRequest) HttpResponse) HttpHandler {
	return httpHandlerFunc(h)
}

type HttpHandler interface {
	Handle(req HttpRequest) HttpResponse
}

type HttpRouter interface {
	HttpHandler

	Register(pattern string, h HttpHandler)
	RegisterFunc(pattern string, h func(req HttpRequest) HttpResponse)
}

type FluentRouter struct {
	patterns map[string]map[string]HttpHandler
	r        HttpRouter
}

func NewFluentRouter(r HttpRouter) HttpRouter {
	return &FluentRouter{
		patterns: make(map[string]map[string]HttpHandler),
		r:        r,
	}
}

func (r *FluentRouter) Handle(req HttpRequest) HttpResponse {
	pattern := req.Url() // TODO: extract pattern
	handlers, ok := r.patterns[pattern]
	if !ok {
		return NotFound()
	}

	handler, ok := handlers[req.Method()]
	if !ok {
		return NewResponse(http.StatusMethodNotAllowed, http.NoBody, defaultHeaders)
	}

	return handler.Handle(req)
}

func (r *FluentRouter) Register(pattern string, h HttpHandler) {
	r.r.Register(pattern, h)
}

func (r *FluentRouter) RegisterFunc(pattern string, h func(req HttpRequest) HttpResponse) {
	r.r.RegisterFunc(pattern, h)
}

func (r *FluentRouter) Get(url string, h HttpHandler) {
	r.patterns[url][http.MethodGet] = h
}

func (r *FluentRouter) Post(url string, h HttpHandler) {
	r.patterns[url][http.MethodPost] = h
}

func (r *FluentRouter) Delete(url string, h HttpHandler) {
	r.patterns[url][http.MethodDelete] = h
}

func (r *FluentRouter) Put(url string, h HttpHandler) {
	r.patterns[url][http.MethodPut] = h
}
