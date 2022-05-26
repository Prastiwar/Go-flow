package di

import (
	"errors"
	"fmt"
	"goflow/reflection"
	"reflect"
	"sync"
)

var ErrNotRegistered = errors.New("service was not registered")

type ServiceProvider interface {
	GetService(serviceType reflect.Type) (interface{}, error)
}

type provider struct {
	sync.Mutex
	descriptors map[reflect.Type]ServiceDescriptor
	singletons  map[reflect.Type]interface{}
}

func newProvider(descriptors []ServiceDescriptor) *provider {
	m := make(map[reflect.Type]ServiceDescriptor, len(descriptors))
	for _, d := range descriptors {
		t := d.Interface()
		m[t] = d
	}
	return &provider{descriptors: m, singletons: make(map[reflect.Type]interface{})}
}

func (p *provider) GetService(serviceType reflect.Type) (instance interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			str, ok := r.(string)
			if !ok {
				str = "error occured while instantiating service"
			}
			e, ok := r.(error)
			if !ok {
				err = errors.New(str)
			} else {
				err = e
			}
		}
	}()

	service, ok := p.singletons[serviceType]
	if ok {
		return service, nil
	}

	descr, ok := p.descriptors[serviceType]
	if !ok {
		return nil, ErrNotRegistered
	}

	service, err = descr.Factory()(p)
	if err != nil {
		return nil, err
	}

	if descr.LifeTime() == Singleton {
		p.Lock()
		defer p.Unlock()
		p.singletons[serviceType] = service
	}

	return service, nil
}

func GetService[T interface{}](provider ServiceProvider) (*T, error) {
	s, err := provider.GetService(reflection.TypeOf[T]())
	if err != nil {
		return nil, err
	}

	service, ok := s.(T)
	if !ok {
		msg := fmt.Sprintf("Unable to cast object of type '%v' to type '%v'", reflect.TypeOf(s), reflect.TypeOf(&service).Elem())
		return nil, errors.New(msg)
	}

	return &service, nil
}

func Provide[T any](provider *provider, service *T) error {
	s, err := GetService[T](provider)
	if err != nil {
		return err
	}

	*service = *s
	return nil
}
