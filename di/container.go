package di

import (
	"errors"
	"fmt"
	"goflow/exception"
	"reflect"
)

type Container interface {
	Validate() error
	Provide(v interface{})
	Register(ctors ...any)
}

type container struct {
	services map[reflect.Type]constructor
	cache    DiCache
}

var (
	ErrNotRegistered    = errors.New("dependency is not registered")
	ErrCyclicDependency = errors.New("cyclic dependency detected")
	ErrNotAddresable    = errors.New("need to pass address")
	ErrNotPointer       = errors.New("must be a pointer")
)

func Register(ctors ...any) (*container, error) {
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

		err := construct.Validate()
		if err != nil {
			return nil, err
		}

		services[construct.typ] = construct
	}

	c := &container{
		services: services,
		cache:    NewRootCache(),
	}

	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

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

func (c *container) Scope() *container {
	scoped := &container{
		services: c.services,
		cache:    NewScopeCache(c.cache),
	}

	return scoped
}

func (c *container) Provide(v interface{}) {
	typ := reflect.TypeOf(v)
	service := c.get(typ)

	c.setValue(v, service)
}

func (c *container) setValue(v interface{}, service interface{}) {
	vval := reflect.ValueOf(v)
	velem := vval.Elem()
	serviceValue := reflect.ValueOf(service)
	if !velem.IsValid() {
		panic(ErrNotAddresable)
	}

	if velem.Kind() == reflect.Interface {
		velem.Set(serviceValue)
	} else if serviceValue.Kind() == reflect.Pointer {
		if velem.Kind() == reflect.Pointer {
			velem.Set(serviceValue)
		} else {
			velem.Set(serviceValue.Elem())
		}
	} else {
		panic(fmt.Sprintf("cannot set value for '%v'", service))
	}
}

func (c *container) get(typ reflect.Type) interface{} {
	if typ.Kind() != reflect.Pointer {
		panic(ErrNotPointer)
	}

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
