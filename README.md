# MGC

This repository contains the monorepo developed for Magalu Cloud (MGC). Each subdirectory within [mgc/](./mgc) corresponds to a Go module:

- **[Core](./mgc/core)**: Definition of data types used for the intermediate
  structure generated by the SDK after parsing either an OpenAPI spec or a static
  definition.

- **[SDK](./mgc/sdk/)**: Implement the concrete data types defined by core
  to generate usable commands. The commands can be used by the CLI via Cobra commands, or
  by Terraform Plugin Provider to perform CRUD on resources.

- **[CLI](./mgc/cli)**: Go CLI, using Cobra, with commands and actions defined by
  the SDK. The commands can either come from dynamic loaded OpenAPI spec or static
  modules, i.e: authentication.

Our code is written in Golang, however there are some utility scripts written
in Python as well.

**Looking for Terraform code? Check out the [Terraform provider](https://github.com/MagaluCloud/terraform-provider-mgc) repository**

## Dependencies

To run the project, the main dependency needed is [Go](https://go.dev/dl/). To
install, visit the official link with the instructions.

There are some utility scripts written in [Python](https://www.python.org/downloads/).
For this, [Poetry](https://python-poetry.org/) is used. Check [Poetry.md](Poetry.md) for instructions.

## Building and running locally

Building needs [goreleaser](https://goreleaser.com/install/) and can be done using Makefile targets.

If you have API spec changes, update them on `specs`, on the corresponding product and run

```bash
$ make refresh-specs
```

After that (or if there are no further spec changes), build the CLI with:

```bash
$ make build-local
```

If all goes well, the output binary will be a platform-dependent directory, where it can be run:

```bash
$ cd dist/mgc_<your_platform>
$ ./mgc
```

## OpenAPI

See [sdk/openapi/README.md](./mgc/sdk/openapi/README.md)

## Adding new APIs

### OpenAPIs

To add a new API spec, see `ADD_NEW_API.md`.

### Static APIs

Manually written APIs should be added to `mgc/sdk/static`, follow the
structure in the exiting modules (`auth`, `config`).

## Contributing

### pre-commit

We use [pre-commit](https://pre-commit.com/) to install git hooks and enforce
lint, formatting, tests, commit messages and others. This tool depends on
Python as well. On pre-commit we enforce:

- On `commit-msg` for all commits:
  - [Conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) pattern
    with [commitzen](https://github.com/commitizen/cz-cli)
- On `pre-commit` for Go files:
  - Complete set of [golangci-lint](https://golangci-lint.run/): `errcheck`,
    `gosimple`, `govet`, `ineffasign`, `staticcheck`, `unused`
- On `pre-commit` for Python files:
  - `flake8` and `black` enforcing pep code styles

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

### Run all

Run pre-commit without any file modified:

```sh
pre-commit run -a
```
