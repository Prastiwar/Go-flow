package httpf_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Prastiwar/Go-flow/httpf"
)

func Example() {
	mux := httpf.NewServeMuxBuilder()

	mux.WithErrorHandler(httpf.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		// define standard structure for error response
		type httpError struct {
			Error string `json:"error"`
			Code  int    `json:"code"`
		}

		// map infrastructure errors to http error response
		var resultError httpError
		if errors.Is(err, context.DeadlineExceeded) {
			resultError = httpError{
				Error: http.StatusText(http.StatusRequestTimeout),
				Code:  http.StatusRequestTimeout,
			}
		} else {
			resultError = httpError{
				Error: http.StatusText(http.StatusInternalServerError),
				Code:  http.StatusInternalServerError,
			}
		}

		// marshal error and write to Response
		result, err := json.Marshal(resultError)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Write(result)
		if err != nil {
			log.Fatal(err)
		}
	}))

	mux.Post("/api/test/", httpf.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		err := errors.New("lel")
		if err != nil {
			return err
		}

		result := struct {
			Id string `json:"id"`
		}{
			Id: "1234",
		}

		return httpf.Json(w, http.StatusCreated, result)
	}))

	go func() {
		_ = httpf.NewServer("localhost:8080", mux.Build()).ListenAndServe()
	}()

	resp, err := http.Post("http://localhost:8080/api/test/", "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(body)

	// output:
	// 201
	// {"id":"1234"}
}
