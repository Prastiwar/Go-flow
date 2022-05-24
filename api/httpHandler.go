package api

import (
	"net/http"
)

type HttpHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
