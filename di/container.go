// Package di provides a dependency container where all dependencies can be registered and cached with different
// lifetime scope. It adds validation for common mistakes like missing or cyclic dependency.
// Container allows to provide service implementation without knowing how to construct it. The container is
// responsible for caching, managing dependencies and creating the objects.
package di

import (
	"errors"
	"fmt"
	"goflow/exception"
	"reflect"
)

// Container is implemented by any value that has a Validate, Provide and Register method.
// The implementation controls how constructors are registered or provided inside container and what
// are requirements must be met to consider container as valid instance.
type Container interface {
	Validate() error
	Provide(v interface{})
	Register(ctors ...any)
}

type container struct {
	services map[reflect.Type]constructor
	cache    Cache
}

var (
	ErrNotRegistered    = errors.New("dependency is not registered")
	ErrCyclicDependency = errors.New("cyclic dependency detected")
	ErrNotAddresable    = errors.New("need to pass address")
	ErrNotPointer       = errors.New("must be a pointer")
)

// Register returns a new container instance with constructor services. Construct or func constructor
// can be passed. Error will be returned if any func constructor is not valid or Validate on
// container will return error.
func Register(ctors ...any) (c *container, err error) {
	defer exception.HandlePanicError(func(rerr error) {
		c = nil
		err = rerr
	})

	services := make(map[reflect.Type]constructor, len(ctors))

	for _, ctor := range ctors {
		construct, ok := ctor.(constructor)
		if !ok {
			constr, ok := ctor.(*constructor)
			if !ok {
				construct = *Construct(Transient, ctor)
			} else {
				construct = *constr
			}
		}

		services[construct.typ] = construct
	}

	c = &container{
		services: services,
		cache:    NewRootCache(),
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Validate verifies if every dependency is registered to provide services without
// missing dependency issue and tests if there is no cyclic dependency.
func (c *container) Validate() error {
	errs := make([]error, 0)
	for serviceType, serviceCtor := range c.services {
		for _, dependencyType := range serviceCtor.params {
			cyclic := serviceType == dependencyType
			if !cyclic {
				var otherType reflect.Type
				if serviceType.Kind() == reflect.Pointer {
					otherType = serviceType.Elem()
				} else {
					otherType = reflect.PointerTo(serviceType)
				}

				cyclic = otherType == dependencyType
			}

			if cyclic {
				errs = append(errs, fmt.Errorf("'%w': '%v'", ErrCyclicDependency, dependencyType))
				continue
			}

			_, ok := checkRegistered(dependencyType, c.services)
			if !ok {
				errs = append(errs, fmt.Errorf("'%w': '%v'", ErrNotRegistered, dependencyType))
			}
		}
	}

	return exception.Aggregate(errs...)
}

// Scope returns a new scoped container which will cache scoped lifetime services.
func (c *container) Scope() *container {
	scoped := &container{
		services: c.services,
		cache:    NewScopeCache(c.cache),
	}

	return scoped
}

// Provide creates service from found constructor and sets the result to v. It will panic
// if constructor cannot be found or v is not a pointer. Singleton services are cahced in
// root cache. If container is scoped it will cache also Scoped services. Transient services
// are always recreated from constructor.
func (c *container) Provide(v interface{}) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Pointer {
		panic(ErrNotPointer)
	}
	service := c.get(typ)

	c.setValue(v, service)
}

// setValue sets service value to v pointer.
func (c *container) setValue(v interface{}, service interface{}) {
	vval := reflect.ValueOf(v)
	velem := vval.Elem()
	serviceValue := reflect.ValueOf(service)
	if !velem.IsValid() {
		panic(ErrNotAddresable)
	}

	if velem.Kind() == reflect.Interface {
		velem.Set(serviceValue)
		return
	}

	if serviceValue.Kind() == reflect.Pointer {
		if velem.Kind() == reflect.Pointer {
			velem.Set(serviceValue)
			return
		}

		velem.Set(serviceValue.Elem())
		return
	}

	panic(fmt.Sprintf("cannot set value for '%v'", service))
}

// get returns service value for typ. Can retrieve it from cache if applicable.
func (c *container) get(typ reflect.Type) interface{} {
	ctor, ok := checkRegistered(typ, c.services)
	if !ok {
		panic(fmt.Errorf("'%w': '%v'", ErrNotRegistered, typ))
	}

	service, ok := c.cache.Get(ctor.life, ctor.typ)
	if ok {
		return service
	}

	service = ctor.Create(c.get)

	c.cache.Put(ctor.life, ctor.typ, service)

	return service
}
