# MGC SDK

This repository holds the SDKs developed for Magalu Cloud (MGC). Each subdirectory
inside [libs/](./libs) translates to a Go library:

* **[OAPI Parser](./libs/parser/)**: reads an OpenAPI YAML or JSON specification and
create an intermediate structure that can be used by code generators or runtime code.
Useful for creating CLI Cobra actions, HTTP REST Client communicating with the MGC,
TF Resources, and others.

* **[CLI](./libs/cli)**: Go CLI, using Cobra, with most commands and actions created and
defined in runtime based on the parsed OpenAPI spec.

* **TF Provider**: Terraform provider plugin with its resources generated from the
parsed OpenAPI spec.

Most of our code is written in Golang, however there are some utility scripts written
in Python as well.

## Dependencies

To run the project, the main dependency needed is [Go](https://go.dev/dl/). To
install, visit the official link with the instructions.

There are some utility scripts written in [Python](https://www.python.org/downloads/).
To install, visit the official website.

## Contributing

### pre-commit

We use [pre-commit](https://pre-commit.com/) to install git hooks and enforce
lint, formatting, tests, commit messages and others. This tool depends on
Python as well. On pre-commit we enforce:

* On `commit-msg` for all commits:
    * [Conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) pattern
    with [commitzen](https://github.com/commitizen/cz-cli)
* On `pre-commit` for Go files:
    * Complete set of [golangci-lint](https://golangci-lint.run/): `errcheck`,
    `gosimple`, `govet`, `ineffasign`, `staticcheck`, `unused`
* On `pre-commit` for Python files:
    * `flake8` and `black` enforcing pep code styles

#### Installation

#### Mac
```sh
brew install pre-commit
```

#### pip

```sh
pip install pre-commit
```

For other types of installation, check their
[official doc](https://pre-commit.com/#install).

#### Configuration

After installing, the developer must configure the git hooks inside its clone:

```sh
pre-commit install
```

### Linters

We install the go linters via `pre-commit`, so it is automatically run by the
pre-commit git hook. However, if one wants to run standalone it can be done via:

```sh
pre-commit run golangci-lint
```
