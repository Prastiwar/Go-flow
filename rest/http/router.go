package http

import (
	"fmt"
	"net/http"

	"github.com/Prastiwar/Go-flow/rest"
)

type router struct {
	mux *http.ServeMux
}

func (r *router) Handle(req rest.HttpRequest) rest.HttpResponse {
	w, httpReq := r.createRequest(req)

	r.mux.ServeHTTP(w, httpReq)

	return rest.NewResponse(w.StatusCode(), w, w.Header())
}

func (r *router) createRequest(req rest.HttpRequest) (responseReaderWriter, *http.Request) {
	if rawReq, ok := req.(*request); ok {
		return &responseWriterDecorator{
			w: rawReq.writer,
		}, rawReq.r
	}

	httpReq, err := http.NewRequestWithContext(req.Context(), req.Method(), req.Url(), req.Body())
	if err != nil {
		panic(fmt.Errorf("cannot handle request: %w", err))
	}
	return &responseWriter{}, httpReq
}

func (r *router) Register(pattern string, h rest.HttpHandler) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(w, r)

		resp := h.Handle(req)

		writeResponse(w, resp)
	})
}

func (r *router) RegisterFunc(pattern string, h func(req rest.HttpRequest) rest.HttpResponse) {
	r.Register(pattern, rest.HttpHandlerFunc(h))
}

func NewHttpRouter() rest.HttpRouter {
	return &router{
		mux: http.NewServeMux(),
	}
}
