package di

import (
	"fmt"
	"goflow/di/mocks"
	"goflow/reflection"
	"goflow/tests/assert"
	"testing"
)

func TestValidateDescriptorSuccess(t *testing.T) {
	d := &descriptor{
		interfaceType: reflection.TypeOf[fmt.Stringer](),
		serviceType:   reflection.TypeOf[mocks.StringerMock](),
		factory:       func(provider ServiceProvider) (interface{}, error) { return mocks.NewStringerMock(""), nil },
		life:          Singleton,
	}

	err := validateDescriptor(d)

	assert.NilError(t, err)
}

func TestValidateDescriptorInvalid(t *testing.T) {
	cases := []struct {
		descr      ServiceDescriptor
		errContent string
	}{
		{
			descr: &descriptor{
				interfaceType: reflection.TypeOf[ServiceCollection](),
				serviceType:   reflection.TypeOf[mocks.StringerMock](),
				factory:       func(provider ServiceProvider) (interface{}, error) { return mocks.NewStringerMock(""), nil },
				life:          Singleton,
			},
			errContent: "not assignable",
		},
		{
			descr: &descriptor{
				interfaceType: reflection.TypeOf[mocks.StringerMock](),
				serviceType:   reflection.TypeOf[mocks.StringerMock](),
				factory:       func(provider ServiceProvider) (interface{}, error) { return mocks.NewStringerMock(""), nil },
				life:          Singleton,
			},
			errContent: "is not interface",
		},
		{
			descr: &descriptor{
				interfaceType: reflection.TypeOf[fmt.Stringer](),
				serviceType:   reflection.TypeOf[fmt.Stringer](),
				factory:       func(provider ServiceProvider) (interface{}, error) { return mocks.NewStringerMock(""), nil },
				life:          Singleton,
			},
			errContent: "is not struct",
		},
		{
			descr: &descriptor{
				interfaceType: reflection.TypeOf[fmt.Stringer](),
				serviceType:   reflection.TypeOf[mocks.StringerMock](),
				factory:       func(provider ServiceProvider) (interface{}, error) { return mocks.NewStringerMock(""), nil },
				life:          5,
			},
			errContent: "not valid Life",
		},
	}

	for _, test := range cases {
		err := validateDescriptor(test.descr)

		assert.ErrorWith(t, err, test.errContent)
	}
}

func TestNewServiceDescriptorSuccess(t *testing.T) {
	defer func() {
		r := recover()

		assert.Equal(t, nil, r)
	}()

	d := NewServiceDescriptor[fmt.Stringer, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return nil, nil
	})

	assert.NotNil(t, d)
}

func TestNewServiceDescriptorInvalid(t *testing.T) {
	defer func() {
		r := recover()

		assert.NotNil(t, r)
	}()

	d := NewServiceDescriptor[mocks.StringerMock, mocks.StringerMock](Singleton, func(provider ServiceProvider) (interface{}, error) {
		return nil, nil
	})

	assert.Equal(t, nil, d)
}
