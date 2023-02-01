package httpf

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestComposeRateKeyFactories(t *testing.T) {
	f := ComposeRateKeyFactories(IPRateKey(), PathRateKey())
	url, _ := url.Parse("https://test.com/api/resource")
	key := f(&http.Request{Method: http.MethodGet, URL: url})
	assert.Equal(t, "0.0.0.0 GET:/api/resource", key)

	defer func() {
		assert.NotNil(t, recover())
	}()
	_ = ComposeRateKeyFactories()
}
