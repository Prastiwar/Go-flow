package di

import (
	"goflow/reflection"
	"goflow/tests/assert"
	"math/rand"
	"reflect"
	"testing"
)

type SomeInterface interface{}

type SomeDependency struct {
	id float64
}

type SomeOtherDependency struct{}

type SomeService struct {
	id float64
}

func NewSomeDependency() *SomeDependency {
	return &SomeDependency{
		id: rand.Float64(),
	}
}

func NewSomeDependencyNoPointer() SomeDependency {
	return SomeDependency{}
}

func NewSomeService() *SomeService {
	return &SomeService{
		id: rand.Float64(),
	}
}

func NewSomeServiceNonPointerDep(*SomeDependency) SomeService {
	return SomeService{
		id: rand.Float64(),
	}
}

func NewSomeServiceWithDep(dep SomeDependency) *SomeService {
	return NewSomeService()
}

func NewSomeServiceWithPointerDep(dep *SomeDependency) *SomeService {
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
		assertErr        assert.ErrorFunc
	}{
		{
			name:  "success-some-service",
			ctors: []interface{}{NewSomeService},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*SomeService](): {
					fn: NewSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-pointer",
			ctors: []interface{}{Construct(Scoped, NewSomeService)},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*SomeService](): {
					fn: NewSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-non-pointer",
			ctors: []interface{}{*Construct(Scoped, NewSomeService)},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*SomeService](): {
					fn: NewSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-ctor-func",
			ctors: []interface{}{NewSomeService, ""},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, NotFuncError, err)
			},
		},
		{
			name:  "invalid-ctor-signature",
			ctors: []interface{}{NewSomeService, func() {}},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, WrongCtorSignature, err)
			},
		},
		// Validation
		{
			name:  "success-no-pointer-single-dep",
			ctors: []any{NewSomeServiceNonPointerDep, NewSomeDependency},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*SomeDependency](): {
					fn: NewSomeDependency,
				},
				reflection.TypeOf[SomeService](): {
					fn:     NewSomeServiceNonPointerDep,
					params: []reflect.Type{reflect.TypeOf(NewSomeDependency())},
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-missing-single-dep",
			ctors: []any{NewSomeServiceWithDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "'dependency is not registered': 'di.SomeDependency'")
			},
		},
		{
			name:  "invalid-missing-two-deps",
			ctors: []any{NewSomeServiceWithTwoDeps},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "dependency is not registered': 'di.SomeDependency")
				assert.ErrorWith(t, err, "dependency is not registered': 'di.SomeOtherDependency")
			},
		},
		{
			name:  "invalid-cyclic-dependency",
			ctors: []any{NewSomeServiceWithCyclicDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "cyclic dependency detected': 'di.SomeService")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := Register(tt.ctors...)
			var services map[reflect.Type]constructor
			if container != nil {
				services = container.services
			}

			tt.assertErr(t, err)
			if t.Failed() {
				t.FailNow()
			}

			assert.Equal(t, len(tt.expectedServices), len(services))
			for typ, ctor := range tt.expectedServices {
				actualCtor, exists := services[typ]
				assert.Equal(t, true, exists, "exists assertion failed")
				assert.Equal(t, reflect.TypeOf(ctor.fn), reflect.TypeOf(actualCtor.fn), "type equality failed")
				assert.Equal(t, len(ctor.params), len(actualCtor.params), "ctor params equality failed", reflect.TypeOf(ctor.fn).String())
			}
		})
	}
}

func TestProvide(t *testing.T) {
	defer func() {
		assert.Equal(t, nil, recover())
	}()

	tests := []struct {
		name         string
		container    func(t *testing.T) (*container, error)
		provideFn    func(t *testing.T, provider func(any)) any
		expectedType reflect.Type
	}{
		{
			name: "success-interface-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeInterface
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
		{
			name: "success-service-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[SomeService](),
		},
		{
			name: "success-service-pointer-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service *SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
		{
			name: "invalid-service-pointer-to-pointer-dep",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeServiceWithDep, NewSomeDependency)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, NotPointerError, recover())
					t.SkipNow()
				}()

				var service *SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
		{
			name: "invalid-service-pointer-no-addresable-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeService)
			},
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
		{
			name: "invalid-service-not-registered-no-deps",
			container: func(t *testing.T) (*container, error) {
				return &container{}, nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					err, _ := recover().(error)
					assert.ErrorWith(t, err, NotRegisteredError.Error())
					t.SkipNow()
				}()

				var service SomeService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*SomeService](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := tt.container(t)
			assert.NilError(t, err)

			service := tt.provideFn(t, container.Provide)

			assert.NotNil(t, service)
			assert.Equal(t, tt.expectedType, reflect.TypeOf(service))
		})
	}
}

func TestProvideCache(t *testing.T) {
	defer func() {
		assert.Equal(t, nil, recover())
	}()

	tests := []struct {
		name      string
		container func(t *testing.T) (*container, error)
		provideFn func(t *testing.T, provider func(any)) any
		cached    bool
	}{
		{
			name: "success-transient-root-no-cached",
			container: func(t *testing.T) (*container, error) {
				return Register(NewSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-singleton-root-cached",
			container: func(t *testing.T) (*container, error) {
				return Register(Construct(Singleton, NewSomeService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-scoped-root-no-cached",
			container: func(t *testing.T) (*container, error) {
				return Register(Construct(Scoped, NewSomeService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-scoped-scope-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Scoped, NewSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-singleton-scope-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Singleton, NewSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-transient-scope-no-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Transient, NewSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service SomeService
				provider(&service)
				return &service
			},
			cached: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := tt.container(t)
			assert.NilError(t, err)

			service := tt.provideFn(t, container.Provide)
			service2 := tt.provideFn(t, container.Provide)

			if tt.cached {
				assert.Equal(t, service, service2, "service was not cached")
			} else {
				assert.NotEqual(t, service, service2, "service was cached")
			}
		})
	}
}
