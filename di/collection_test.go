package di

import (
	"errors"
	"fmt"
	"goflow/di/mocks"
	"goflow/tests/assert"
	"reflect"
	"testing"
)

type Service interface{}
type ServiceImplementation struct{}

// newStringerImplementation is constructor with dependencies
func newStringerImplementation(s Service) *mocks.StringerMock {
	return &mocks.StringerMock{}
}

// provideStringerImplementation is factory method leading to constructor
func provideStringerImplementation(provider ServiceProvider) (*mocks.StringerMock, error) {
	s, err := GetService[Service](provider)
	if err != nil {
		return nil, err
	}

	return newStringerImplementation(s), nil
}

func TestRegisterSingleton(t *testing.T) {
	services := NewServiceCollection()
	provider := mocks.NewServiceProviderMock(func(serviceType reflect.Type) (interface{}, error) { return ServiceImplementation{}, nil })

	err := RegisterSingleton[fmt.Stringer, mocks.StringerMock](services, newStringerImplementation)
	assert.NilError(t, err)

	service, err := services.Descriptors()[0].Factory()(provider)
	assert.NilError(t, err)
	assert.NotNil(t, service)
}

func TestRegisterSingletonInvalidCtor(t *testing.T) {
	services := NewServiceCollection()

	err := RegisterSingleton[fmt.Stringer, mocks.StringerMock](services, "test")

	assert.ErrorWith(t, err, "constructor")
	assert.Equal(t, 0, len(services.Descriptors()))
}

func TestRegisterSingletonInvalidImplementation(t *testing.T) {
	services := NewServiceCollection()

	err := RegisterSingleton[fmt.Stringer, mocks.ErrorMock](services, newStringerImplementation)

	assert.ErrorWith(t, err, "not assignable")
	assert.Equal(t, 0, len(services.Descriptors()))
}

func TestRegisterSingletonDependencyNotRegistered(t *testing.T) {
	services := NewServiceCollection()
	provider := mocks.NewServiceProviderMock(func(serviceType reflect.Type) (interface{}, error) { return nil, errors.New("not registered") })

	err := RegisterSingleton[fmt.Stringer, mocks.StringerMock](services, newStringerImplementation)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(services.Descriptors()))

	_, err = services.Descriptors()[0].Factory()(provider)
	assert.ErrorWith(t, err, "not registered")
}

func TestRegisterSingletonWithFactory(t *testing.T) {
	services := NewServiceCollection()
	provider := mocks.NewServiceProviderMock(func(serviceType reflect.Type) (interface{}, error) { return ServiceImplementation{}, nil })

	err := RegisterSingletonWithFactory[fmt.Stringer](services, provideStringerImplementation)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(services.Descriptors()))
	service, err := services.Descriptors()[0].Factory()(provider)
	assert.NilError(t, err)
	assert.NotNil(t, service)
}

func TestRegisterSingletonWithFactoryInvalid(t *testing.T) {
	services := NewServiceCollection()

	err := RegisterSingletonWithFactory[fmt.Stringer](services, func(provider ServiceProvider) (*string, error) {
		str := "test"
		return &str, nil
	})

	assert.Error(t, err)
	assert.Equal(t, 0, len(services.Descriptors()))
}

func TestRegisterSingletonWithInstance(t *testing.T) {
	services := NewServiceCollection()

	instance := mocks.NewStringerMock("test")
	err := RegisterSingletonWithInstance[fmt.Stringer](services, instance)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(services.Descriptors()))

	actualInstance, err := services.Descriptors()[0].Factory()(mocks.ServiceProviderMock{})
	assert.NilError(t, err)
	assert.Equal(t, instance, actualInstance)
}

func TestRegisterSingletonWithInstanceInvalid(t *testing.T) {
	services := NewServiceCollection()

	instance := "test"
	err := RegisterSingletonWithInstance[fmt.Stringer](services, &instance)

	assert.Error(t, err)
	assert.Equal(t, 0, len(services.Descriptors()))
}

func TestBuildProvider(t *testing.T) {
	services := NewServiceCollection()

	provider := BuildProvider(services)

	assert.NotNil(t, provider)
}
