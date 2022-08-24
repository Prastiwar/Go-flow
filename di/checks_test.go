package di

import (
	"fmt"
	"goflow/reflection"
	"goflow/tests/assert"
	"reflect"
	"testing"
)

func TestCheckInterface(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		services map[reflect.Type]constructor
		expected bool
	}{
		{
			name: "success-found",
			typ:  reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]constructor{
				reflection.TypeOf[fmt.Stringer](): *Construct(Singleton, func() string { return "" }),
			},
			expected: true,
		},
		{
			name:     "invalid-not-found",
			typ:      reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]constructor{},
			expected: false,
		},
		{
			name: "invalid-not-interface",
			typ:  reflect.TypeOf(""),
			services: map[reflect.Type]constructor{
				reflect.TypeOf(""): *Construct(Singleton, func() string { return "" }),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := checkInterface(tt.typ, tt.services)

			assert.Equal(t, tt.expected, ok)
		})
	}
}

func TestCheckRegistered(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		services map[reflect.Type]constructor
		expected bool
	}{
		{
			name: "success-found-service",
			typ:  reflection.TypeOf[someService](),
			services: map[reflect.Type]constructor{
				reflection.TypeOf[someService](): *Construct(Singleton, newSomeService),
			},
			expected: true,
		},
		{
			name: "success-found-interface",
			typ:  reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]constructor{
				reflection.TypeOf[fmt.Stringer](): *Construct(Singleton, func() string { return "" }),
			},
			expected: true,
		},
		{
			name:     "invalid-not-found",
			typ:      reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]constructor{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := checkRegistered(tt.typ, tt.services)
			assert.Equal(t, tt.expected, ok)
		})
	}
}
