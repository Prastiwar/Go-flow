package rest

import (
	"net/http"
	"sync"
)

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
	mu sync.RWMutex

	router   HttpRouter
	patterns map[string]map[string]HttpHandler
}

func NewFluentRouter(r HttpRouter) HttpRouter {
	return &FluentRouter{
		router:   r,
		patterns: make(map[string]map[string]HttpHandler),
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
		handler, ok = handlers["*"]
		if !ok {
			return NewResponse(http.StatusMethodNotAllowed, http.NoBody, defaultHeaders)
		}
	}

	return handler.Handle(req)
}

func (r *FluentRouter) Register(pattern string, h HttpHandler) {
	r.registerMethod(pattern, "*", h)
}

func (r *FluentRouter) RegisterFunc(pattern string, h func(req HttpRequest) HttpResponse) {
	r.registerMethod(pattern, "*", HttpHandlerFunc(h))
}

func (r *FluentRouter) Get(url string, h HttpHandler) {
	r.registerMethod(url, http.MethodGet, h)
}

func (r *FluentRouter) Post(url string, h HttpHandler) {
	r.registerMethod(url, http.MethodPost, h)
}

func (r *FluentRouter) Delete(url string, h HttpHandler) {
	r.registerMethod(url, http.MethodDelete, h)
}

func (r *FluentRouter) Put(url string, h HttpHandler) {
	r.registerMethod(url, http.MethodPut, h)
}

func (r *FluentRouter) registerMethod(url string, method string, h HttpHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.patterns[url]
	if p == nil {
		p = make(map[string]HttpHandler)
	}
	p[method] = h
	r.patterns[url] = p
}
