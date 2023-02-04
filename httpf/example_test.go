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
	"strconv"
	"sync"
	"time"

	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

const (
	hostPrefix = "localhost:"
)

var (
	port     = 8080
	serverMu = sync.Mutex{}
)

// runServer creates new server and listens in new goroutine on address that is returned from this function.
func runServer(router httpf.Router) string {
	serverMu.Lock()
	defer serverMu.Unlock()

	port++
	address := hostPrefix + strconv.Itoa(port)
	go func() {
		_ = httpf.NewServer(address, router).ListenAndServe()
	}()
	return "http://" + address
}

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

	mux.Post("/api/test/", httpf.HandlerFunc(func(w httpf.ResponseWriter, r *http.Request) error {
		result := struct {
			Id string `json:"id"`
		}{
			Id: "1234",
		}

		return w.Response(http.StatusCreated, result)
	}))

	serverAddress := runServer(mux.Build())

	resp, err := http.Post(serverAddress+"/api/test/", httpf.ApplicationJsonType, bytes.NewBufferString("{}"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// output:
	// 201
	// {"id":"1234"}
}

func ExampleRateLimitMiddleware() {
	mux := httpf.NewServeMuxBuilder()

	mux.WithErrorHandler(httpf.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, rate.ErrRateLimitExceeded) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		panic(err)
	}))

	var (
		resetPeriod time.Duration = time.Second
		maxTokens   uint64        = 2
		tokens      uint64        = maxTokens
	)
	storeLimiter := mocks.LimiterStoreMock{
		OnLimit: func(key string) rate.Limiter {
			return mocks.LimiterMock{
				OnLimit:  func() uint64 { return maxTokens },
				OnTokens: func() uint64 { return tokens },
				OnTake: func() rate.Token {
					return mocks.TokenMock{
						OnUse: func() error {
							if tokens <= 0 {
								return rate.ErrRateLimitExceeded
							}
							tokens--
							go func() {
								time.Sleep(resetPeriod)
								tokens++
							}()
							return nil
						},
						OnResetsAt: func() time.Time { return time.Now().Add(resetPeriod) },
					}
				},
			}
		},
	}

	mux.Get("/api/test/", httpf.RateLimitMiddleware(httpf.HandlerFunc(func(w httpf.ResponseWriter, r *http.Request) error {
		return w.Response(http.StatusOK, nil)
	}), storeLimiter, httpf.PathRateKey()))

	serverAddress := runServer(mux.Build())

	callGet := func() {
		resp, err := http.Get(serverAddress + "/api/test/")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(resp.StatusCode)
	}

	for i := uint64(0); i <= maxTokens; i++ {
		callGet()
	}

	time.Sleep(resetPeriod)
	callGet()

	// output:
	// 200
	// 200
	// 429
	// 200
}
