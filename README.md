# MGC SDK

This repository holds the SDKs developed for Magalu Cloud (MGC). Each subdirectory
inside [mgc/](./mgc) translates to a Go library:

* **[core](./mgc/core)**: Definition of data types used for the intermediate
structure generated by the SDK after parsing either an OpenAPI spec or a static
definition.

* **[sdk](./mgc/sdk/)**: Implement visitor pattern on the data types defined by core
to generate usable commands. The commands can be used by the CLI via Cobra commands, or
by Terraform Plugin Provider to perform CRUD on resources.

* **[CLI](./mgc/cli)**: Go CLI, using Cobra, with commands and actions defined by
the SDK. The commands can either come from dynamic loaded OpenAPI spec or static
modules, i.e: authentication.

Most of our code is written in Golang, however there are some utility scripts written
in Python as well.

## Dependencies

To run the project, the main dependency needed is [Go](https://go.dev/dl/). To
install, visit the official link with the instructions.

There are some utility scripts written in [Python](https://www.python.org/downloads/).
To install, visit the official website.


## Running the CLI

The quickest way to run in development is by going into the folder and running as:

```shell
cd mgc/cli
go run main.go
```

Or, build and run:

```shell
cd mgc/cli && go build -o cli && ./cli
```

### Authentication

We provide a Postman collection that handles authentication as well:
1. Import the Collection
2. Edit the Collection
3. Go to the `Authorization` Tab
4. Click on `Get New Access Token`
5. Copy the access token value and set as a env var:

```shell
export MGC_SDK_ACCESS_TOKEN=long-token
```

To ensure it is working, perform a CLI command that requires authentication:

```shell
cd mgc/cli
go run main.go virtual-machine instances get
```

## Adding new APIs

To add a new API spec, add the URL and name to `./scripts/add_specs.sh`. Then, run:

```shell
./scripts/add_specs.sh mgc/cli/openapis
```

This will fetch the URL, apply modifications and save in the given path `mgc/cli/openapis`.

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
