# How to contribute

You're free to contribute with us! Read the instructions before do it.

## Prerequisites

- Use search to find if issue/pull request already exists to not duplicate it.
- Start an issue or discussion over problem or new feature. The decision if this is applicable feature for this library should be analyzed and approved by community.
- You can point out in discussion that you want to work on it and start writing code.
- Before submitting pull request make sure your work contains unit (or other) tests and compiles successfully.

## Templates

Please, respect the github templates which were made for descriptive enough information realted with issue/pull request.
If something in template is unrelevant, remove it instead of leaving it blank.
Always use clear and concise title for your request to allow easier identification for the contents.

## Coding convention

Make sure your code follows the code style of this project. Always format your code accordingly to shared .editorconfig in repository.

### Testing

The approach for tests is as follows:
- Make `_test.go` file and the package name should be `{pkg}_test` where `{pkg}` is the name of tested package - this assures unit tests cover the exported behaviour and all success/failure paths are from exported point.

NOTE: This rule can be broken in certain circumstances like important smoke tests or testing alghoritm which would require testing for unexported behaviour. In such situation it's allowed to

- Every interface should have corresponding mock. Mocks are located in the `tests/mocks/{pkg}.go` where `{pkg}` is the package name that mock belongs to. Mock should be suffixed with `Mock` and all mocked functions should be named `On{func_name}`. On top of the file there should be declared discarded variable with type assertion for mocking type to force compiler throw compile-time error when mock was not changed accordingly.

- 80% coverage is absolutely minimum but it's highly recommended to cover as much as possible. This is open-source project which does not make money therefore has no deadlines on features so there is nothing to lose by covering rare or obvious block of code.

## Pull request

In order to review pull request it must match coding convention, be properly described and recommended to be disscussed earlier.
Most of things should be automatically check with CI jobs like build, linter or test coverage so make sure to pass the checks before submitting
your pull request from draft state to ready for review.
