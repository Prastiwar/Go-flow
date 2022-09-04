# Go-flow

[![Go Reference](https://pkg.go.dev/badge/github.com/Prastiwar/go-flow.svg)](https://pkg.go.dev/github.com/Prastiwar/go-flow)

Framework for Go services in go with zero dependency rule, so you can use it in any project without other third-party dependencies or writing your own code for common tasks.

- [Go-flow](#go-flow)
  - [Library purpose](#library-purpose)
  - [Packages](#packages)
    - [config](#config)
    - [di](#di)
    - [exception](#exception)
    - [logging](#logging)
    - [middleware](#middleware)
    - [reflection](#reflection)
    - [tests](#tests)
  - [Contributing](#contributing)
  - [License](#license)

## Library purpose

The idea is to provide and maintain by community single framework without other third-party dependencies to facilitate software development without worrying and dealing with obsolete libraries which hugely increases technical debt.
The technical debt is the reason why this library follows no dependency rule which means it does not depend on any other library. It can be clearly visible in go.mod file.
This framework's mission is to extend the built-in GO standard library in a non-invasive way with common systems like configuration, logging and dependency management
meaning it should have feeling like it's part of standard one but it should not give up on simplifying building systems by adding GOs like boilerplate.
Writing production-ready system developer often must make decision which will not change and will not apply to every possible case but still should be modifiable enough to make development easier not harder.

## Packages

### config

Configuration module wchich provides functionality to load configuration from file, environment variables and command line arguments with binding to a struct functionality.
It allows to extend the behavior with interfaces for providers and KeyInterceptor option to change the way it looks for matching key for field name.

See [example file](config\example_test.go) for runnable examples.

### di

Dependency injection module with container. This pattern is encouraged to use in large projects where dependency hierarchy is deep and complex and cannot be improved by design decisions.
In such case dependency maintenance can be a problem that container can solve.
It's not recommended to use it in small or medium projects where dependency graph is simple and could be improved by design decisions.
Try to use dependency injection without container first and then use container if you really need it.
Providing a service implementation does not return error - it panics instead. User is responsible for verifying if service he wants to use is registered - this is the easiest problem user need to deal with
while working with dependency container. The other common mistakes like cyclic dependency or missing dependency is solved by validating the container registration and returning and error at this point.

See [example file](di\example_test.go) for runnable examples.

### exception

It provides helper functions to facilitate work with errors. It allows to handle panic with ensured error (when panic is commonly mixed strings or errors), aggregate the errors and more.

See [example file](exception\example_test.go) for runnable examples.

### logging

Logf package is very simple wrapper over io.Writer with provided Formatter and scope(Fields) added. It provides leveling printing as Info, Error and Debug.
As far as you could add custom level, because it's just a field variable, it's not recommended until it's necessary.
The default three: info, error and debug matches almost every project where Trace is used in very much the same matter as Debug one.
Warning is not common and mostly used in wrong way since it's just an information, so info level can be used instead. It does not contain well known global logging as in standard or other libraries
which should be considered as anti-pattern due to hidden dependency and lack of way to encapsulate(mock) its behavior.
The global printing should be used only for testing or playground cases where fmt.Print suits best but due to global logging existence it often encourages to use it in actual project where correct way would be to use logger as dependency.

See [example file](logf\example_test.go) for runnable examples.

### middleware

Generic middleware implements middleware architecture pattern to use in generic way. It's used to create delegates for processing the request or response and handle common tasks like
logging, authentication, compressing data in single contact point which is called pipeline. The pattern is commonly used in http packages.

See [example file](middleware\example_test.go) for runnable examples.

### reflection

Extended package for standard "reflect" package. Provides functions to help with tasks where reflection is needed like parsing, setting field values or casting an array.

See [example file](reflection\example_test.go) for runnable reflection examples.

See [example file](reflection\cast\example_test.go) for runnable casting examples.

### tests

Package contains assertions for equality, matching slice, map elements and counter for asserting function call. The amount of assert features are very limited due to its convention to provide only mostly used functions based on this particular project. It can be extended with any additional reasonable assertions in future. Currently the most used library is [github.com/stretchr/testify](https://github.com/stretchr/testify/) which I'd recommend read about if your project requires complex and lots more features in assertions or mocking.

See any test files to discover usage.

## Contributing

You can freely contribute to this library! Report issues and make pull requests to help us improve this project.
Please read [CONTRIBUTING.md](https://github.com/Prastiwar/Go-Flow/blob/main/.github/CONTRIBUTING.md) for details before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/Prastiwar/Go-Flow/blob/main/LICENSE) file for details.
