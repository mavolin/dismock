# Contributing

We would love to see the ideas you want to bring in to improve this project.
Before you get started, make sure to read the guidelines below. 

## Contributing through issues

If have an idea how to improve this project or if you found a bug, let us know by submitting an issue.
The issue templates will take care of most of the requirements, but there is one thing you should note:

### Titles

We not only use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) for commits, but also for issue titles.
If you propose a feature, use `feat(PACKAGE_NAME): TITLE`, for a bug replace the `feat` with a `fix`, for `docs` use `docs` etc.
Other types are `style`, `refactor` and `test`.
If your change is breaking (in the semantic versioning sense) add an exclamation mark behind the scope, e.g. `feat(package)!: title`.

If your issue proposes changes to multiple packages, or you don't know which package is affected, leave the `(PACKAGE_NAME)` part out, e.g. `feat: that one thing that's missing`.


## Code Contributions
### Opening an Issue

Before you hand in a PR open an issue describing what you want to change, and tell us you'll be handing in a PR for it.
This gives us the ability to point out important things you should keep in mind, and give you feedback for your idea, before you get to implementing the feature.

### Committing

This is by far the most important guideline.
Please make small, thoughtful commits, a commit like `feat: add xy` with 20 new files is rarely appropriate.

#### Conventional Commits

Please use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) for your contributions, once you get the hang of it, you'll see that they are pretty easy to use.
Just have a look at the [quick start guide](https://www.conventionalcommits.org/en/v1.0.0/#summary) on their website.
The scope is typically the package name, but for non-go files appropriate scopes may also be: `git`, `readme`, `golangci`, or `go.mod`.

##### Types
We use the following types:

- ci: changes to our CI configuration files and scripts
- docs: changes to the documentation
- feat: a new feature
- fix: a bug fix
- perf: an improvement to performance
- refactor: a code change that neither fixes a bug nor adds a feature
- style: a change that does not affect the meaning of the code
- test: a change to an existing test, or a new test

##### Breaking Changes

Breaking changes must have a `!` after the type/scope and a `BREAKING CHANGE:` footer.

### Fixing a Bug

If you're fixing a bug make sure to add a test case for that bug, to ensure it's gone for good.
This of course only applies if the function is testable.

### Code Style

Make sure all code is `gofmt -s`'ed, and passes the golangci-lint checks.
If your code fails a lint task, but the way you did it is justified, add a `//nolint:{{name_of_linter}}` comment to the line or block.

### Testing

If possible and appropriate you should fully test the code you submit.
Each function exposed and unexposed should have a single test, which either tests directly or is split into subtests, preferably table-driven.

In an effort to ease the writing of tests, we use [testify](https://github.com/stretchr/testify) for all testing.

#### Table-Driven Tests

If there is a single table, it should be called `testCases`, multiple use the name `<type>Cases`, e.g. `successCases` and `failureCases`, for tests that test the output for a valid input (a success case), and those that aim to provoke an error (a failure case) and therefore work different from a success case.
The same goes if there is a table that's only testing a portion of a function, and multiple non-table-driven tests in addition.

The structs used in tables should always anonymous.

Every sub-test including table-driven ones should have a name that clearly shows what is being done.
For table-driven tests this name is either obtained from a `name` field or computed using the other fields in the table entry.

Every case in a table should run in its own subtest (`t.Run`).
Additionally, if there are multiple tables, each table should have its own subtest, in which it calls its cases:

```
TestSomething
    testCase
    testCase

    additionalTest
```

```
TestSomething
    successCases
        successCase
        successCase
    failureCases
        failureCase
        failureCase

    additionalTest
```

```
TestSomething
    successCases
        successCase
        successCase
    failureCases
        failureCase
        failureCase
        additionalFailureTest
```

### Opening a Pull Request

When opening a pull request, use the title of the issue as PR title.

A Pull Request must pass all tests to be merged.
