package httpf

import (
	"context"
	"net/http"
)

type paramsKeyValue struct{}

var paramsKey = &paramsKeyValue{}

// ParamsParser is parser for request to retrieve path parameters.
// Implementation should decide how to retrieve params based on values in specified http request.
type ParamsParser interface {
	ParseParams(r *http.Request) map[string]string
}

// The ParamsParserFunc type is an adapter to allow the use of ordinary
// functions as ParamsParser. If p is a function with
// the appropriate signature, ParamsParserFunc(p) is a ParamsParser that will return p
type ParamsParserFunc func(r *http.Request) map[string]string

// ParseParams returns p(r)
func (p ParamsParserFunc) ParseParams(r *http.Request) map[string]string {
	return p(r)
}

// HasParam returns true if key exists in path params map
func HasParam(r *http.Request, key string) bool {
	_, ok := Params(r)[key]
	return ok
}

// Params returns raw value for path param by key. If no key was set it returns "".
// To distinguish between empty value and value was not set use HasParam or Params directly
func Param(r *http.Request, key string) string {
	return Params(r)[key]
}

// Params retrieves path param values from request context or empty map if it was not set.
// Router is responsible to decorate http request with WithParams using ParamsParser
func Params(r *http.Request) map[string]string {
	params := r.Context().Value(paramsKey)
	if params == nil {
		return make(map[string]string)
	}
	return params.(map[string]string)
}

// WithParams returns a shallow copy of r with its context changed to contain params as context value
func WithParams(r *http.Request, params map[string]string) *http.Request {
	ctx := context.WithValue(r.Context(), paramsKey, params)
	return r.WithContext(ctx)
}
