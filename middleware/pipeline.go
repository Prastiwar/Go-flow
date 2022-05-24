package middleware

type Middleware[Request any, Response any] interface {
	Wrap(func(r Request) Response) func(r Request) Response

	Use(pipe func(r Request, next func(r Request) Response) Response)
}

type pipeFunc[Request any, Response any] func(func(Request) Response) func(Request) Response

type middlewareDesigner[Request any, Response any] struct {
	pipes []pipeFunc[Request, Response]
}

func NewMiddleware[Request any, Response any]() Middleware[Request, Response] {
	return &middlewareDesigner[Request, Response]{}
}

func (f *middlewareDesigner[Request, Response]) Wrap(fn func(r Request) Response) func(r Request) Response {
	// build pipeline in FIFO order
	for i := len(f.pipes) - 1; i >= 0; i-- {
		fn = f.pipes[i](fn)
	}
	return fn
}

// Use appends pipe function to middleware pipeline
func (f *middlewareDesigner[Request, Response]) Use(pipe func(r Request, next func(r Request) Response) Response) {
	wrapper := func(f func(Request) Response) func(Request) Response {
		return func(t Request) Response {
			return pipe(t, f)
		}
	}

	f.pipes = append(f.pipes, wrapper)
}
