package reflection

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

type Foo struct{}

func TestTypeOf(t *testing.T) {
	stringType := TypeOf[string]()
	assert.Equal(t, "string", stringType.Name())

	fooType := TypeOf[Foo]()
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
				TypeOf[int](),
			},
		},
		{
			name: "success-two",
			typ:  reflect.TypeOf(func(int, string) string { return "" }),
			want: []reflect.Type{
				TypeOf[int](),
				TypeOf[string](),
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
			got := InParamTypes(tt.typ)
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
				TypeOf[int](),
			},
		},
		{
			name: "success-two",
			typ:  reflect.TypeOf(func() (int, string) { return 0, "" }),
			want: []reflect.Type{
				TypeOf[int](),
				TypeOf[string](),
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
			got := OutParamTypes(tt.typ)
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
			typ:     TypeOf[*bool](),
			pointer: false,
		},
		{
			name:    "success-nonpointer-to-pointer",
			typ:     TypeOf[bool](),
			pointer: true,
		},
		{
			name:    "success-double-pointer-to-pointer",
			typ:     TypeOf[**any](),
			pointer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TogglePointer(tt.typ)

			assert.NotNil(t, got)
			assert.Equal(t, tt.pointer, got.Kind() == reflect.Pointer)
		})
	}
}
