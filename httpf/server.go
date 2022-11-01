package httpf

import (
	"context"
	"net"
	"net/http"
)

type Server interface {
	Close() error
	Shutdown(ctx context.Context) error
	RegisterOnShutdown(f func())

	ListenAndServe() error
	Serve(l net.Listener) error

	ListenAndServeTLS(certFile, keyFile string) error
	ServeTLS(l net.Listener, certFile, keyFile string) error
}

func NewServer(addr string, router Router) *http.Server {
	return &http.Server{
		Handler: router,
	}
}
