package di_test

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/di"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestConstructorCreate(t *testing.T) {
	tests := []struct {
		name     string
		c        di.Constructor
		provider func(reflect.Type) interface{}
		want     interface{}
		panic    bool
	}{
		{
			name: "success-no-dep",
			c: di.Construct(di.Singleton, func() string {
				return ""
			}),
			provider: func(t reflect.Type) interface{} {
				return ""
			},
			want: "",
		},
		{
			name: "success-dep",
			c: di.Construct(di.Singleton, func(v int) string {
				return ""
			}),
			provider: func(t reflect.Type) interface{} {
				return 1
			},
			want: "",
		},
		{
			name: "panic-dep",
			c: di.Construct(di.Singleton, func(v int) string {
				return ""
			}),
			provider: func(t reflect.Type) interface{} {
				return ""
			},
			panic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				defer func() {
					assert.NotNil(t, recover())
				}()
			}
			got := tt.c.Create(tt.provider)
			assert.Equal(t, tt.want, got, "create failed")
		})
	}
}

func TestConstructorValidate(t *testing.T) {
	tests := []struct {
		name      string
		life      di.LifeTime
		ctor      any
		assertErr assert.ErrorFunc
	}{
		{
			name: "success-no-dep",
			life: di.Singleton,
			ctor: func() string { return "" },
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "success-dep",
			life: di.Singleton,
			ctor: func(int) string { return "" },
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "invalid-instance",
			life: di.Singleton,
			ctor: "",
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, di.ErrCtorNotFunc, err)
			},
		},
		{
			name: "invalid-no-return",
			life: di.Singleton,
			ctor: func() {},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, di.ErrWrongCtorSignature, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				var err error
				if r := recover(); r != nil {
					rerr, ok := r.(error)
					assert.Equal(t, true, ok, "recover error cast failed")
					err = rerr
				}

				tt.assertErr(t, err)
			}()

			_ = di.Construct(tt.life, tt.ctor)
		})
	}
}

func TestConstructorFunc(t *testing.T) {
	expectedValue := 0
	f := di.ConstructorFunc(func(provider func(reflect.Type) interface{}) interface{} {
		return expectedValue
	})

	val := f.Create(nil)

	assert.Equal(t, expectedValue, val)
}
