MGC Terraform
=============

## Provider

### Install

This will create an override to make terraform look for the provider in the
local environment(this folder) and not in the remote registry.

```
./install.sh
```

To utilize the provider for now, the user must use the CLI for authentication
with the following command:

```sh
./mgc auth login
```

It will request your MagaLu Cloud ID, and save your credentials for the provider
usage.

### SDK

The provider makes use of the MGC SDK to identify and create the terraform
resources. The SDK has an internal representation of operations that a user can
execute in the MGC APIs, and the provider clusters related operations together
to create a resource.

To make use of resources auto-generated from the SDK through OpenAPI and Blueprint
specs, we must tell the SDK in which folder those specs are defined.

#### OpenAPIs And Blueprints

To use the OAPI and Blueprint files with the terraform examples it's necessary to set
where the files are. We recommend using absolute paths, but if you want to use
relative paths, make sure they are defined in relationship to the `*.tf` file
you are using:

```sh
export MGC_SDK_OPENAPI_DIR=~/<repo_root_folder>/mgc/cli/openapis
export MGC_SDK_BLUEPRINTS_DIR=~/<repo_root_folder>/mgc/cli/blueprints
```

>For now, the specs are defined in a specific folder, we probably will move it
to allow users to "install" specific products into the SDK.

### Build

Builds the provider to be used by the terraform application:

```sh
go build
```

### Examples

For development we have a `./run.sh` script that builds and executes a
terraform/opentf command with a specific `.tf` example file.

To check if the provider is being correctly installed and generating resources
you can run the `provider` example (present in the `examples/provider` folder),
by executing the following command:

```
./run.sh provider plan
OR
./run.sh provider apply
```

> Running this command requires authentication to the magalu cloud. See install
section.

Other examples can be executed by changing the example name in the command:

```sh
./run.sh <example> <terraform cmd>
```

### Generate docs

To automatically generate the documentation for the Provider and it's Resources, run the command below. The documentation will be found at `terraform-provider-mgc/docs`.

```sh
(cd ../.. && ./scripts/tf_generate_docs.sh)
```

The docs will be generated at `./docs` folder.

> This script is integrated into the `build_release.sh` script. This ensures
that all releases have up to date TF documentation

## OpenTofu
To use OpenTofu for install scripts, set the environment variable `MGC_OPENTF` to one

```
export MGC_OPENTF=1
```
