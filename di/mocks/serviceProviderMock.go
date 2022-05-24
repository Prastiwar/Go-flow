package mocks

import "reflect"

type ServiceProviderMock struct {
	getService func(serviceType reflect.Type) (interface{}, error)
}

func (m ServiceProviderMock) GetService(serviceType reflect.Type) (interface{}, error) {
	return m.getService(serviceType)
}

func NewServiceProviderMock(getService func(serviceType reflect.Type) (interface{}, error)) *ServiceProviderMock {
	return &ServiceProviderMock{
		getService: getService,
	}
}
