package reflection_test

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/reflection"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type Foo struct{}

func TestTypeOf(t *testing.T) {
	stringType := reflection.TypeOf[string]()
	assert.Equal(t, "string", stringType.Name())

	fooType := reflection.TypeOf[Foo]()
	assert.Equal(t, "Foo", fooType.Name())
}

func TestInParamTypes(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want []reflect.Type
	}{
		{
			name: "success-zero",
			typ:  reflect.TypeOf(func() string { return "" }),
			want: []reflect.Type{},
		},
		{
			name: "success-single",
			typ:  reflect.TypeOf(func(int) string { return "" }),
			want: []reflect.Type{
				reflection.TypeOf[int](),
			},
		},
		{
			name: "success-two",
			typ:  reflect.TypeOf(func(int, string) string { return "" }),
			want: []reflect.Type{
				reflection.TypeOf[int](),
				reflection.TypeOf[string](),
			},
		},
		{
			name: "invalid-type",
			typ:  reflect.TypeOf("func(int, string) int { return 1 }"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflection.InParamTypes(tt.typ)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestOutParamTypes(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want []reflect.Type
	}{
		{
			name: "success-zero",
			typ:  reflect.TypeOf(func(string) {}),
			want: []reflect.Type{},
		},
		{
			name: "success-single",
			typ:  reflect.TypeOf(func() int { return 0 }),
			want: []reflect.Type{
				reflection.TypeOf[int](),
			},
		},
		{
			name: "success-two",
			typ:  reflect.TypeOf(func() (int, string) { return 0, "" }),
			want: []reflect.Type{
				reflection.TypeOf[int](),
				reflection.TypeOf[string](),
			},
		},
		{
			name: "invalid-type",
			typ:  reflect.TypeOf("func(int, string) int { return 1 }"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflection.OutParamTypes(tt.typ)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestTogglePointer(t *testing.T) {
	tests := []struct {
		name    string
		typ     reflect.Type
		pointer bool
	}{
		{
			name:    "success-pointer-to-nonpointer",
			typ:     reflection.TypeOf[*bool](),
			pointer: false,
		},
		{
			name:    "success-nonpointer-to-pointer",
			typ:     reflection.TypeOf[bool](),
			pointer: true,
		},
		{
			name:    "success-double-pointer-to-pointer",
			typ:     reflection.TypeOf[**any](),
			pointer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflection.TogglePointer(tt.typ)

			assert.NotNil(t, got)
			assert.Equal(t, tt.pointer, got.Kind() == reflect.Pointer)
		})
	}
}
