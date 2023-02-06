package config_test

import (
	"errors"
	"testing"

	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestSetFields(t *testing.T) {
	tests := []struct {
		name            string
		assertWithValue func(t *testing.T) (any, func(err error))
		opts            config.LoadOptions
		findFn          config.FieldValueFinder
	}{
		{
			name: "success-raw-value-pointer",
			assertWithValue: func(t *testing.T) (any, func(err error)) {
				v := struct {
					Field string
				}{}
				return &v, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, "str", v.Field)
				}
			},
			opts: *config.NewLoadOptions(),
			findFn: func(key string) (any, error) {
				field := "str"
				return &field, nil
			},
		},
		{
			name: "failure-find-fn-error",
			assertWithValue: func(t *testing.T) (any, func(err error)) {
				v := struct {
					Field string
				}{}
				return &v, func(err error) {
					assert.ErrorWith(t, err, "find-error")
					assert.Equal(t, "", v.Field)
				}
			},
			opts: *config.NewLoadOptions(),
			findFn: func(key string) (any, error) {
				return nil, errors.New("find-error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, asserts := tt.assertWithValue(t)
			setter := config.NewFieldSetter("", tt.opts)

			err := setter.SetFields(v, tt.findFn)

			asserts(err)
		})
	}
}
