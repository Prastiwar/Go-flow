package di

import (
	"errors"
	"fmt"
	"goflow/reflection"
	"reflect"
)

type ServiceCollection interface {
	Register(d ServiceDescriptor)
	Descriptors() []ServiceDescriptor
}

type serviceCollection struct {
	descriptors []ServiceDescriptor
}

func (c *serviceCollection) Register(d ServiceDescriptor) {
	c.descriptors = append(c.descriptors, d)
}

func (c *serviceCollection) Descriptors() []ServiceDescriptor {
	return c.descriptors
}

func NewServiceCollection() ServiceCollection {
	return &serviceCollection{}
}

func BuildProvider(c ServiceCollection) ServiceProvider {
	return newProvider(c.Descriptors())
}

func RegisterSingleton[I any, T any, C any](s ServiceCollection, ctor C) error {
	method := reflect.ValueOf(ctor)
	if method.Kind() != reflect.Func {
		return errors.New("constructor must be a function")
	}
	paramLen := method.Type().NumIn()
	paramValues := make([]reflect.Value, paramLen)
	paramTypes := make([]reflect.Type, paramLen)
	for i := 0; i < paramLen; i++ {
		paramTypes = append(paramTypes, method.Type().In(i))
	}

	fac := func(provider ServiceProvider) (interface{}, error) {
		for i := 0; i < paramLen; i++ {
			t := paramTypes[i]
			object, err := provider.GetService(t)
			if err != nil {
				return nil, err
			}
			paramValues[i] = reflect.ValueOf(object)
		}
		service, ok := method.Call(paramValues)[0].Interface().(I)
		if !ok {
			msg := fmt.Sprintf("Unable to cast object of type '%v' to type '%v'", reflect.TypeOf(service), reflection.TypeOf[I]())
			return nil, errors.New(msg)
		}
		return service, nil
	}

	d, err := NewServiceDescriptor[I, T](Singleton, fac)
	if err != nil {
		return err
	}

	s.Register(d)
	return nil
}

func RegisterSingletonWithFactory[I any, T any](s ServiceCollection, fac func(ServiceProvider) (*T, error)) error {
	ctor := func(provider ServiceProvider) (interface{}, error) { return fac(provider) }
	d, err := NewServiceDescriptor[I, T](Singleton, ctor)
	if err != nil {
		return err
	}

	s.Register(d)
	return nil
}

func RegisterSingletonWithInstance[I any, T any](s ServiceCollection, service *T) error {
	ctor := func(provider ServiceProvider) (interface{}, error) { return service, nil }
	d, err := NewServiceDescriptor[I, T](Singleton, ctor)
	if err != nil {
		return err
	}

	s.Register(d)
	return nil
}
