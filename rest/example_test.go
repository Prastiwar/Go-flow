package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	stdHttp "net/http"

	"github.com/Prastiwar/Go-flow/rest"
	"github.com/Prastiwar/Go-flow/rest/http"
)

func Example() {
	stdClient := &stdHttp.Client{}
	rawClient := http.NewHttpClient(stdClient)
	client := rest.NewFluentClient(rawClient)

	resp, err := client.Get(context.Background(), "https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		// cannot send request
		panic(err)
	}

	type Result struct {
		Id     int `json:"id"`
		UserId int `json:"userId"`
	}
	var result Result

	decoder := json.NewDecoder(resp.Body())
	if err := decoder.Decode(&result); err != nil {
		// cannot decode
		panic(err)
	}

	fmt.Println("Id:", result.Id)
	fmt.Println("User Id:", result.UserId)

	// Output:
	// Id: 1
	// User Id: 1
}
