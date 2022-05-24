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
	descriptors := []ServiceDescriptor{
		NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			return expectedService, nil
		}),
	}

	provider := newProvider(descriptors)

	assert.Equal(t, len(provider.descriptors), 1)
}

func TestProvide(t *testing.T) {
	expectedService := mocks.NewStringerMock("test")
	descriptors := []ServiceDescriptor{
		NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			return expectedService, nil
		}),
	}
	provider := newProvider(descriptors)

	var stringer fmt.Stringer
	err := Provide(provider, &stringer)

	assert.NilError(t, err)
}

func TestGetServiceSuccess(t *testing.T) {
	descriptors := []ServiceDescriptor{
		NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			return nil, mocks.NewErrorMock("not registered")
		}),
		NewServiceDescriptor[error, mocks.ErrorMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			panic("panic")
		}),
		NewServiceDescriptor[fmt.GoStringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			panic(mocks.NewErrorMock("something went wrong"))
		}),
		NewServiceDescriptor[ServiceProvider, mocks.ServiceProviderMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
			panic(time.Now())
		}),
	}
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
