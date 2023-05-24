# How to contribute

You're free to contribute with us! Read the instructions before doing so.

## Prerequisites

- Use the search function to check if the issue or pull request already exists to avoid duplication.
- Start an issue or discussion about a problem or new feature. The decision on whether the feature is suitable for this library should be analyzed and approved by the community.
- In the discussion, express your interest in working on it and start writing code.
- Before submitting a pull request, ensure that your work includes unit (or other) tests and compiles successfully.

## Templates

Please respect the GitHub templates, which provide descriptive information related to issues or pull requests. If something in the template is irrelevant, remove it instead of leaving it blank. Always use a clear and concise title for your request to allow for easier identification of the contents.

## Coding convention

Make sure your code follows the code style of this project. Always format your code according to the shared .editorconfig file in the repository.

### Testing

The approach for tests is as follows:
- Create an `_test.go` file, and the package name should be `{pkg}_test`, where `{pkg}` is the name of the tested package. This ensures that unit tests cover the exported behavior and all success/failure paths are from an exported point.

NOTE: This rule can be broken in certain circumstances, such as important smoke tests or testing an algorithm that requires testing for unexported behavior. In such situations, it's allowed to break it.

- Ensure that every interface has a corresponding mock. Mocks are located in the `tests/mocks/{pkg}.go` file, where `{pkg}` is the package name to which the mock belongs. The mock should be suffixed with `Mock`, and all mocked functions should be named `On{func_name}`. At the top of the file, declare a discarded variable with a type assertion for the mocking type to force the compiler to throw a compile-time error when the mock is not changed accordingly.

- Aim for a minimum coverage of 80%, but it's highly recommended to cover as much code as possible. This is an open-source project that does not generate revenue, so there are no deadlines for features. Therefore, covering rare or obvious blocks of code is encouraged.

## Pull request

Before attempting to make any changes, make sure you are familiar with the Stable Release Strategy described in the [ARCHITECTURE.md](ARCHITECTURE.md) file. The document outlines the strategy for creating stable solutions and the agreement on architecture decisions to provide high-quality software solutions with pull requests.

To have your pull request reviewed, it must adhere to the coding convention, be properly described, and it's recommended to discuss it earlier. Architecture decisions must be made and logged before or during the pull request.
Most things should be automatically checked with CI jobs like build, linter, or test coverage. So, make sure to pass the checks before submitting your pull request from the draft state to ready for review.

### Branch convention

The branch should be named with prefixes indicating the type and related package (*only the "docs" prefix does not require a package prefix*). If the branch is related to changes in more than one package, the most important one (or the one that will have the most changes) should be part of the prefix.

Branch format with examples:

- `{type}/{package?}/{name}`
  - `docs/contributing-update`
  - `adr/config/add-context`
  - `poc/config/add-context`
  - `feat/config/refactor-loading`
  - `fix/config/bug-with-context`

Branch prefix types:
- `docs/` - indicates that the changes are related strictly to documentation only or repository/GitHub files.
- `adr/` - indicates that the changes are related to architecture record documents, which are associated with either a new feature or a significant change.
- `poc/` - indicates that the changes are related to another not-yet-merged `adr/` branch in order to demonstrate it in the code. The suffix should be exactly the same as in the `adr/{suffix}` branch. This type of branch is not meant to be merged into the main branch.
- `feat/` - indicates that the changes are related to a new feature described by another merged `adr/` branch. The suffix should be exactly the same as in the `adr/{suffix}` branch.
- `fix/` - indicates bug fix.
