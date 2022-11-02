package httpf

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

func TestJson(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		data      interface{}
		writer    func(t *testing.T) http.ResponseWriter
		assertion assert.ErrorFunc
	}{
		{
			name:   "success",
			status: 201,
			data: struct {
				Id string `json:"id"`
			}{Id: "123"},
			writer: func(t *testing.T) http.ResponseWriter {
				writeCounter := assert.Count(t, 1)
				writeHeaderCounter := assert.Count(t, 1)
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						t.Fatal("unexpected Header() function call")
						return nil
					},
					OnWrite: func(b []byte) (int, error) {
						assert.Equal(t, `{"id":"123"}`, string(b))
						writeCounter.Inc()
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						assert.Equal(t, 201, code)
						writeHeaderCounter.Inc()
					},
				}
				return m
			},
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:   "success-nil-data",
			status: 204,
			data:   nil,
			writer: func(t *testing.T) http.ResponseWriter {
				writeHeaderCounter := assert.Count(t, 1)
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						t.Fatal("unexpected Header() function call")
						return nil
					},
					OnWrite: func(b []byte) (int, error) {
						t.Fatal("unexpected Write() function call")
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						assert.Equal(t, 204, code)
						writeHeaderCounter.Inc()
					},
				}
				return m
			},
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:   "invalid-pointer-data",
			status: 204,
			data:   &struct{ Chan chan (int) }{},
			writer: func(t *testing.T) http.ResponseWriter {
				m := &mocks.ResponseWriter{
					OnHeader: func() http.Header {
						t.Fatal("unexpected Header() function call")
						return nil
					},
					OnWrite: func(b []byte) (int, error) {
						t.Fatal("unexpected Write() function call")
						return 0, nil
					},
					OnWriteHeader: func(code int) {
						t.Fatal("unexpected Write() function call")
					},
				}
				return m
			},
			assertion: func(t *testing.T, err error) {
				assert.ErrorType(t, err, &json.UnsupportedTypeError{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Json(tt.writer(t), tt.status, tt.data)

			tt.assertion(t, err)
		})
	}
}
