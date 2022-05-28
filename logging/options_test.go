package logging

import (
	"goflow/tests/assert"
	"testing"
)

func TestBuildOptions(t *testing.T) {
	opts := []LogOption{
		WithLogFormat("log"),
		WithInfoFormat("info"),
		WithWarnFormat("warn"),
		WithErrorFormat("error"),
	}

	options := Options(opts...)

	assert.NotNil(t, options)
	assert.Equal(t, "log", options.logf)
	assert.Equal(t, "info", options.infof)
	assert.Equal(t, "warn", options.warnf)
	assert.Equal(t, "error", options.errorf)
}
