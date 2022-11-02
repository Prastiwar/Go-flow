package httpf

import (
	"net/http"
	"sync"
)

type serveMuxBuilder struct {
	mu sync.RWMutex

	routes          map[string]map[string]Handler
	errorHandler    ErrorHandler
	writerDecorator func(http.ResponseWriter) ResponseWriter
}

func NewServeMuxBuilder() *serveMuxBuilder {
	return &serveMuxBuilder{
		routes: make(map[string]map[string]Handler),
	}
}

func (b *serveMuxBuilder) Get(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodGet, pattern, handler)
}

func (b *serveMuxBuilder) Post(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPost, pattern, handler)
}

func (b *serveMuxBuilder) Put(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPut, pattern, handler)
}

func (b *serveMuxBuilder) Delete(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodDelete, pattern, handler)
}

func (b *serveMuxBuilder) Patch(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodPatch, pattern, handler)
}

func (b *serveMuxBuilder) Options(pattern string, handler Handler) RouteBuilder {
	return b.handle(http.MethodOptions, pattern, handler)
}

func (b *serveMuxBuilder) WithErrorHandler(handler ErrorHandler) RouteBuilder {
	b.errorHandler = handler
	return b
}

func (b *serveMuxBuilder) WithWriterDecorator(decorator func(http.ResponseWriter) ResponseWriter) RouteBuilder {
	b.writerDecorator = decorator
	return b
}

func (b *serveMuxBuilder) Build() Router {
	mux := http.NewServeMux()
	for route, handlers := range b.routes {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			h, ok := handlers[r.Method]
			if !ok {
				http.Error(w, "", http.StatusMethodNotAllowed)
				return
			}

			var writer ResponseWriter
			if b.writerDecorator == nil {
				writer = &jsonWriterDecorator{w}
			} else {
				writer = b.writerDecorator(w)
			}

			if err := h.ServeHTTP(writer, r); err != nil {
				if b.errorHandler != nil {
					b.errorHandler.Handle(w, r, err)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
	return mux
}

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
