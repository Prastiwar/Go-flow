package di

import (
	"goflow/di/mocks"
	"goflow/tests/assert"
	"testing"
)

type Service interface{}
type ServiceImplementation struct{}

var StringerImplementationCtor = NewConstructor[mocks.StringerMock](NewStringerImplementation)

func NewStringerImplementation(s Service) *mocks.StringerMock {
	return &mocks.StringerMock{}
}

func ProvideStringerImplementation(provider ServiceProvider) (*mocks.StringerMock, error) {
	s, err := GetService[Service](provider)
	if err != nil {
		return nil, err
	}

	return NewStringerImplementation(s), nil
}

func TestRegisterSingleton(t *testing.T) {
	services := NewServiceCollection()

	// RegisterSingleton[fmt.Stringer](services, StringerImplementationCtor)

	// RegisterSingleton[fmt.Stringer, StringerImplementation](services, ProvideStringerImplementation)

	assert.Equal(t, 1, len(services.Descriptors()))
}
