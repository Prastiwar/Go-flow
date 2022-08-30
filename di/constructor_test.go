package di

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestConstructorCreate(t *testing.T) {
	tests := []struct {
		name     string
		c        *constructor
		provider func(reflect.Type) interface{}
		want     interface{}
		panic    bool
	}{
		{
			name: "success-no-dep",
			c: Construct(Singleton, func() string {
				return ""
			}),
			provider: func(t reflect.Type) interface{} {
				return ""
			},
			want: "",
		},
		{
			name: "success-dep",
			c: Construct(Singleton, func(v int) string {
				return ""
			}),
			provider: func(t reflect.Type) interface{} {
				return 1
			},
			want: "",
		},
		{
			name: "panic-dep",
			c: Construct(Singleton, func(v int) string {
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
		life      LifeTime
		ctor      any
		assertErr assert.ErrorFunc
	}{
		{
			name: "success-no-dep",
			life: Singleton,
			ctor: func() string { return "" },
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "success-dep",
			life: Singleton,
			ctor: func(int) string { return "" },
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "invalid-instance",
			life: Singleton,
			ctor: "",
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, ErrCtorNotFunc, err)
			},
		},
		{
			name: "invalid-no-return",
			life: Singleton,
			ctor: func() {},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, ErrWrongCtorSignature, err)
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

			_ = Construct(tt.life, tt.ctor)
		})
	}
}
