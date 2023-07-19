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

## Dependencies

To run the project, the only dependency needed is [Go](https://go.dev/dl/). To
install, visit the official link with the instructions.
