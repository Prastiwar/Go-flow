package config_test

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"
	"unicode"

	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func setArgs(args ...string) {
	os.Args = []string{
		"go-flow", // os.Args[0] is always program name
	}
	os.Args = append(os.Args, args...)
}

func TestFlagProviderLoad(t *testing.T) {
	const desc = "test usage description"
	nowUtc := time.Now().UTC().Format(time.RFC3339)

	optionsWithLowerFirtCase := []config.LoadOption{
		config.WithInterceptor(func(name string, sf reflect.StructField) string {
			a := []rune(sf.Name)
			a[0] = unicode.ToLower(a[0])
			return string(a)
		}),
	}

	tests := []struct {
		name    string
		flags   []flag.Flag
		init    func(t *testing.T) (any, func())
		options []config.LoadOption
		wantErr bool
	}{
		{
			name: "success-all-flags",
			flags: []flag.Flag{
				config.BoolFlag("varBool", desc),
				config.StringFlag("varString", desc),
				config.Int32Flag("varInt32", desc),
				config.Int64Flag("varInt64", desc),
				config.Uint32Flag("varUint32", desc),
				config.Uint64Flag("varUint64", desc),
				config.Float32Flag("varFloat32", desc),
				config.Float64Flag("varFloat64", desc),
				config.DurationFlag("varDuration", desc),
				config.TimeFlag("varTime", desc),
			},
			init: func(t *testing.T) (any, func()) {
				setArgs(
					"-varBool=true",
					"-varString=text",
					"-varInt32=32",
					"-varInt64=64",
					"-varUint32=32",
					"-varUint64=64",
					"-varFloat32=32.32",
					"-varFloat64=64.64",
					"-varDuration=1s",
					"-varTime="+nowUtc,
				)

				v := struct {
					VarBool     bool
					VarString   string
					VarInt32    int32
					VarInt64    int64
					VarUint32   uint32
					VarUint64   uint64
					VarFloat32  float32
					VarFloat64  float64
					VarDuration time.Duration
					VarTime     time.Time
				}{}

				return &v, func() {
					assert.Equal(t, true, v.VarBool, "bool flag expectation failed")
					assert.Equal(t, "text", v.VarString, "string flag expectation failed")
					assert.Equal(t, int32(32), v.VarInt32, "int32 flag expectation failed")
					assert.Equal(t, int64(64), v.VarInt64, "int64 flag expectation failed")
					assert.Equal(t, uint32(32), v.VarUint32, "uint32 flag expectation failed")
					assert.Equal(t, uint64(64), v.VarUint64, "uint64 flag expectation failed")
					assert.Equal(t, float32(32.32), v.VarFloat32, "float32 flag expectation failed")
					assert.Equal(t, float64(64.64), v.VarFloat64, "float64 flag expectation failed")
					assert.Equal(t, 1*time.Second, v.VarDuration, "duration flag expectation failed")
					assert.Equal(t, nowUtc, v.VarTime.Format(time.RFC3339), "time flag expectation failed")
				}
			},
			options: optionsWithLowerFirtCase,
			wantErr: false,
		},
		{
			name: "success-pointers",
			flags: []flag.Flag{
				config.StringFlag("varString", desc),
				config.BoolFlag("varBool", desc),
			},
			init: func(t *testing.T) (any, func()) {
				setArgs(
					"-varBool",
					"-varString=pointing",
				)

				v := struct {
					VarBool   *bool
					VarString *string
				}{}

				return &v, func() {
					assert.Equal(t, "pointing", *v.VarString, "string pointer expectation failed")
					assert.Equal(t, true, *v.VarBool, "bool pointer expectation failed")
				}
			},
			options: optionsWithLowerFirtCase,
			wantErr: false,
		},
		{
			name: "success-convertible",
			flags: []flag.Flag{
				config.Int32Flag("varIntPointer", desc),
				config.Int32Flag("varInt", desc),
			},
			init: func(t *testing.T) (any, func()) {
				setArgs(
					"-varInt=64",
					"-varIntPointer=64",
				)

				v := struct {
					VarInt        int64
					VarIntPointer *int64
				}{}

				return &v, func() {
					assert.Equal(t, int64(64), v.VarInt, "convertible integer expectation failed")
					assert.Equal(t, int64(64), *v.VarIntPointer, "convertible integer pointer expectation failed")
				}
			},
			options: optionsWithLowerFirtCase,
			wantErr: false,
		},
		{
			name: "success-overriding",
			flags: []flag.Flag{
				config.StringFlag("notOverrideString", desc),
				config.StringFlag("overriddenEmptyString", desc),
			},
			init: func(t *testing.T) (any, func()) {
				setArgs(
					"-overriddenEmptyString=overridden",
				)

				v := struct {
					NotOverrideString     string
					OverriddenEmptyString string
				}{}

				v.NotOverrideString = "not-overridden"
				v.OverriddenEmptyString = "not-overridden"

				return &v, func() {
					assert.Equal(t, "not-overridden", v.NotOverrideString, "not override expectation failed")
					assert.Equal(t, "overridden", v.OverriddenEmptyString, "override expectation failed")
				}
			},
			options: optionsWithLowerFirtCase,
			wantErr: false,
		},
		{
			name:  "success-no-flag",
			flags: []flag.Flag{},
			init: func(t *testing.T) (any, func()) {
				setArgs()
				return &struct{}{}, func() {}
			},
			wantErr: false,
		},
		{
			name:  "invalid-provided-not-defined",
			flags: []flag.Flag{},
			init: func(t *testing.T) (any, func()) {
				setArgs("-boolean")
				return &struct{}{}, func() {}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := config.NewFlagProvider(tt.flags...)
			v, asserts := tt.init(t)

			err := p.Load(context.Background(), v, tt.options...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NilError(t, err)
			}

			if t.Failed() {
				t.FailNow()
			}

			asserts()
		})
	}
}

func TestNewFlagProvider(t *testing.T) {
	valid := config.NewFlagProvider(config.BoolFlag("name", "usage"))
	assert.NotNil(t, valid)

	defer func() {
		assert.ErrorIs(t, recover().(error), config.ErrMustImplementGetter)
	}()

	_ = config.NewFlagProvider(flag.Flag{})
}
