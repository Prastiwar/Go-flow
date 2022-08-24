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
  - [Contributing](#contributing)
  - [License](#license)

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

See [example file](di\example_test.go) for runnable examples.

## logging

Extended package for standard "log" package.

See [example file](logf\example_test.go) for runnable examples.

## middleware

Generic middleware pipeline pattern for any request and responses

See [example file](middleware\example_test.go) for runnable examples.

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

## Contributing

You can freely contribute to this library! Report issues and make pull requests to help us improve this project.
Please read [CONTRIBUTING.md](https://github.com/Prastiwar/Go-Flow/blob/main/.github/CONTRIBUTING.md) for details before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/Prastiwar/Go-Flow/blob/main/LICENSE) file for details.
