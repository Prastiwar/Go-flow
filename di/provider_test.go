package di

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"goflow/di/mocks"
	"goflow/reflection"
	"goflow/tests/assert"
)

func TestNewProvider(t *testing.T) {
	expectedService := mocks.NewStringerMock("test")
	d, _ := NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return expectedService, nil
	})
	descriptors := []ServiceDescriptor{d}

	provider := newProvider(descriptors)

	assert.Equal(t, len(provider.descriptors), 1)
}

func TestProvide(t *testing.T) {
	expectedService := mocks.NewStringerMock("test")
	d, _ := NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return expectedService, nil
	})
	descriptors := []ServiceDescriptor{d}
	provider := newProvider(descriptors)

	var stringer fmt.Stringer
	err := Provide(provider, &stringer)

	assert.NilError(t, err)
}

func TestProvideError(t *testing.T) {
	provider := newProvider([]ServiceDescriptor{})

	var stringer fmt.Stringer
	err := Provide(provider, &stringer)

	assert.Error(t, err)
}

func TestGetServiceSuccess(t *testing.T) {
	expectedService := mocks.NewStringerMock("test")
	d, _ := NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return expectedService, nil
	})
	descriptors := []ServiceDescriptor{d}
	provider := newProvider(descriptors)

	s, err := provider.GetService(reflection.TypeOf[fmt.Stringer]())

	assert.NilError(t, err)
	assert.Equal(t, expectedService, s)
	assert.Equal(t, 1, len(provider.singletons))

	cachedService, err := provider.GetService(reflection.TypeOf[fmt.Stringer]())
	assert.NilError(t, err)
	assert.Equal(t, s, cachedService)
}

func TestGetServiceErrors(t *testing.T) {
	stringer, _ := NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return nil, mocks.NewErrorMock("not registered")
	})
	errorer, _ := NewServiceDescriptor[error, mocks.ErrorMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		panic("panic")
	})
	goStringer, _ := NewServiceDescriptor[fmt.GoStringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		panic(mocks.NewErrorMock("something went wrong"))
	})
	providerMock, _ := NewServiceDescriptor[ServiceProvider, mocks.ServiceProviderMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		panic(time.Now())
	})
	descriptors := []ServiceDescriptor{stringer, errorer, goStringer, providerMock}
	provider := newProvider(descriptors)

	cases := []struct {
		name        string
		serviceType reflect.Type
		errContent  string
	}{
		{
			name:        "factory error",
			serviceType: reflection.TypeOf[fmt.Stringer](),
			errContent:  "not registered",
		},
		{
			name:        "panic string",
			serviceType: reflection.TypeOf[error](),
			errContent:  "panic",
		},
		{
			name:        "panic error",
			serviceType: reflection.TypeOf[fmt.GoStringer](),
			errContent:  "something went wrong",
		},
		{
			name:        "panic custom",
			serviceType: reflection.TypeOf[ServiceProvider](),
			errContent:  "error occured",
		},
		{
			name:        "not registered",
			serviceType: reflection.TypeOf[fmt.Formatter](),
			errContent:  ErrNotRegistered.Error(),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := provider.GetService(test.serviceType)
			assert.ErrorWith(t, err, test.errContent)
			assert.Equal(t, 0, len(provider.singletons))
		})
	}
}

func TestGenericGetServiceSuccess(t *testing.T) {
	expectedString := "expectation"
	stringer := mocks.NewStringerMock(expectedString)
	serviceProvider := mocks.NewServiceProviderMock(
		func(serviceType reflect.Type) (interface{}, error) {
			return stringer, nil
		},
	)

	service, _ := GetService[fmt.Stringer](serviceProvider)

	str := (*service).String()
	assert.Equal(t, expectedString, str)
}

func TestGenericGetServiceUnableToCast(t *testing.T) {
	serviceProvider := mocks.NewServiceProviderMock(
		func(serviceType reflect.Type) (interface{}, error) {
			return "", nil
		},
	)

	_, err := GetService[fmt.Stringer](serviceProvider)

	assert.ErrorWith(t, err, "Unable to cast")
}
