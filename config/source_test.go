package config

import (
	"encoding/json"
	"errors"
	"goflow/config/decoders"
	"goflow/tests/assert"
	"os"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

var (
	errCannotRead = errors.New("cannot read")
)

type InvalidReader struct{}

func (r *InvalidReader) Read(p []byte) (n int, err error) {
	return 0, errCannotRead
}

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

func TestSourceLoad(t *testing.T) {
	tests := []struct {
		name string
		init func(t *testing.T) (*Source, any, func(error))
	}{
		{
			name: "success-empty",
			init: func(t *testing.T) (*Source, any, func(error)) {
				s := Provide()
				v := struct {
					Key string
				}{}
				v.Key = "test"

				return s, &v, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, "test", v.Key)
				}
			},
		},
		{
			name: "success-complex",
			init: func(t *testing.T) (*Source, any, func(error)) {
				s := Provide(
					NewFlagProvider(
						StringFlag("flagKey", "just a string"),
						StringFlag("overriden", "just a string"),
						StringFlag("notOverriden", "just a string"),
					),
					NewEnvProvider(),
				)

				err := s.SetDefault(
					*Opt("DefaultKey", "1234567890"),
				)
				assert.NilError(t, err)

				s.ShareOptions(
					WithInterceptor(func(providerName string, field reflect.StructField) string {
						if providerName == EnvProviderName {
							return strings.ToUpper(field.Name)
						}

						a := []rune(field.Name)
						a[0] = unicode.ToLower(a[0])
						return string(a)
					}),
				)

				v := struct {
					DefaultKey   string
					EnvKey       string
					FlagKey      string
					Overriden    string
					NotOverriden string
				}{}
				v.Overriden = "to-override"
				v.NotOverriden = "not-overriden"

				assert.NilError(t, os.Setenv("ENVKEY", "envved"))
				assert.NilError(t, os.Setenv("OVERRIDEN", "overriden"))

				setArgs(
					"-flagKey=flagged",
					"-overriden=yes",
				)

				return s, &v, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, "1234567890", v.DefaultKey)
					assert.Equal(t, "envved", v.EnvKey)
					assert.Equal(t, "flagged", v.FlagKey)
					assert.Equal(t, "overriden", v.Overriden)
					assert.Equal(t, "not-overriden", v.NotOverriden)

					assert.NilError(t, os.Unsetenv("ENVKEY"))
					assert.NilError(t, os.Unsetenv("OVERRIDEN"))
				}
			},
		},
		{
			name: "invalid-non-pointer",
			init: func(t *testing.T) (*Source, any, func(error)) {
				s := Provide()
				v := struct{}{}

				return s, v, func(err error) {
					assert.ErrorIs(t, err, ErrNonPointer)
				}
			},
		},
		{
			name: "invalid-default",
			init: func(t *testing.T) (*Source, any, func(error)) {
				s := Provide()
				err := s.SetDefault(
					*Opt("Key", 10),
				)
				assert.NilError(t, err)

				v := struct {
					Key chan int
				}{}

				return s, &v, func(err error) {
					assert.ErrorType(t, err, &json.UnmarshalTypeError{})
				}
			},
		},
		{
			name: "invalid-provider-error",
			init: func(t *testing.T) (*Source, any, func(error)) {
				s := Provide(
					NewReaderProvider(&InvalidReader{}, decoders.NewJson()),
				)

				v := struct{}{}

				return s, &v, func(err error) {
					assert.ErrorIs(t, err, errCannotRead)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, v, asserts := tt.init(t)

			err := source.Load(v)

			asserts(err)
		})
	}
}

func TestBind(t *testing.T) {
	tests := []struct {
		name string
		init func(t *testing.T) (any, any, func(error))
	}{
		{
			name: "success",
			init: func(t *testing.T) (any, any, func(error)) {
				from := struct {
					Field string
					Int   int32
				}{}
				from.Field = "test"

				to := struct {
					Field string
					Int   int64
				}{}

				return &from, &to, func(err error) {
					assert.NilError(t, err)
					assert.Equal(t, from.Field, to.Field)
					assert.Equal(t, int64(from.Int), to.Int)
				}
			},
		},
		{
			name: "invalid-non-pointer",
			init: func(t *testing.T) (any, any, func(error)) {
				from := struct{}{}
				to := struct{}{}

				return from, to, func(err error) {
					assert.ErrorIs(t, err, ErrNonPointer)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to, asserts := tt.init(t)

			err := Bind(from, to)

			asserts(err)
		})
	}
}