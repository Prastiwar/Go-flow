package httpf_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/datas"
	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type nameStructFixture struct {
	Name string `json:"name"`
}

type httpErrorFixture struct {
	Message string `json:"message"`
}

func (err *httpErrorFixture) Error() string {
	return err.Message
}

func TestBodyUnmarshaler(t *testing.T) {
	tests := []struct {
		name      string
		u         httpf.BodyUnmarshaler
		r         http.Response
		v         any
		assertion assert.ResultErrorFunc[any]
	}{
		{
			name: "handle-200-status-with-json",
			u: httpf.NewBodyUnmarshaler(datas.Json(), func(r *http.Response, u datas.ReaderUnmarshaler) error {
				panic("unexpected call")
			}),
			r: http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"name":"foo"}`)),
			},
			v: &nameStructFixture{},
			assertion: func(t *testing.T, result any, err error) {
				assert.NilError(t, err)
				v, ok := result.(*nameStructFixture)
				assert.Equal(t, true, ok)
				assert.Equal(t, "foo", v.Name)
			},
		},
		{
			name: "handle-404-status-custom-error",
			u: httpf.NewBodyUnmarshaler(datas.Json(), func(r *http.Response, u datas.ReaderUnmarshaler) error {
				if r.StatusCode == http.StatusNotFound {
					return errors.New("not-found")
				}
				return errors.New("invalid-status")
			}),
			r: http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(nil),
			},
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorWith(t, err, "not-found")
			},
		},
		{
			name: "error-404-status-with-struct",
			u:    httpf.NewBodyUnmarshalerWithError(datas.Json(), &httpErrorFixture{}),
			r: http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"resource-not-found"}`)),
			},
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorType(t, err, &httpErrorFixture{})
				assert.Equal(t, "resource-not-found", err.Error())
			},
		},
		{
			name: "error-400-status-with-invalid-json-body",
			u:    httpf.NewBodyUnmarshalerWithError(datas.Json(), &httpErrorFixture{}),
			r: http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`{`)),
			},
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.u.Unmarshal(&tt.r, tt.v)
			tt.assertion(t, tt.v, err)
		})
	}
}
