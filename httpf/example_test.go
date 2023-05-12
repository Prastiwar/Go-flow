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
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/Prastiwar/Go-flow/datas"
	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

// runServer runs new server and returns its address and close function.
func runServer(router httpf.Router) (string, func()) {
	server := httptest.NewServer(router)
	return server.URL, server.Close
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

	serverAddress, cleanup := runServer(mux.Build())
	defer cleanup()

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
		OnLimit: func(ctx context.Context, key string) (rate.Limiter, error) {
			return mocks.LimiterMock{
				OnLimit:  func() uint64 { return maxTokens },
				OnTokens: func(ctx context.Context) (uint64, error) { return tokens, nil },
				OnTake: func(ctx context.Context) (rate.Token, error) {
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
					}, nil
				},
			}, nil
		},
	}

	mux.Get("/api/test/", httpf.RateLimitMiddleware(httpf.HandlerFunc(func(w httpf.ResponseWriter, r *http.Request) error {
		return w.Response(http.StatusOK, nil)
	}), storeLimiter, httpf.PathRateKey()))

	serverAddress, cleanup := runServer(mux.Build())
	defer cleanup()

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

type DummyJsonProducts struct {
	Products []DummyJsonProduct `json:"products"`
}

type DummyJsonProduct struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Price int    `json:"price"`
}

type HttpErr struct {
	Message string `json:"message"`
}

func (err *HttpErr) Error() string {
	return err.Message
}

type DummyJsonClient struct {
	unmarshaler httpf.BodyUnmarshaler
}

func NewDummyJsonClient() *DummyJsonClient {
	return &DummyJsonClient{
		unmarshaler: httpf.NewBodyUnmarshalerWithError(datas.Json(), &HttpErr{}),
	}
}

func (c *DummyJsonClient) GetProducts() (*DummyJsonProducts, error) {
	resp, err := http.Get("https://dummyjson.com/products")
	if err != nil {
		return nil, err
	}

	var data DummyJsonProducts
	if err := c.unmarshaler.Unmarshal(resp, &data); err != nil {
		if _, ok := err.(*HttpErr); ok {
			// handle api error
			return nil, err
		}
		return nil, err
	}

	return &data, nil
}

func (c *DummyJsonClient) GetProduct(id int) (*DummyJsonProduct, error) {
	resp, err := http.Get("https://dummyjson.com/products/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	var data DummyJsonProduct
	if err := c.unmarshaler.Unmarshal(resp, &data); err != nil {
		if _, ok := err.(*HttpErr); ok {
			// handle api error
			return nil, err
		}
		return nil, err
	}

	return &data, nil
}

func ExampleNewBodyUnmarshaler() {
	client := NewDummyJsonClient()

	data, err := client.GetProducts()
	if err != nil {
		panic(err)
	}

	product, err := client.GetProduct(data.Products[0].ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Products[0] == *product)

	// output:
	// true
}
