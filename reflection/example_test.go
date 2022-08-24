package reflection_test

import (
	"fmt"
	"goflow/reflection"
)

func ExampleTypeOf() {
	strTyp := reflection.TypeOf[fmt.Stringer]()

	fmt.Println(strTyp)

	// Output:
	// fmt.Stringer
}

func ExampleParse() {
	parsed, err := reflection.Parse("1", int8(0))
	if err != nil {
		panic(err)
	}

	asInt, ok := parsed.(int8)
	fmt.Println(asInt, ok)

	// Output:
	// 1 true
}
