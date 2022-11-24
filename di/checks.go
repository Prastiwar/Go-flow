package di

import (
	"reflect"

	"github.com/Prastiwar/Go-flow/reflection"
)

// checkInterface returns constructor found from services for service which implements typ.
func checkInterface(typ reflect.Type, services map[reflect.Type]constructor) (constructor, bool) {
	if typ.Kind() != reflect.Interface {
		return constructor{}, false
	}

	for serviceType, ctor := range services {
		ok := serviceType.Implements(typ)
		if ok {
			return ctor, true
		}
	}

	return constructor{}, false
}

// checkInterface returns constructor found from services for matched typ.
func checkRegistered(typ reflect.Type, services map[reflect.Type]constructor) (constructor, bool) {
	if typ.Kind() == reflect.Interface {
		return checkInterface(typ, services)
	}

	ctor, ok := services[typ]
	if !ok {
		otherType := reflection.TogglePointer(typ)
		if otherType.Kind() == reflect.Interface {
			return checkInterface(otherType, services)
		}

		ctor, ok := services[otherType]
		return ctor, ok
	}

	return ctor, ok
}
