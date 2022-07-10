package v2

import (
	"goflow/reflection"
	"goflow/tests/assert"
	"reflect"
	"testing"
)

type SomeInterface interface{}

type SomeDependency struct{}

type SomeOtherDependency struct{}

type SomeService struct{}

func NewSomeDependency() *SomeDependency {
	return &SomeDependency{}
}

func NewSomeService() *SomeService {
	return &SomeService{}
}

func NewSomeServiceWithDep(dep SomeDependency) *SomeService {
	return NewSomeService()
}

func NewSomeServiceWithTwoDeps(dep SomeDependency, otherDep SomeOtherDependency) *SomeService {
	return NewSomeService()
}

func NewSomeServiceWithCyclicDep(dep SomeService) *SomeService {
	return NewSomeService()
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name             string
		ctors            []interface{}
		expectedServices map[reflect.Type]constructor
		expectedErr      string
	}{
		{
			name:  "valid-some-service",
			ctors: []interface{}{NewSomeService},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*SomeService](): {
					fn: NewSomeService,
				},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := Register(tt.ctors...)
			var services map[reflect.Type]constructor
			if container != nil {
				services = container.services
			}

			assert.Equal(t, tt.expectedErr, getString(err))
			assert.Equal(t, len(tt.expectedServices), len(services))
			for typ, ctor := range tt.expectedServices {
				actualCtor, exists := services[typ]
				assert.Equal(t, true, exists)
				assert.Equal(t, reflect.TypeOf(ctor.fn), reflect.TypeOf(actualCtor.fn))
				assert.Equal(t, len(ctor.params), len(actualCtor.params))
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		ctors       []any
		expectedErr string
	}{
		{
			name:        "valid-no-deps",
			ctors:       []any{NewSomeService},
			expectedErr: "",
		},
		{
			name:        "valid-single-dep",
			ctors:       []any{NewSomeServiceWithDep, NewSomeDependency},
			expectedErr: "",
		},
		{
			name:        "missing-single-dep",
			ctors:       []any{NewSomeServiceWithDep},
			expectedErr: "['dependency is not registered': 'v2.SomeDependency']",
		},
		{
			name:        "missing-two-deps",
			ctors:       []any{NewSomeServiceWithTwoDeps},
			expectedErr: "['dependency is not registered': 'v2.SomeDependency', 'dependency is not registered': 'v2.SomeOtherDependency']",
		},
		{
			name:        "cyclic-dependency",
			ctors:       []any{NewSomeServiceWithCyclicDep},
			expectedErr: "['cyclic dependency detected': 'v2.SomeService']",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := Register(tt.ctors...)
			assert.NilError(t, err)

			err = container.Validate()

			assert.Equal(t, tt.expectedErr, getString(err))
		})
	}
}

func TestProvideSuccess(t *testing.T) {
	defer func() {
		assert.Equal(t, nil, recover())
	}()

	container, err := Register(NewSomeService)
	assert.NilError(t, err)

	tests := []struct {
		name         string
		provideFn    func(t *testing.T, provider func(any)) any
		expectedType reflect.Type
	}{
		{
			name: "valid-interface-no-deps",
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeInterface
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
		{
			name: "valid-service-no-deps",
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[SomeService](),
		},
		{
			name: "valid-service-pointer-no-deps",
			provideFn: func(t *testing.T, provider func(any)) any {
				var service *SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
		{
			name: "invalid-service-pointer-no-deps",
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, NotAddresableError, recover())
					t.SkipNow()
				}()

				var service *SomeService
				provider(service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.provideFn(t, container.Provide)

			assert.NotNil(t, service)
			assert.Equal(t, tt.expectedType, reflect.TypeOf(service))
		})
	}
}

func getString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
