package di

import (
	"errors"
	"goflow/reflection"
	"reflect"
)

type LifeTime int

const (
	Transient LifeTime = iota
	Singleton
	Scoped
)

var (
	NotFuncError       = errors.New("ctor is not func")
	WrongCtorSignature = errors.New("ctor must return only service value")
)

type constructor struct {
	typ    reflect.Type
	fn     interface{}
	params []reflect.Type
	life   LifeTime
}

func Construct(life LifeTime, ctor any) *constructor {
	var typ reflect.Type

	ctorValue := reflect.ValueOf(ctor)
	var inParamTypes []reflect.Type
	if ctorValue.Kind() == reflect.Func {
		paramTypes := reflection.OutParamTypes(ctorValue.Type())
		if len(paramTypes) > 0 {
			typ = paramTypes[0]
		}
		inParamTypes = reflection.InParamTypes(ctorValue.Type())
	}

	return &constructor{
		typ:    typ,
		fn:     ctor,
		life:   life,
		params: inParamTypes,
	}
}

func (c *constructor) Validate() error {
	ctorValue := reflect.ValueOf(c.fn)
	if ctorValue.Kind() != reflect.Func {
		return NotFuncError
	}

	paramTypes := reflection.OutParamTypes(ctorValue.Type())
	if len(paramTypes) != 1 {
		return WrongCtorSignature
	}

	return nil
}

func (c *constructor) Create(provider func(reflect.Type) interface{}) interface{} {
	paramValues := make([]reflect.Value, len(c.params))
	for i := 0; i < len(c.params); i++ {
		t := c.params[i]
		object := provider(t)
		paramValues[i] = reflect.ValueOf(object)
	}

	method := reflect.ValueOf(c.fn)
	service := method.Call(paramValues)[0].Interface()

	return service
}
