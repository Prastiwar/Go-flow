package di_test

import (
	"fmt"
	"goflow/di"
)

type Dependency interface{}
type someDependency struct{} // implements Dependency

func NewSomeDependency() *someDependency {
	return &someDependency{}
}

type SomeInterface interface{}
type someService struct { // implements SomeInterface
	serv Dependency
}

func NewSomeService(serv Dependency) *someService {
	return &someService{
		serv: serv,
	}
}

func Example() {
	// register constructors for services and dependencies. By default all services are transient
	container, err := di.Register(
		NewSomeService,
		NewSomeDependency,
	)

	// alternatively you can setup lifetime
	// container, err := di.Register(
	//     di.Construct(di.Singleton, NewSomeService),
	//     di.Construct(di.Transient, NewSomeDependency),
	//     di.Construct(di.Scoped, NewSomeDependency),
	// )

	if err != nil {
		// each ctor must be func Kind with single output parameter
		panic(err)
	}

	// di.Register() already calls this validation - you don't need to do it again
	err = container.Validate()
	if err != nil {
		// any service couldn't be created due to missing or cyclic dependencies
		panic(err)
	}

	// use container.Scope() to create new scoped container which caches scoped services
	// scopedContainer := container.Scope()

	var s SomeInterface
	fmt.Println(s == nil)

	// panics when there is not service implementing SomeInterface
	container.Provide(&s)

	fmt.Println(s == nil)

	// Output:
	// true
	// false
}
