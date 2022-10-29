package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	stdHttp "net/http"

	"github.com/Prastiwar/Go-flow/rest"
	"github.com/Prastiwar/Go-flow/rest/http"
)

func ExampleNewFluentRouter() {
	rawRouter := http.NewHttpRouter()
	router := rest.NewFluentRouter(rawRouter)
	router.RegisterFunc("/api/test", func(req rest.HttpRequest) rest.HttpResponse {
		fmt.Println("printed from api endpoint")
		return rest.Ok()
	})

	server := &stdHttp.Server{}
	s := http.NewServer(server, router)

	go func() {
		_ = s.Run("localhost:8080", router)
	}()

	stdClient := &stdHttp.Client{}
	rawClient := http.NewHttpClient(stdClient)
	client := rest.NewFluentClient(rawClient)

	resp, err := client.Get(context.Background(), "http://localhost:8080/api/test")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() >= 400 {
		panic(fmt.Errorf("unsuccessfull http request: %v", resp))
	}

	if err := s.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	// Output:
	// printed from api endpoint
}

func ExampleNewFluentClient() {
	stdClient := &stdHttp.Client{}
	rawClient := http.NewHttpClient(stdClient)
	client := rest.NewFluentClient(rawClient)

	resp, err := client.Get(context.Background(), "https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		panic(err)
	}

	type Result struct {
		Id     int `json:"id"`
		UserId int `json:"userId"`
	}
	var result Result

	decoder := json.NewDecoder(resp.Body())
	if err := decoder.Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Id:", result.Id)
	fmt.Println("User Id:", result.UserId)

	// Output:
	// Id: 1
	// User Id: 1
}
