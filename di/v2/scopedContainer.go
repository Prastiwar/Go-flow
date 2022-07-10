package v2

import (
	"fmt"
	"reflect"
)

type scopedContainer struct {
	parent container
	cache  map[reflect.Type]interface{}
}

func (c *scopedContainer) Provide(v interface{}) {
	typ := reflect.TypeOf(v)
	service := c.get(typ)

	c.parent.setValue(v, service)
}

func (c *scopedContainer) get(typ reflect.Type) interface{} {
	if typ.Kind() != reflect.Pointer {
		panic("must be pointer")
	}

	ctor, ok := c.parent.checkRegistered(typ)
	if !ok {
		panic(fmt.Errorf("'%w': '%v'", NotRegisteredError, typ))
	}

	if ctor.life == Singleton {
		service, ok := c.parent.cache[ctor.typ]
		if ok {
			return service
		}
	} else if ctor.life == Scoped {
		service, ok := c.cache[ctor.typ]
		if ok {
			return service
		}
	}

	service := ctor.Create(c.get)

	if ctor.life == Singleton {
		c.parent.cache[ctor.typ] = service
	} else if ctor.life == Scoped {
		c.cache[ctor.typ] = service
	}

	return service
}
