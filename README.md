# Go-flow

[![Go Reference](https://pkg.go.dev/badge/github.com/Prastiwar/go-flow.svg)](https://pkg.go.dev/github.com/Prastiwar/Go-flow)
[![Go Report Card](https://goreportcard.com/badge/github.com/Prastiwar/go-flow)](https://goreportcard.com/report/github.com/Prastiwar/Go-flow)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Prastiwar_Go-flow&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Prastiwar_Go-flow)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=Prastiwar_Go-flow&metric=coverage)](https://sonarcloud.io/summary/new_code?id=Prastiwar_Go-flow)

Backend framework library written in Go for Gophers who want to focus on building great quality backend. It mostly provides dependency inversion to most common built systems to allow easily modularize our application and change implementation without changing our business logic. This is very cost-effective for organizations moving forward from MVP. The framework is viable for both personal and organization level usage where there are multiple products, APIs and development consistency is needed.

- [Go-flow](#go-flow)
  - [Download](#download)
  - [Library purpose](#library-purpose)
  - [Packages](#packages)
    - [caching](#caching)
    - [config](#config)
    - [datas](#datas)
    - [di](#di)
    - [exception](#exception)
    - [httpf](#httpf)
    - [logging](#logging)
    - [middleware](#middleware)
    - [policy](#policy)
      - [retry](#retry)
    - [rate](#rate)
    - [reflection](#reflection)
    - [tests](#tests)
  - [Contributing](#contributing)
  - [License](#license)

## Download

See latest version: https://github.com/Prastiwar/Go-flow/releases/latest

```bash
go get -u github.com/Prastiwar/Go-flow@{version}
```

## Library purpose

The idea is to provide and maintain by community single framework with the most used backend systems like caching and logging without the need for repeatable writing the same interfaces over third-party dependencies to facilitate software development with dependency inversion rule. It does provide an implementation for some systems so Gophers would not waste time researching GitHub and worrying about massive third-party obsolete libraries which hugely increases technical debt and takes time to find the one that fits your needs.
This does not mean the framework is completely resilient to time and language development. A single library is simpler to maintain by the community to prevent deprecating the library or forcing developer to wait for changes.
The technical debt is the reason why this library follows the no-dependency rule which means it does not depend on any other library. It can be visible in clean go.mod file.
This framework's mission is to extend the built-in GO standard library in a non-invasive way with common systems like configuration, logging, and dependency management meaning it should have a feeling like it's part of standard one but it should not give up on simplifying building systems by adding boilerplate.
Writing production-ready system developers often must make a decision that will not change and will not apply to every possible case but still should be modular enough to make development easier not harder.

## Packages

### caching

Caching package introduces dependency inversion over third-party cache libraries. In case the third party does not follow the contract your infrastructure should implement an adapter for this specific library to fulfill Cache interface. In-memory caching interface requires a minimum of TTL support and any storing acceptance for any value. If a third-party package does accept only byte slice (e.g [bigcache](https://github.com/allegro/bigcache)) as a value it would need to make use of encoding/decoding internally to fulfill a contract.

### config

Configuration module which provides functionality to load configuration from file, environment variables and command line arguments with binding to a struct functionality.
It allows to extend the behavior with interfaces for providers and KeyInterceptor option to change the way it looks for matching key for field name.

See [example file](config/example_test.go) for runnable examples.

### datas

Datas is a package that provides functionality for data formatting, encoding, and decoding especially for serialization purposes. It includes support for JSON and XML encoding, as well as generic interfaces that allow you to work with data in a flexible and extensible way. This package is particularly useful for projects that need to work with data in a variety of formats, or that require a high degree of customization in how data conversion is handled. The package will be highly appreciated by people who plan to use third-party library to convert data using standard format, like json instead standard library for performance(or any other) reason.

See [example file](datas/example_test.go) for runnable examples.

### di

Dependency injection module with container. This pattern is encouraged to use in large projects where dependency hierarchy is deep and complex and cannot be improved by design decisions. In such case dependency maintenance can be a problem that container can solve.
It's not recommended to use it in small or medium projects where dependency graph is simple and could be improved by design decisions.
Try to use dependency injection without container first and then use container if you really need it.
Providing a service implementation does not return error - it panics instead. User is responsible for verifying if service he wants to use is registered - this is the easiest problem user need to deal with
while working with dependency container. The other common mistakes like cyclic dependency or missing dependency is solved by validating the container registration and returning and error at this point.

See [example file](di/example_test.go) for runnable examples.

### exception

It provides helper functions to facilitate work with errors. It allows to handle panic with ensured error (when panic is commonly mixed strings or errors), aggregate the errors and more.

See [example file](exception/example_test.go) for runnable examples.

### httpf

httpf package provides abstraction over standard net/http to introduce dependency inversion rule. Mosly routing and server are abstracted which should help with mocking and facilitate using it without mistakes while providing harder to misuse API.
Additionaly it adds simple configurable rate limiter middleware for request per IP or Endpoint.

See [example file](httpf/example_test.go) for runnable examples.

### logging

Logf package is very simple wrapper over io.Writer with provided Formatter and scope(Fields) added. It provides leveling printing as Info, Error and Debug.
As far as you could add custom level, because it's just a field variable, it's not recommended until it's necessary.
The default three: info, error and debug matches almost every project where Trace is used in very much the same matter as Debug one.
Warning is not common and mostly used in wrong way since it's just an information, so info level can be used instead. It does not contain well known global logging as in standard or other libraries
which should be considered as anti-pattern due to hidden dependency and lack of way to encapsulate(mock) its behavior.
The global printing should be used only for testing or playground cases where fmt.Print suits best but due to global logging existence it often encourages to use it in actual project where correct way would be to use logger as dependency.

See [example file](logf/example_test.go) for runnable examples.

### middleware

Generic middleware implements middleware architecture pattern to use in generic way. It's used to create delegates for processing the request or response and handle common tasks like
logging, authentication, compressing data in single contact point which is called pipeline. The pattern is commonly used in http packages.

See [example file](middleware/example_test.go) for runnable examples.

### policy

Defines policies that could help developers to handle gracefully faults like retry or any other resilience functionality.

#### retry

Policy helps to handle transient errors by repeating the function call. It includes configuration features like retry count, wait time before next retry execution or cancellation control which can be used to stop retry execution on error which is not transient.

See [example file](policy/retry/example_test.go) for runnable retry policy examples.

### rate

Rate package introduces dependency inversion over third-party rate-limit libraries. In case the third party does not follow the contract your infrastructure should implement an adapter for this specific library to fulfill Limiter interfaces.
Package contains simple and unified API that can be implemented by any algorithm like token/leaky bucket, fixed/sliding window, and others. There are a few interfaces that extend the standard Limiter interface like BurstLimiter which supports bursting or ReservationLimiter which supports token reservation.

See [example file](rate/example_test.go) for runnable rate examples.

### reflection

Extended package for standard "reflect" package. Provides functions to help with tasks where reflection is needed like parsing, setting field values or casting an array.

See [example file](reflection/example_test.go) for runnable reflection examples.

See [example file](reflection/cast/example_test.go) for runnable casting examples.

### tests

Package contains assertions for equality, matching slice, map elements and counter for asserting function call. The amount of assert features are very limited due to its convention to provide only mostly used functions based on this particular project. It can be extended with any additional reasonable assertions in future. Currently the most used library is [github.com/stretchr/testify](https://github.com/stretchr/testify/) which I'd recommend read about if your project requires complex and lots more features in assertions or mocking.

See any test files to discover usage.

## Contributing

You can freely contribute to this library! Report issues and make pull requests to help us improve this project.
Please read [CONTRIBUTING.md](https://github.com/Prastiwar/Go-Flow/blob/main/.github/CONTRIBUTING.md) for details before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/Prastiwar/Go-Flow/blob/main/LICENSE) file for details.
