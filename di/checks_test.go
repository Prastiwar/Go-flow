package di

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/reflection"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type someService struct{}

func NewSomeService() *someService {
	return &someService{}
}

func TestCheckInterface(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		services map[reflect.Type]Constructor
		expected bool
	}{
		{
			name: "success-found",
			typ:  reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]Constructor{
				reflection.TypeOf[fmt.Stringer](): Construct(Singleton, func() string { return "" }),
			},
			expected: true,
		},
		{
			name:     "invalid-not-found",
			typ:      reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]Constructor{},
			expected: false,
		},
		{
			name: "invalid-not-interface",
			typ:  reflect.TypeOf(""),
			services: map[reflect.Type]Constructor{
				reflect.TypeOf(""): Construct(Singleton, func() string { return "" }),
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
		services map[reflect.Type]Constructor
		expected bool
	}{
		{
			name: "success-found-service",
			typ:  reflection.TypeOf[someService](),
			services: map[reflect.Type]Constructor{
				reflection.TypeOf[someService](): Construct(Singleton, NewSomeService),
			},
			expected: true,
		},
		{
			name: "success-found-interface",
			typ:  reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]Constructor{
				reflection.TypeOf[fmt.Stringer](): Construct(Singleton, func() string { return "" }),
			},
			expected: true,
		},
		{
			name:     "invalid-not-found",
			typ:      reflection.TypeOf[fmt.Stringer](),
			services: map[reflect.Type]Constructor{},
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
