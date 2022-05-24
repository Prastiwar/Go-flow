package di

import "reflect"

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

type Constructor[T any] struct {
	ctor interface{}
}

func NewConstructor[T any, C any](ctor C) *Constructor[T] {
	return &Constructor[T]{ctor: ctor}
}

func RegisterSingleton[I any, T any, C *Constructor[T]](s ServiceCollection, ctor C) {
	fac := func(provider ServiceProvider) (interface{}, error) {
		method := reflect.ValueOf(ctor).MethodByName("ctor")
		paramLen := method.Type().NumIn()
		paramValues := make([]reflect.Value, paramLen)
		for i := 0; i < paramLen; i++ {
			t := method.Type().In(i)
			object, err := provider.GetService(t)
			if err != nil {
				return nil, err
			}
			paramValues[i] = reflect.ValueOf(object)
		}
		service := method.Call(paramValues)[0] // TODO: cast to T
		return service, nil
	}
	d := NewServiceDescriptor[I, T](Singleton, fac)
	s.Register(d)
}

func RegisterSingletonWithFactory[I any, T any](s ServiceCollection, fac func(ServiceProvider) (*T, error)) {
	ctor := func(provider ServiceProvider) (interface{}, error) { return fac(provider) }
	d := NewServiceDescriptor[I, T](Singleton, ctor)
	s.Register(d)
}

func RegisterSingletonWithInstance[I any, T any](s ServiceCollection, service *T) {
	ctor := func(provider ServiceProvider) (interface{}, error) { return service, nil }
	d := NewServiceDescriptor[I, T](Singleton, ctor)
	s.Register(d)
}
