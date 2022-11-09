package httpf

import (
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestParams(t *testing.T) {
	req := &http.Request{}
	has := HasParam(req, "one")
	assert.Equal(t, false, has, "expected to have 'one' key")

	r := WithParams(req, map[string]string{"one": "two"})

	has = HasParam(r, "one")
	assert.Equal(t, true, has, "expected to have 'one' key")

	v := Param(r, "one")
	assert.Equal(t, "two", v)
}
