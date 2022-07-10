package v2

import (
	"errors"
	"fmt"
	"goflow/exception"
	"goflow/reflection"
	"reflect"
)

type Container interface {
	Validate() error
}

type container struct {
	services map[reflect.Type]constructor
	cache    map[reflect.Type]interface{}
}

var (
	NotRegisteredError    = errors.New("dependency is not registered")
	CyclicDependencyError = errors.New("cyclic dependency detected")
	NotAddresableError    = errors.New("need to pass address to v")
)

func Register(ctors ...any) (*container, error) {
	services := make(map[reflect.Type]constructor, len(ctors))

	for _, ctor := range ctors {
		construct, ok := ctor.(constructor)
		if !ok {
			construct = *Construct(Transient, ctor)
		}

		err := construct.Validate()
		if err != nil {
			return nil, err
		}

		services[construct.typ] = construct
	}

	return &container{
		services: services,
	}, nil
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
				errs = append(errs, fmt.Errorf("'%w': '%v'", CyclicDependencyError, dependencyType))
				continue
			}

			_, ok := c.checkRegistered(dependencyType)
			if !ok {
				errs = append(errs, fmt.Errorf("'%w': '%v'", NotRegisteredError, dependencyType))
			}
		}
	}

	return exception.Aggregate(errs...)
}

func (c *container) Scope() *container {
	scoped := &container{
		services: c.services,
		cache:    c.cache,
		// TODO: cache handling
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
		panic(NotAddresableError)
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
		panic("unhandled case")
	}
}

func (c *container) get(typ reflect.Type) interface{} {
	if typ.Kind() != reflect.Pointer {
		panic("must be pointer")
	}

	ctor, ok := c.checkRegistered(typ)
	if !ok {
		panic(fmt.Errorf("'%w': '%v'", NotRegisteredError, typ))
	}

	if ctor.life == Singleton {
		service, ok := c.cache[ctor.typ]
		if ok {
			return service
		}
	}

	service := ctor.Create(c.get)

	if ctor.life == Singleton {
		c.cache[ctor.typ] = service
	}

	return service
}

func (c *container) checkInterface(u reflect.Type) (constructor, bool) {
	for serviceType, ctor := range c.services {
		ok := serviceType.Implements(u)
		if ok {
			return ctor, true
		}
	}

	return constructor{}, false
}

func (c *container) checkRegistered(u reflect.Type) (constructor, bool) {
	if u.Kind() == reflect.Interface {
		return c.checkInterface(u)
	}

	ctor, ok := c.services[u]
	if !ok {
		otherType := reflection.TogglePointer(u)
		if otherType.Kind() == reflect.Interface {
			return c.checkInterface(otherType)
		}

		ctor, ok := c.services[otherType]
		return ctor, ok
	}

	return ctor, ok
}
