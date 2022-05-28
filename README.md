# Go-flow
 Backend framework for services in go with zero dependency rule, so you can use it in any project without other third-party dependencies or writing your own code for common tasks.

- [Go-flow](#go-flow)
  - [configs](#configs)
  - [di](#di)
  - [logging](#logging)
  - [middleware](#middleware)
  - [observability](#observability)
  - [reflection](#reflection)

## configs
TBA

## di
Dependency injection module with container

```go
// newStringerImplementation is constructor with dependencies
func newStringerImplementation(s Service) *mocks.StringerMock {
	return &mocks.StringerMock{service: s}
}

// newStringerImplementation is factory constructor
func newFactoryStringerImplementation(provider di.ServiceProvider) (*mocks.StringerMock, error) {
    s, err := di.GetService[StringerImplementation](provider)
    if err != nil {
        return nil, err
    }
 	return &mocks.StringerMock{service: s}, nil
}


// registration container for dependency injection
services := di.NewServiceCollection()

err := di.RegisterSingleton[fmt.Stringer, StringerImplementation](services, newStringerImplementation)
// alternatively use:
// err := RegisterSingletonWithFactory[fmt.Stringer](services, newFactoryStringerImplementation)
// or you can use:
// err := RegisterSingletonWithInstance[fmt.Stringer](services, &mocks.StringerMock{})
if err != nil {
     // invalid constructor or type mismatch
    panic(err)
}

// provider container for retrieving services
provider := di.BuildProvider(services)

stringer, err := di.GetService[fmt.Stringer](provider)
// alternatively use:
// var stringer fmt.Stringer
// err := di.Provide(provider, &stringer)
if err != nil {
    // service was not registered or constructor panics
    panic(err)
}
```

## logging
Extended package for standard "log" package

TODO:
- [ ] custom output option 
- [ ] formatter for WithLogFormat 

```go
// settings for global logger
logging.WithOptions(
    logging.WithLogFormat("[{level}]: %v") // this is default format for each level which does not have custom format
    logging.WithInfoFormat("[Info]: %v") // custom format for info level
    logging.WithWarnFormat("[Warn]: %v") // custom format for warn level
    logging.WithErrorFormat("[Error]: %v") // custom format for error level
)

logging.Info("something happened")
logging.Warn("should not happen often")
logging.Error(errors.New("something went wrong"))

// creates new instance of logger
logger := logging.NewLogger(logging.WithLogFormat("[%level]: %v"))

logging.LogInfo(logger, "something happened")
logging.LogWarn(logger, "should not happen often")
logging.LogError(logger, errors.New("something went wrong"))

// log with defined level and custom format
logger.Log(logging.Errorl, "[%v]: %v", "ERROR", errors.New("something went wrong"))
```

## middleware
Generic middleware pipeline pattern for any request and responses
```go
type pipeRequest string
type pipeResponse error

// middleware pipeline for request of type 'pipeRequest' and response of type 'pipeResponse'
middleware := NewMiddleware[pipeRequest, pipeResponse]()

middleware.Use(
	func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
        logging.Info("request started")
		response := next(r)
        if response != nil {
            logging.Error("request failed")
        }
		return response
	},
)

middleware.Use(
	func(r pipeRequest, next func(r pipeRequest) pipeResponse) pipeResponse {
        ok := validate(r)
        if !ok {
            // stop pipeline and return error
            return errors.New("validation failed")
        }
		return next(r)
	},
)

handler := func(r pipeRequest) pipeResponse {
	// actual handler
	return nil
}

// wrap middleware to handler 
wrappedHandler := middleware.Wrap(handler)
request := pipeRequest("request")

// run pipeline
response := wrappedHandler(request)
logging.Info(response)
```

## observability
TBA

## reflection
Extended package for standard "reflect" package

```go
// Types
strTyp := reflection.TypeOf[fmt.Stringer]()
logging.Info(strTyp) // reflect.Type for fmt.Stringer

// Casts
arr := []interface{}{"1", "2", "3"}

stringArr, ok := cast.As[string](arr)
if !ok {
    panic("cannot cast between provided two types")
}
```
