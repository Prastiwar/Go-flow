package http

import (
	"context"
	"net/http"

	"github.com/Prastiwar/Go-flow/rest"
)

type httpServer struct {
	server *http.Server
	router rest.HttpRouter
}

func (s *httpServer) Run(addr string, h rest.HttpHandler) error {
	s.server.Addr = addr
	s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(w, r)

		resp := h.Handle(req)

		writeResponse(w, resp)
	})
	return s.server.ListenAndServe()
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *httpServer) OnShutdown(fn func()) {
	s.server.RegisterOnShutdown(fn)
}

func NewServer(server *http.Server, router rest.HttpRouter) rest.Server {
	return &httpServer{
		server: server,
		router: router,
	}
}
