package di

import (
	"goflow/reflection"
	"goflow/tests/assert"
	"math/rand"
	"reflect"
	"testing"
)

type someInterface interface{}

type someDependency struct {
	id float64
}

type someOtherDependency struct{}

type someService struct {
	id float64
}

func newSomeDependency() *someDependency {
	return &someDependency{
		id: rand.Float64(),
	}
}

func newSomeService() *someService {
	return &someService{
		id: rand.Float64(),
	}
}

func newSomeServiceNonPointerDep(*someDependency) someService {
	return someService{
		id: rand.Float64(),
	}
}

func newSomeServiceWithDep(dep someDependency) *someService {
	return newSomeService()
}

func newSomeServiceWithTwoDeps(dep someDependency, otherDep someOtherDependency) *someService {
	return newSomeService()
}

func newSomeServiceWithCyclicDep(dep someService) *someService {
	return newSomeService()
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
			ctors: []interface{}{newSomeService},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*someService](): {
					fn: newSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-pointer",
			ctors: []interface{}{Construct(Scoped, newSomeService)},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*someService](): {
					fn: newSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-non-pointer",
			ctors: []interface{}{*Construct(Scoped, newSomeService)},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*someService](): {
					fn: newSomeService,
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-ctor-func",
			ctors: []interface{}{newSomeService, ""},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, ErrCtorNotFunc, err)
			},
		},
		{
			name:  "invalid-ctor-signature",
			ctors: []interface{}{newSomeService, func() {}},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, ErrWrongCtorSignature, err)
			},
		},
		// Validation
		{
			name:  "success-no-pointer-single-dep",
			ctors: []any{newSomeServiceNonPointerDep, newSomeDependency},
			expectedServices: map[reflect.Type]constructor{
				reflection.TypeOf[*someDependency](): {
					fn: newSomeDependency,
				},
				reflection.TypeOf[someService](): {
					fn:     newSomeServiceNonPointerDep,
					params: []reflect.Type{reflect.TypeOf(newSomeDependency())},
				},
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-missing-single-dep",
			ctors: []any{newSomeServiceWithDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "'dependency is not registered': 'di.someDependency'")
			},
		},
		{
			name:  "invalid-missing-two-deps",
			ctors: []any{newSomeServiceWithTwoDeps},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "dependency is not registered': 'di.someDependency")
				assert.ErrorWith(t, err, "dependency is not registered': 'di.someOtherDependency")
			},
		},
		{
			name:  "invalid-cyclic-dependency",
			ctors: []any{newSomeServiceWithCyclicDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "cyclic dependency detected': 'di.someService")
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
				return Register(newSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someInterface
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*someService](),
		},
		{
			name: "success-service-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(newSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[someService](),
		},
		{
			name: "success-service-pointer-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(newSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service *someService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*someService](),
		},
		{
			name: "invalid-service-pointer-to-pointer-dep",
			container: func(t *testing.T) (*container, error) {
				return Register(newSomeServiceWithDep, newSomeDependency)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, ErrNotPointer, recover())
					t.SkipNow()
				}()

				var service *someService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*someService](),
		},
		{
			name: "invalid-service-pointer-no-addresable-no-deps",
			container: func(t *testing.T) (*container, error) {
				return Register(newSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, ErrNotAddresable, recover())
					t.SkipNow()
				}()

				var service *someService
				provider(service)
				return service
			},
			expectedType: reflection.TypeOf[*someService](),
		},
		{
			name: "invalid-service-not-registered-no-deps",
			container: func(t *testing.T) (*container, error) {
				return &container{}, nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					err, _ := recover().(error)
					assert.ErrorWith(t, err, ErrNotRegistered.Error())
					t.SkipNow()
				}()

				var service someService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*someService](),
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
				return Register(newSomeService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-singleton-root-cached",
			container: func(t *testing.T) (*container, error) {
				return Register(Construct(Singleton, newSomeService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-scoped-root-no-cached",
			container: func(t *testing.T) (*container, error) {
				return Register(Construct(Scoped, newSomeService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-scoped-scope-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Scoped, newSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-singleton-scope-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Singleton, newSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-transient-scope-no-cached",
			container: func(t *testing.T) (*container, error) {
				c, err := Register(Construct(Transient, newSomeService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someService
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
