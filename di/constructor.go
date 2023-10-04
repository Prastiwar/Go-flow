package di

import (
	"errors"
	"reflect"

	"github.com/Prastiwar/Go-flow/reflection"
)

// LifeTime is value defining scope of life for object.
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

// Constructor is interface for delegate the Create function which has passed the provider for dependency resolving while
// instantiating a new instace of concrete type.
type Constructor interface {
	// Create returns created instance from called ctor with parameters retrieved with provider.
	Create(provider func(reflect.Type) interface{}) interface{}

	// Type returns a reflect.Type for object that can be created by this constructor.
	Type() reflect.Type

	// Dependencies returns an array of reflect.Type defining type of dependencies for object to construct.
	Dependencies() []reflect.Type

	// Life returns LifeTime which defines scope of existence for constructed object.
	Life() LifeTime
}

// ConstructorFunc is simple func type that implements Constructor.
type ConstructorFunc func(provider func(reflect.Type) interface{}) interface{}

func (f ConstructorFunc) Create(provider func(reflect.Type) interface{}) interface{} {
	return f(provider)
}

var _ Constructor = &constructor{}

type constructor struct {
	typ    reflect.Type
	fn     interface{}
	params []reflect.Type
	life   LifeTime
}

// Construct returns a new Constructor instance for specified function.
// It panics if ctor is not a function or it does not return any value.
func Construct(life LifeTime, ctor any) Constructor {
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

func (c *constructor) Type() reflect.Type {
	return c.typ
}

func (c *constructor) Dependencies() []reflect.Type {
	return c.params
}
func (c *constructor) Life() LifeTime {
	return c.life
}
