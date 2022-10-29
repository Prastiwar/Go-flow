package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Prastiwar/Go-flow/rest"
)

type httpServer struct {
	server *http.Server
	router rest.HttpRouter
}

func (s *httpServer) Run(addr string, h rest.HttpHandler) error {
	s.server.Addr = addr
	s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := rest.NewRequest(r.Method, r.URL.Path, r.Body, r.Header).WithContext(r.Context())

		resp := h.Handle(req)

		for key, values := range resp.Headers() {
			for _, val := range values {
				w.Header().Add(key, val)
			}
		}
		w.WriteHeader(resp.StatusCode())

		contentLengthHeader := resp.Headers().Get(rest.ContentLengthHeader)
		contentLength, _ := strconv.ParseInt(contentLengthHeader, 10, 0)
		bytes := make([]byte, 0, contentLength)

		// TODO: handle errors
		_, err := resp.Body().Read(bytes)
		if err == nil {
			_, _ = w.Write(bytes)
		}
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
