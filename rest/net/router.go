package net

import (
	"net/http"

	"github.com/Prastiwar/Go-flow/rest"
)

type router struct {
	mux *http.ServeMux
}

func (r *router) Register(pattern string, h rest.HttpHandler) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		req := &request{req: r}

		resp := h.Handle(req)

		w.WriteHeader(resp.StatusCode())

		// TODO: improve
		b, _ := resp.Body()
		_, _ = w.Write(b)
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
