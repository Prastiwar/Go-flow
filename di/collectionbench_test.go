package di

import (
	"fmt"
	"goflow/di/mocks"
	"testing"
)

type Dependency1 interface{}
type Dependency2 interface{}
type Dependency3 interface{}
type Dependency4 interface{}
type Dependency5 interface{}
type Dependency6 interface{}
type Dependency7 interface{}

var singleCtor = func(d Dependency1) *mocks.StringerMock {
	return &mocks.StringerMock{}
}

var sevenCtor = func(d1 Dependency1, d2 Dependency2, d3 Dependency3, d4 Dependency4, d5 Dependency5, d6 Dependency6, d7 Dependency7) *mocks.StringerMock {
	return &mocks.StringerMock{}
}

var singleFac = func(provider ServiceProvider) (*mocks.StringerMock, error) {
	_, err := GetService[Dependency1](provider)
	if err != nil {
		return nil, err
	}
	return &mocks.StringerMock{}, nil
}

var sevenFac = func(provider ServiceProvider) (*mocks.StringerMock, error) {
	_, err := GetService[Dependency1](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency2](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency3](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency4](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency5](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency6](provider)
	if err != nil {
		return nil, err
	}
	_, err = GetService[Dependency7](provider)
	if err != nil {
		return nil, err
	}
	return &mocks.StringerMock{}, nil
}

func BenchmarkRegister1Dep(b *testing.B) {
	runBenchmarkRegister(b, singleCtor)
}

func BenchmarkRegister7Dep(b *testing.B) {
	runBenchmarkRegister(b, sevenCtor)
}

func BenchmarkRegisterWithFactory1Dep(b *testing.B) {
	runBenchmarkRegisterWithFactory(b, singleFac)
}

func BenchmarkRegisterWithFactory7Dep(b *testing.B) {
	runBenchmarkRegisterWithFactory(b, sevenFac)
}

func BenchmarkRegisterWithInstance(b *testing.B) {
	services := NewServiceCollection()
	instance := newStringerImplementation(nil)
	_ = RegisterSingletonWithInstance[fmt.Stringer](services, instance)
	provider := BuildProvider(services)
	descFac := services.Descriptors()[0].Factory()
	for i := 0; i < b.N; i++ {
		_, _ = descFac(provider)
	}
}

func BenchmarkProvide1Dep(b *testing.B) {
	runBenchmarkProvide(b, singleCtor)
}

func BenchmarkProvide7Dep(b *testing.B) {
	runBenchmarkProvide(b, sevenCtor)
}

func BenchmarkProvideWithFactory1Dep(b *testing.B) {
	runBenchmarkProvideWithFactory(b, singleFac)
}

func BenchmarkProvideWithFactory7Dep(b *testing.B) {
	runBenchmarkProvideWithFactory(b, sevenFac)
}

func BenchmarkProvideWithInstance(b *testing.B) {
	services := NewServiceCollection()
	instance := newStringerImplementation(nil)
	_ = RegisterSingletonWithInstance[fmt.Stringer](services, instance)
	provider := BuildProvider(services)
	descFac := services.Descriptors()[0].Factory()
	for i := 0; i < b.N; i++ {
		_, _ = descFac(provider)
	}
}

func runBenchmarkRegister(b *testing.B, fac interface{}) {
	for i := 0; i < b.N; i++ {
		services := NewServiceCollection()
		_ = RegisterSingleton[fmt.Stringer, mocks.StringerMock](services, fac)
	}
}

func runBenchmarkProvide(b *testing.B, fac interface{}) {
	services := NewServiceCollection()
	provider := BuildProvider(services)
	_ = RegisterSingleton[fmt.Stringer, mocks.StringerMock](services, fac)
	descFac := services.Descriptors()[0].Factory()
	for i := 0; i < b.N; i++ {
		_, _ = descFac(provider)
	}
}

func runBenchmarkRegisterWithFactory(b *testing.B, fac func(ServiceProvider) (*mocks.StringerMock, error)) {
	for i := 0; i < b.N; i++ {
		services := NewServiceCollection()
		_ = RegisterSingletonWithFactory[fmt.Stringer](services, fac)
	}
}

func runBenchmarkProvideWithFactory(b *testing.B, fac func(ServiceProvider) (*mocks.StringerMock, error)) {
	services := NewServiceCollection()
	provider := BuildProvider(services)
	_ = RegisterSingletonWithFactory[fmt.Stringer](services, fac)
	descFac := services.Descriptors()[0].Factory()
	for i := 0; i < b.N; i++ {
		_, _ = descFac(provider)
	}
}
