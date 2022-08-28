// Package middleware implements middleware architecture pattern to use in generic way.
// It's used to create delegates for processing the request or response and handle common tasks like
// logging, authentication, compressing data in single contact point which is called pipeline.
package middleware

// Middleware is implemented by any value that has a Wrap and Use method.
// The implementation provides a way to build pipeline over specifiec function
// and store pipes using simple Use method.
type Middleware[Request any, Response any] interface {
	Wrap(func(r Request) Response) func(r Request) Response

	Use(pipe func(r Request, next func(r Request) Response) Response)
}

type pipeFunc[Request any, Response any] func(func(Request) Response) func(Request) Response

type middlewareDesigner[Request any, Response any] struct {
	pipes []pipeFunc[Request, Response]
}

// NewMiddleware returns a new generic middleware instance.
func NewMiddleware[Request any, Response any]() Middleware[Request, Response] {
	return &middlewareDesigner[Request, Response]{}
}

// Wrap builds pipeline over specified functiom. Pipeline will be called in FIFO order.
func (f *middlewareDesigner[Request, Response]) Wrap(fn func(r Request) Response) func(r Request) Response {
	for i := len(f.pipes) - 1; i >= 0; i-- {
		fn = f.pipes[i](fn)
	}
	return fn
}

// Use appends pipe function to middleware pipeline at the end.
func (f *middlewareDesigner[Request, Response]) Use(pipe func(r Request, next func(r Request) Response) Response) {
	wrapper := func(f func(Request) Response) func(Request) Response {
		return func(t Request) Response {
			return pipe(t, f)
		}
	}

	f.pipes = append(f.pipes, wrapper)
}
