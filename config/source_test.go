package config

import (
	"encoding/json"
	"errors"
	"goflow/tests/assert"
	"testing"
)

func TestSourceSetDefault(t *testing.T) {
	tests := []struct {
		name     string
		defaults []opt
		init     func(t *testing.T) (*Source, func(error))
	}{
		{
			name: "success",
			defaults: []opt{
				*Opt("key", "value"),
				*Opt("key2", 1),
			},
			init: func(t *testing.T) (*Source, func(error)) {
				s := Provide()
				return s, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, "{\"key\":\"value\",\"key2\":1}", string(s.defaults))
				}
			},
		},
		{
			name: "invalid-duplicate-key",
			defaults: []opt{
				*Opt("key", "value"),
				*Opt("key", 1),
			},
			init: func(t *testing.T) (*Source, func(error)) {
				s := Provide()
				return s, func(err error) {
					isExpectedError := errors.Is(err, ErrDuplicateKey)
					assert.Equal(t, true, isExpectedError, "error expectation failed")
				}
			},
		},
		{
			name: "invalid-json-value",
			defaults: []opt{
				*Opt("key", make(chan int)),
			},
			init: func(t *testing.T) (*Source, func(error)) {
				s := Provide()
				return s, func(err error) {
					_, isExpectedError := err.(*json.UnsupportedTypeError)
					assert.Equal(t, true, isExpectedError, "error expectation failed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, asserts := tt.init(t)

			err := source.SetDefault(tt.defaults...)

			asserts(err)
		})
	}
}

func TestSourceDefault(t *testing.T) {
	tests := []struct {
		name string
		init func(t *testing.T) (*Source, any, func(err error))
	}{
		{
			name: "success-empty-default",

			init: func(t *testing.T) (*Source, any, func(err error)) {
				s := Provide()
				v := struct{}{}
				return s, &v, func(err error) {
					assert.NilError(t, err)
				}
			},
		},
		{
			name: "success-not-empty-default",

			init: func(t *testing.T) (*Source, any, func(err error)) {
				s := Provide()
				err := s.SetDefault(
					*Opt("key", "value"),
				)
				assert.NilError(t, err)

				v := struct {
					Key            string `json:"key"`
					NotOverrideKey string
				}{}
				v.NotOverrideKey = "not-overriden"

				return s, &v, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, "value", v.Key)
					assert.Equal(t, "not-overriden", v.NotOverrideKey)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, v, asserts := tt.init(t)

			err := source.Default(v)

			asserts(err)
		})
	}
}
