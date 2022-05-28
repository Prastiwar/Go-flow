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

func RegisterWithConstructor[I any, T any, C any](s ServiceCollection, life LifeTime, ctor C) error {
	fac, err := createServiceFactory[I](ctor)
	if err != nil {
		return err
	}

	d, err := NewServiceDescriptor[I, T](life, fac)
	if err != nil {
		return err
	}

	s.Register(d)
	return nil
}

func RegisterWithFactory[I any, T any](s ServiceCollection, life LifeTime, fac func(ServiceProvider) (*T, error)) error {
	ctor := func(provider ServiceProvider) (interface{}, error) { return fac(provider) }
	d, err := NewServiceDescriptor[I, T](life, ctor)
	if err != nil {
		return err
	}

	s.Register(d)
	return nil
}

func RegisterTransient[I any, T any, C any](s ServiceCollection, ctor C) error {
	return RegisterWithConstructor[I, T](s, Transient, ctor)
}

func RegisterTransientWithFactory[I any, T any](s ServiceCollection, fac func(ServiceProvider) (*T, error)) error {
	return RegisterWithFactory[I](s, Transient, fac)
}

func RegisterSingleton[I any, T any, C any](s ServiceCollection, ctor C) error {
	return RegisterWithConstructor[I, T](s, Singleton, ctor)
}

func RegisterSingletonWithFactory[I any, T any](s ServiceCollection, fac func(ServiceProvider) (*T, error)) error {
	return RegisterWithFactory[I](s, Singleton, fac)
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

func createServiceFactory[I any, T any](ctor T) (ServiceFactory, error) {
	method := reflect.ValueOf(ctor)
	if method.Kind() != reflect.Func {
		return nil, errors.New("constructor must be a function")
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

	return fac, nil
}
