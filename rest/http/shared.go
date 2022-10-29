package http

import (
	"net/http"
	"strconv"

	"github.com/Prastiwar/Go-flow/rest"
)

func writeResponse(w http.ResponseWriter, r rest.HttpResponse) {
	for key, values := range r.Headers() {
		for _, val := range values {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(r.StatusCode())

	contentLengthHeader := r.Headers().Get(rest.ContentLengthHeader)
	contentLength, _ := strconv.ParseInt(contentLengthHeader, 10, 0)
	bytes := make([]byte, 0, contentLength)

	r.Body().Read(bytes)
	w.Write(bytes)
}
