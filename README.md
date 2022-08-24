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
  - [reflection](#reflection)
  - [Contributing](#contributing)
  - [License](#license)

## Library purpose

The idea is to provide and maintain by community single framework without other third-party dependencies to facilitate software development without worrying and dealing with obsolete libraries which hugely increases technical debt. This framework's mission is to extend the built-in GO standard library in a non-invasive way with common systems like configuration, logging and dependency management meaning it should have feeling like it's part of standard one but it should not give up on simplifying building systems by adding GOs like boilerplate.
Writing production-ready system developer often must make decision which will not change and will not apply to every possible case but still should be modifiable enough to make development easier not harder.

## config

Loading configuration from file, environment variables and command line arguments with binding functionality.

See [example file](config\example_test.go) for runnable examples.

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

## reflection

Extended package for standard "reflect" package

See [example file](reflection\example_test.go) for runnable reflection examples.

See [example file](reflection\cast\example_test.go) for runnable casting examples.

## Contributing

You can freely contribute to this library! Report issues and make pull requests to help us improve this project.
Please read [CONTRIBUTING.md](https://github.com/Prastiwar/Go-Flow/blob/main/.github/CONTRIBUTING.md) for details before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/Prastiwar/Go-Flow/blob/main/LICENSE) file for details.
