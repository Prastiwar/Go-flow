package httpf

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

func TestServeMux(t *testing.T) {
	mux := NewServeMuxBuilder()

	handlerFunc := func(name string) Handler {
		counter := assert.Count(t, 1, name+" route was not called")
		return HandlerFunc(func(w ResponseWriter, r *http.Request) error {
			counter.Inc()
			return nil
		})
	}

	router := mux.Get("/api/albums/", handlerFunc("Get")).
		Post("/api/albums/", handlerFunc("Post")).
		Put("/api/albums/", handlerFunc("Put")).
		Delete("/api/albums/", handlerFunc("Delete")).
		Options("/api/albums/", handlerFunc("Options")).
		Patch("/api/albums/", handlerFunc("Patch")).
		Build()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
	}

	for _, method := range methods {
		r, err := http.NewRequest(method, "http://localhost/api/albums/", nil)
		assert.NilError(t, err)

		router.ServeHTTP(&mocks.ResponseWriter{}, r)
	}
}

func TestErrorHandler(t *testing.T) {
	mux := NewServeMuxBuilder()
	errHandler := errors.New("handler-error")

	errorCounter := assert.Count(t, 1, "expected call to error handler")

	router := mux.WithErrorHandler(ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		errorCounter.Inc()
		assert.Equal(t, errHandler, err)
	})).Get("/api/albums/", HandlerFunc(func(w ResponseWriter, r *http.Request) error {
		return errHandler
	})).Build()

	r, err := http.NewRequest(http.MethodGet, "http://localhost/api/albums/", nil)
	assert.NilError(t, err)

	router.ServeHTTP(&mocks.ResponseWriter{}, r)
}

func TestServeMuxHandleGetError(t *testing.T) {
	errHandle := errors.New("handler-error")

	tests := []struct {
		name          string
		requestMethod string
		errHandler    func(t *testing.T) ErrorHandler
		writer        func(t *testing.T) http.ResponseWriter
	}{
		{
			name:          "with-custom-handler",
			requestMethod: http.MethodGet,
			errHandler: func(t *testing.T) ErrorHandler {
				errorCounter := assert.Count(t, 1, "expected call to error handler")
				return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {
					errorCounter.Inc()
					assert.Equal(t, errHandle, err)
				})
			},
			writer: func(t *testing.T) http.ResponseWriter {
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						t.Fatal("unexpected Header() function call")
						return nil
					},
					OnWrite: func(b []byte) (int, error) {
						t.Fatal("unexpected Header() function call")
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						t.Fatal("unexpected Header() function call")
					},
				}
				return m
			},
		},
		{
			name:          "with-default-handler",
			requestMethod: http.MethodGet,
			errHandler: func(t *testing.T) ErrorHandler {
				return nil
			},
			writer: func(t *testing.T) http.ResponseWriter {
				writeCounter := assert.Count(t, 1, "expected Write() call")
				writeHeaderCounter := assert.Count(t, 1, "expected WriteHeader() call")
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						return http.Header{}
					},
					OnWrite: func(b []byte) (int, error) {
						assert.Equal(t, errHandle.Error()+"\n", string(b))
						writeCounter.Inc()
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						assert.Equal(t, 500, code)
						writeHeaderCounter.Inc()
					},
				}
				return m
			},
		},
		{
			name:          "with-method-not-allowed",
			requestMethod: http.MethodPost,
			errHandler: func(t *testing.T) ErrorHandler {
				return nil
			},
			writer: func(t *testing.T) http.ResponseWriter {
				writeCounter := assert.Count(t, 1, "expected Write() call")
				writeHeaderCounter := assert.Count(t, 1, "expected WriteHeader() call")
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						return http.Header{}
					},
					OnWrite: func(b []byte) (int, error) {
						assert.Equal(t, "\n", string(b))
						writeCounter.Inc()
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						assert.Equal(t, http.StatusMethodNotAllowed, code)
						writeHeaderCounter.Inc()
					},
				}
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewServeMuxBuilder()

			router := mux.WithErrorHandler(tt.errHandler(t)).
				Get("/api/albums/", HandlerFunc(func(w ResponseWriter, r *http.Request) error {
					return errHandle
				})).Build()

			r, err := http.NewRequest(tt.requestMethod, "http://localhost/api/albums/", nil)
			assert.NilError(t, err)

			router.ServeHTTP(tt.writer(t), r)
		})
	}
}

func TestWriterDecoration(t *testing.T) {
	tests := []struct {
		name      string
		code      int
		data      interface{}
		decorator func(t *testing.T) func(http.ResponseWriter) ResponseWriter
	}{
		{
			name: "succes-default-writer",
			code: 201,
			data: struct {
				Id string `json:"id"`
			}{Id: "123"},
			decorator: func(t *testing.T) func(http.ResponseWriter) ResponseWriter {
				return nil
			},
		},
		{
			name: "succes-custom-writer",
			code: 204,
			data: nil,
			decorator: func(t *testing.T) func(http.ResponseWriter) ResponseWriter {
				callCounter := assert.Count(t, 1)
				return func(w http.ResponseWriter) ResponseWriter {
					return &mocks.ResponseWriter{
						OnResponse: func(code int, data interface{}) error {
							callCounter.Inc()
							return nil
						},
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewServeMuxBuilder()

			router := mux.WithWriterDecorator(tt.decorator(t)).
				Get("/api/albums/", HandlerFunc(func(w ResponseWriter, r *http.Request) error {
					return w.Response(tt.code, tt.data)
				})).Build()

			r, err := http.NewRequest(http.MethodGet, "http://localhost/api/albums/", nil)
			assert.NilError(t, err)

			router.ServeHTTP(&mocks.ResponseWriter{
				OnHeader:      func() http.Header { return http.Header{} },
				OnWriteHeader: func(statusCode int) {},
				OnWrite:       func(b []byte) (int, error) { return 0, err },
			}, r)
		})
	}
}
