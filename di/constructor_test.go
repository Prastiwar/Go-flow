package di

import (
	"goflow/tests/assert"
	"reflect"
	"testing"
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
		c         *constructor
		assertErr assert.ErrorFunc
	}{
		{
			name: "success-no-dep",
			c: Construct(Singleton, func() string {
				return ""
			}),
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "success-dep",
			c: Construct(Singleton, func(int) string {
				return ""
			}),
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name: "invalid-instance",
			c:    Construct(Singleton, ""),
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, NotFuncError, err)
			},
		},
		{
			name: "invalid-no-return",
			c:    Construct(Singleton, func() {}),
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, WrongCtorSignature, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			tt.assertErr(t, err)
		})
	}
}
