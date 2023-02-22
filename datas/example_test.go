package datas_test

import (
	"fmt"

	"github.com/Prastiwar/Go-flow/datas"
)

func Example() {
	var serializer datas.Marshaler = datas.Json()

	type Data struct {
		Foo string `json:"foo"`
	}

	testData := Data{Foo: "success"}

	b, err := serializer.Marshal(testData)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	// Output:
	// {"foo":"success"}
}
