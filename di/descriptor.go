package di

import (
	"errors"
	"fmt"
	"goflow/reflection"
	"reflect"
)

type LifeTime int

const (
	Transient LifeTime = iota
	Singleton
)

type ServiceFactory func(provider ServiceProvider) (interface{}, error)

type ServiceDescriptor interface {
	Interface() reflect.Type
	Service() reflect.Type
	LifeTime() LifeTime
	Factory() ServiceFactory
}

type descriptor struct {
	interfaceType reflect.Type
	serviceType   reflect.Type
	factory       ServiceFactory
	life          LifeTime
}

func (d *descriptor) Interface() reflect.Type {
	return d.interfaceType
}

func (d *descriptor) Service() reflect.Type {
	return d.serviceType
}

func (d *descriptor) LifeTime() LifeTime {
	return d.life
}

func (d *descriptor) Factory() ServiceFactory {
	return d.factory
}

func NewServiceDescriptor[I any, S any](life LifeTime, fac ServiceFactory) (ServiceDescriptor, error) {
	interfaceType := reflection.TypeOf[I]()
	serviceType := reflection.TypeOf[S]()

	d := &descriptor{
		interfaceType: interfaceType,
		serviceType:   serviceType,
		factory:       fac,
		life:          life,
	}

	err := validateDescriptor(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func validateDescriptor(descriptor ServiceDescriptor) error {
	life := descriptor.LifeTime()
	if life > Singleton {
		msg := fmt.Sprintf("'%v' is not valid LifeTime", life)
		return errors.New(msg)
	}

	interfaceType := descriptor.Interface()
	if interfaceType.Kind() != reflect.Interface {
		msg := fmt.Sprintf("'%v' is not interface", interfaceType.Name())
		return errors.New(msg)
	}

	serviceType := descriptor.Service()
	if serviceType.Kind() != reflect.Struct {
		msg := fmt.Sprintf("'%v' is not struct", interfaceType.Name())
		return errors.New(msg)
	}

	ok := serviceType.AssignableTo(interfaceType) && serviceType.Implements(interfaceType)
	if !ok {
		msg := fmt.Sprintf("'%v' is not assignable from '%v'", interfaceType.Name(), serviceType.Name())
		return errors.New(msg)
	}

	return nil
}
