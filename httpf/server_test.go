package httpf_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestNewServer(t *testing.T) {
	got := httpf.NewServer("localhost:8086", nil)

	assert.Equal(t, reflect.TypeOf(&http.Server{}), reflect.TypeOf(got))
}
