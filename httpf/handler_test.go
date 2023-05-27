package httpf

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/exception"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestRecoverMiddleware(t *testing.T) {
	t.Run("with-panic", func(t *testing.T) {
		defer exception.HandlePanicError(func(err error) {
			assert.NilError(t, err)
		})

		got := RecoverMiddleware(HandlerFunc(func(w ResponseWriter, r *http.Request) error {
			panic("panic-message")
		}))

		err := got.ServeHTTP(nil, nil)

		assert.ErrorWith(t, err, "panic-message")
	})

	t.Run("with-no-panic", func(t *testing.T) {
		defer exception.HandlePanicError(func(err error) {
			assert.NilError(t, err)
		})

		got := RecoverMiddleware(HandlerFunc(func(w ResponseWriter, r *http.Request) error {
			return errors.New("no-panic-error")
		}))

		err := got.ServeHTTP(nil, nil)

		assert.ErrorWith(t, err, "no-panic-error")
	})
}
