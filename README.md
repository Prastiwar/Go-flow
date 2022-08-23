# Go-flow

[![Go Reference](https://pkg.go.dev/badge/github.com/Prastiwar/go-flow.svg)](https://pkg.go.dev/github.com/Prastiwar/go-flow)

Framework for Go services in go with zero dependency rule, so you can use it in any project without other third-party dependencies or writing your own code for common tasks.
 
- [Go-flow](#go-flow)
  - [Library purpose](#library-purpose)
  - [config](#config)
      - [TODO](#todo)
  - [di](#di)
  - [logging](#logging)
  - [middleware](#middleware)
  - [observability](#observability)
  - [reflection](#reflection)
 
## Library purpose

The idea is to provide and maintain by community single framework without other third-party dependencies to facilitate software development without worrying and dealing with obsolete libraries which hugely increases technical debt. This framework's mission is to extend the built-in GO standard library in a non-invasive way with common systems like configuration, logging and dependency management meaning it should have feeling like it's part of standard one but it should not give up on simplifying building systems by adding GOs like boilerplate.
Writing production-ready system developer often must make decision which will not change and will not apply to every possible case but still should be modifiable enough to make development easier not harder.

## config

Loading configuration from file, environment variables and command line arguments with binding functionality.
```go
// Provide creates new Source instance with provided configs.
cfg := config.Provide(
    // { "queryTimeout": "10s" }
    config.NewFileProvider("config.json", decoders.NewJson()), 
    // --dbName="my-collection" --errorDetails=true
    config.NewFlagProvider(
        config.StringFlag("dbName", "name for database"),
        config.BoolFlag("errorDetails", "should show error details"),
    ),
    // CONNECTION_STRING="mongodb://localhost:8089"; ERROR_DETAILS="false"
    config.NewEnvProvider(),
)

// ShareOptions provides options to set default options for all provider.
cfg.ShareOptions(
    // KeyInterceptor allows to intercept field name before it's used to find it in provider.
    // It's useful when you want to use different field names than they're defined in struct.
    // For example you can use tag to define field name and intercept it there.
    WithKeyInterceptor(func(providerName string, field reflect.StructField) string {
        if providerName == config.EnvProviderName {
            return strings.ToUpper(field.Name)
        }
        return field.Name
    })
)

// Use default values for options in case they are not included in providers.
cfg.SetDefault(
    config.Opt("connectionString", "mongodb://localhost:27017"),
    config.Opt("dbName", "go-flow"),
    config.Opt("errorDetails", true),
    config.Opt("queryTimeout", time.Second * 15),
    config.Opt("access-key", "ABC123EFGH456IJK789"),
)

type DbOptions struct {
    DbName           string
    ConnectionString string
    ErrorDetails     bool
    QueryTimeout     time.Duration
    AccessKey        string
}

var dbOptions DbOptions
err := cfg.Load(&dbOptions)
// dbOptions were loaded starting from file -> flag -> env and all values were overriden in that order.
// The default value is not overriden if it doesn't exist in any provider
if err != nil {
    // One of the providers failed to load config values
    panic(err)
}
// options.DbName == "my-collection"
// options.ConnectionString == "mongodb://localhost:8089"
// options.ErrorDetails == false
// options.QueryTimeout == time.Second * 10
// settings.AccessKey == "ABC123EFGH456IJK789"

type AccessOptions struct {
    AccessKey string
}
var aOptions AccessOptions
// Bind will try to copy corresponding field from dbOptions to aOptions
err = config.Bind(dbOptions, &aOptions)
if err != nil {
    // Probably field type mismatch
    panic(err)
}
// aOptions.AccessKey == "ABC123EFGH456IJK789"
```

#### TODO

- [ ] cache bindings

## di

Dependency injection module with container. This pattern is encouraged to use in large projects where dependency hierarchy is deep and complex and cannot be improved by design decisions.
It's not recommended to use it in small or medium projects where dependency graph is simple and could be improved by design decisions. 
Use dependency injection without container first and then use container if you really need it.

```go

type Dependency interface {}
type someDependency struct {} // implements Dependency

func NewSomeDependency() *someDependency {
	return &someDependency{}
}

type SomeInterface interface {}
type someService struct { // implements SomeInterface
	serv Dependency
}

func NewSomeService(serv Dependency) *someService {
	return &SomeService{
		serv: serv,
	}
}

// register constructors for services and dependencies
// by default all services are transient
container, err := di.Register(
	NewSomeService,
	NewSomeDependency,
)

// alternatively you can setup lifetime
container, err := di.Register(
	di.Construct(di.Singleton, NewSomeService),
	di.Construct(di.Transient, NewSomeDependency),
	di.Construct(di.Scoped, NewSomeDependency),
)

if err != nil {
    // each ctor must be func Kind with single output parameter
    panic(err)
}

// return error if any service cannot be created due to dependencies
err := container.Validate()

// alternative: s := di.New[SomeInterface]()
var s SomeInterface
// panics when there is not service implementing SomeInterface
container.Provide(&s)

scopedContainer := container.Scope()
scopedContainer.Provide(&scopedService) // will cache this service in this scope
```

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
// Create new logger with optional settiings like output, formatter, scope fields
logger = logf.NewLogger(
    logf.WithOutput(writer),
    logf.WithFormatter(logf.NewTextFormatter()),
    logf.WithFields(logf.Fields{logf.LogTime: logf.NewTimeField(time.RFC3999)}),
)


// Create logger based on parent logger with additional scope
logger = logf.WithScope(
    logger, 
    logf.Fields{
        "currentTime": time.Now().UTC().Format("2006-01-02 15:04:05")
    }
)

logger.Error("error message")
logger.Errorf("error occured: %v", err)

logger.Info("info message")
logger.Infof("count: %v", 1)

logger.Debug("debug message")
logger.Debugf("debug message: %v", "debug message")
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
