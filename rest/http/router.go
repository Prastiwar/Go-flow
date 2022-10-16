package http

import (
	"net/http"
	"strconv"

	"github.com/Prastiwar/Go-flow/rest"
)

type router struct {
	mux *http.ServeMux
}

func (r *router) Register(pattern string, h rest.HttpHandler) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		req := rest.NewRequest(r.Method, r.URL.RawPath, r.Body, r.Header)

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
}

func (r *router) RegisterFunc(pattern string, h func(req rest.HttpRequest) rest.HttpResponse) {
	r.Register(pattern, rest.HttpHandlerFunc(h))
}

func NewHttpRouter() rest.HttpRouter {
	return &router{
		mux: http.NewServeMux(),
	}
}
