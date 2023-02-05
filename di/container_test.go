package di_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/di"
	"github.com/Prastiwar/Go-flow/reflection"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type someInterface interface{}

type fooDependency struct {
	id float64
}

type someOtherDependency struct{}

type fooService struct {
	id float64
}

func newfooDependency() *fooDependency {
	return &fooDependency{
		id: rand.Float64(),
	}
}

func newFooService() *fooService {
	return &fooService{
		id: rand.Float64(),
	}
}

func newFooServiceNonPointerDep(*fooDependency) fooService {
	return fooService{
		id: rand.Float64(),
	}
}

func newFooServiceWithDep(dep fooDependency) *fooService {
	return newFooService()
}

func newFooServiceWithTwoDeps(dep fooDependency, otherDep someOtherDependency) *fooService {
	return newFooService()
}

func newFooServiceWithCyclicDep(dep fooService) *fooService {
	return newFooService()
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name             string
		ctors            []interface{}
		expectedServices map[reflect.Type]di.Constructor
		assertErr        assert.ErrorFunc
	}{
		{
			name:  "success-some-service",
			ctors: []interface{}{newFooService},
			expectedServices: map[reflect.Type]di.Constructor{
				reflection.TypeOf[*fooService](): di.Construct(di.Transient, newFooService),
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-pointer",
			ctors: []interface{}{di.Construct(di.Scoped, newFooService)},
			expectedServices: map[reflect.Type]di.Constructor{
				reflection.TypeOf[*fooService](): di.Construct(di.Transient, newFooService),
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "success-ctor-non-pointer",
			ctors: []interface{}{di.Construct(di.Scoped, newFooService)},
			expectedServices: map[reflect.Type]di.Constructor{
				reflection.TypeOf[*fooService](): di.Construct(di.Transient, newFooService),
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-ctor-func",
			ctors: []interface{}{newFooService, ""},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, di.ErrCtorNotFunc, err)
			},
		},
		{
			name:  "invalid-ctor-signature",
			ctors: []interface{}{newFooService, func() {}},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, di.ErrWrongCtorSignature, err)
			},
		},
		// Validation
		{
			name:  "success-no-pointer-single-dep",
			ctors: []any{newFooServiceNonPointerDep, newfooDependency},
			expectedServices: map[reflect.Type]di.Constructor{
				reflection.TypeOf[*fooDependency](): di.Construct(di.Transient, newfooDependency),
				reflection.TypeOf[fooService]():     di.Construct(di.Transient, newFooServiceNonPointerDep),
			},
			assertErr: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:  "invalid-missing-single-dep",
			ctors: []any{newFooServiceWithDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "'dependency is not registered': 'di_test.fooDependency'")
			},
		},
		{
			name:  "invalid-missing-two-deps",
			ctors: []any{newFooServiceWithTwoDeps},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "dependency is not registered': 'di_test.fooDependency")
				assert.ErrorWith(t, err, "dependency is not registered': 'di_test.someOtherDependency")
			},
		},
		{
			name:  "invalid-cyclic-dependency",
			ctors: []any{newFooServiceWithCyclicDep},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorWith(t, err, "cyclic dependency detected': 'di_test.fooService")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := di.Register(tt.ctors...)
			services := make(map[reflect.Type]di.Constructor, 4)
			if container != nil {
				containerServices := container.Services()
				for _, s := range containerServices {
					services[s.Type()] = s.Constructor()
				}
			}

			tt.assertErr(t, err)
			if t.Failed() {
				t.FailNow()
			}

			assert.Equal(t, len(tt.expectedServices), len(services))
			for typ, ctor := range tt.expectedServices {
				actualCtor, exists := services[typ]
				assert.Equal(t, true, exists, "exists assertion failed")
				if exists {
					// assert.Equal(t, reflect.TypeOf(ctor.fn), reflect.TypeOf(actualCtor.fn), "type equality failed")
					assert.ElementsMatch(t, ctor.Dependencies(), actualCtor.Dependencies())
					assert.Equal(t, len(ctor.Dependencies()), len(actualCtor.Dependencies()), "ctor params equality failed")
				}
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
		container    func(t *testing.T) (di.Container, error)
		provideFn    func(t *testing.T, provider func(any)) any
		expectedType reflect.Type
	}{
		{
			name: "success-interface-no-deps",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service someInterface
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*fooService](),
		},
		{
			name: "success-service-no-deps",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[fooService](),
		},
		{
			name: "success-service-pointer-no-deps",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service *fooService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*fooService](),
		},
		{
			name: "success-service-dep-pointer-to-non-pointer",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooServiceWithDep, newfooDependency)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, nil, recover())
				}()

				var service fooService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[fooService](),
		},
		{
			name: "success-service-pointer-to-pointer-dep",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooServiceWithDep, newfooDependency)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, nil, recover())
				}()

				var service *fooService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*fooService](),
		},
		{
			name: "invalid-service-pointer-no-addressable-no-deps",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					assert.Equal(t, di.ErrNotAddresable, recover())
					t.SkipNow()
				}()

				var service *fooService
				provider(service)
				return service
			},
			expectedType: reflection.TypeOf[*fooService](),
		},
		{
			name: "invalid-service-not-registered-no-deps",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register()
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					err, _ := recover().(error)
					assert.ErrorWith(t, err, di.ErrNotRegistered.Error())
					t.SkipNow()
				}()

				var service fooService
				provider(&service)
				return service
			},
			expectedType: reflection.TypeOf[*fooService](),
		},
		{
			name: "invalid-service-not-pointer",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register()
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				defer func() {
					err, _ := recover().(error)
					assert.ErrorWith(t, err, di.ErrNotPointer.Error())
					t.SkipNow()
				}()

				var service fooService
				provider(service)
				return service
			},
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
		container func(t *testing.T) (di.Container, error)
		provideFn func(t *testing.T, provider func(any)) any
		cached    bool
	}{
		{
			name: "success-transient-root-no-cached",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(newFooService)
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-singleton-root-cached",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(di.Construct(di.Singleton, newFooService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-di.Scoped-root-no-cached",
			container: func(t *testing.T) (di.Container, error) {
				return di.Register(di.Construct(di.Scoped, newFooService))
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			cached: false,
		},
		{
			name: "success-di.Scoped-scope-cached",
			container: func(t *testing.T) (di.Container, error) {
				c, err := di.Register(di.Construct(di.Scoped, newFooService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-singleton-scope-cached",
			container: func(t *testing.T) (di.Container, error) {
				c, err := di.Register(di.Construct(di.Singleton, newFooService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
				provider(&service)
				return service
			},
			cached: true,
		},
		{
			name: "success-transient-scope-no-cached",
			container: func(t *testing.T) (di.Container, error) {
				c, err := di.Register(di.Construct(di.Transient, newFooService))
				if err != nil {
					return nil, err
				}
				return c.Scope(), nil
			},
			provideFn: func(t *testing.T, provider func(any)) any {
				var service fooService
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
