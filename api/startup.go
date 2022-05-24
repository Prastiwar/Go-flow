package api

type Bootstrap interface {
	RegisterRoutes(router Router)
	ListenAndServe(server HttpServer)
}
