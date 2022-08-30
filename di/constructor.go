package di

import (
	"errors"
	"reflect"

	"github.com/Prastiwar/Go-flow/reflection"
)

type LifeTime int

const (
	Transient LifeTime = iota
	Singleton
	Scoped
)

var (
	ErrCtorNotFunc        = errors.New("ctor is not func")
	ErrWrongCtorSignature = errors.New("ctor must return only service value")
)

type constructor struct {
	typ    reflect.Type
	fn     interface{}
	params []reflect.Type
	life   LifeTime
}

// Construct returns a new constructor instance for specified function.
// It panics if ctor is not a function or it does not return any value.
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

	c := &constructor{
		typ:    typ,
		fn:     ctor,
		life:   life,
		params: inParamTypes,
	}

	err := c.validate()
	if err != nil {
		panic(err)
	}

	return c
}

// validate verifies if ctor function is func and returns single value.
func (c *constructor) validate() error {
	ctorValue := reflect.ValueOf(c.fn)
	if ctorValue.Kind() != reflect.Func {
		return ErrCtorNotFunc
	}

	paramTypes := reflection.OutParamTypes(ctorValue.Type())
	if len(paramTypes) != 1 {
		return ErrWrongCtorSignature
	}

	return nil
}

// Create returns created instance from called ctor with parameters retrieved with provider.
func (c *constructor) Create(provider func(reflect.Type) interface{}) interface{} {
	paramValues := make([]reflect.Value, len(c.params))
	for i := 0; i < len(c.params); i++ {
		t := c.params[i]
		object := provider(t)

		v, err := reflection.GetFieldValueFor(t, object)
		if err != nil {
			panic(err)
		}

		paramValues[i] = v
	}

	method := reflect.ValueOf(c.fn)
	service := method.Call(paramValues)[0].Interface()

	return service
}
