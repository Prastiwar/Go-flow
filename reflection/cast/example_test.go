package cast_test

import (
	"fmt"
	"goflow/reflection/cast"
)

func ExampleAs() {
	arr := []interface{}{"1", "2", "3"}

	stringArr, ok := cast.As[string](arr)
	if !ok {
		panic("cannot cast between provided two types")
	}
	fmt.Println(stringArr)

	// Output:
	// [1 2 3]
}
