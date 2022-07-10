package v2

import (
	"goflow/reflection"
	"reflect"
)

func checkInterface(u reflect.Type, services map[reflect.Type]constructor) (constructor, bool) {
	for serviceType, ctor := range services {
		ok := serviceType.Implements(u)
		if ok {
			return ctor, true
		}
	}

	return constructor{}, false
}

func checkRegistered(u reflect.Type, services map[reflect.Type]constructor) (constructor, bool) {
	if u.Kind() == reflect.Interface {
		return checkInterface(u, services)
	}

	ctor, ok := services[u]
	if !ok {
		otherType := reflection.TogglePointer(u)
		if otherType.Kind() == reflect.Interface {
			return checkInterface(otherType, services)
		}

		ctor, ok := services[otherType]
		return ctor, ok
	}

	return ctor, ok
}
