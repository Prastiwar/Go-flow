package rest

import "context"

type Server interface {
	Run(addr string, h HttpHandler) error
	Shutdown(ctx context.Context) error
	OnShutdown(fn func())
}
