package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Prastiwar/Go-flow/rest"
)

type responseWriter struct {
	data       *bytes.Buffer
	headers    http.Header
	statusCode int
}

func (r *responseWriter) Header() http.Header {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	return r.headers
}

func (r *responseWriter) Write(data []byte) (int, error) {
	if r.data == nil {
		r.data = bytes.NewBuffer(data)
		return len(data), nil
	}
	return r.data.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if statusCode < 100 || statusCode > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", statusCode))
	}
	r.statusCode = statusCode
}

func (r *responseWriter) Read(p []byte) (n int, err error) {
	return r.data.Read(p)
}

func (r *responseWriter) Close() error {
	return nil
}

type router struct {
	mux *http.ServeMux
}

func (r *router) Handle(req rest.HttpRequest) rest.HttpResponse {
	httpReq, err := http.NewRequestWithContext(req.Context(), req.Method(), req.Url(), req.Body())
	if err != nil {
		panic(fmt.Errorf("cannot handle request: %w", err))
	}

	responseRw := &responseWriter{}
	r.mux.ServeHTTP(responseRw, httpReq)

	if responseRw.data == nil {
		responseRw.data = bytes.NewBuffer([]byte{})
	}

	if responseRw.statusCode == 0 {
		responseRw.statusCode = 200
	}

	return rest.NewResponse(responseRw.statusCode, responseRw, responseRw.headers)
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
