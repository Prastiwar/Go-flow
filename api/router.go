package api

import (
	"net/http"
)

type Router interface {
	http.Handler
}
