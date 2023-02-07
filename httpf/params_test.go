package httpf_test

import (
	"net/http"
	"testing"

	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestParams(t *testing.T) {
	req := &http.Request{}
	has := httpf.HasParam(req, "one")
	assert.Equal(t, false, has, "expected to have 'one' key")

	r := httpf.WithParams(req, map[string]string{"one": "two"})

	has = httpf.HasParam(r, "one")
	assert.Equal(t, true, has, "expected to have 'one' key")

	v := httpf.Param(r, "one")
	assert.Equal(t, "two", v)
}
