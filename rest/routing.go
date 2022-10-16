package rest

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
	Register(pattern string, h HttpHandler)
	RegisterFunc(pattern string, h func(req HttpRequest) HttpResponse)
}
