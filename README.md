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

```go
// logging functions with added level field
logf.PrintError("error message")
logf.PrintErrorf("error occured: %v", err)

logf.PrintInfo("info message")
logf.PrintInfof("count: %v", 1)

logf.PrintWarn("warning message")
logf.PrintWarnf("probably should be: %v", 2)

logf.PrintDebug("debug message")
logf.PrintDebugf("debug message: %v", "debug message")

logf.PrintTrace("trace message")
logf.PrintTracef("trace message: %v", "trace message")

logf.PrintFatal("fatal message")
logf.PrintFatalf("fatal message: %v", "fatal message")

// SetFormatter wraps default (global) writer with formatter and sets flags to 0
SetFormatter(formatter)

// SetScope wraps default (global) writer and replaces formatter with additional scope
SetScope(formatter)

// CreateWithScope creates new instance of *log.Logger with provided formatter
logger := CreateWithFormatter(formatter)

// CreateWithScope creates new instance of *log.Logger with provided fields.
// Formatter is preserved or initialized with logf.DefaultFormatter() if not set
logger = CreateWithScope(logf.Fields{"currentTime": time.Now().UTC().Format("2006-01-02 15:04:05")})

logf.Error(logger, "error message")
logf.Errorf(logger, "error occured: %v", err)

logf.Info(logger, "info message")
logf.Infof(logger, "count: %v", 1)

logf.Warn(logger, "warning message")
logf.Warnf(logger, "probably should be: %v", 2)

logf.Debug(logger, "debug message")
logf.Debugf(logger, "debug message: %v", "debug message")

logf.Trace(logger, "trace message")
logf.Tracef(logger, "trace message: %v", "trace message")

logf.Fatal(logger, "fatal message")
logf.Fatalf(logger, "fatal message: %v", "fatal message")

// Creates logger based on parent logger with additional scope
logger = WithScope(logger, logf.Fields{"currentTime": time.Now().UTC().Format("2006-01-02 15:04:05")})

// Creates logger based on parent logger with different formatter
logger = WithFormatter(logger, formatter)
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
