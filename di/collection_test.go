package di

import (
	"fmt"
	"goflow/di/mocks"
	"goflow/tests/assert"
	"reflect"
	"testing"
)

type Service interface{}
type ServiceImplementation struct{}

func newStringerImplementation(s Service) *mocks.StringerMock {
	return &mocks.StringerMock{}
}

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

func TestBuildProvider(t *testing.T) {
	services := NewServiceCollection()

	provider := BuildProvider(services)

	assert.NotNil(t, provider)
}
