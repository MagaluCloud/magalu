# Development

## Install

Install dependencies using:

```sh
go install
```

## Build

Build the command line using:

```sh
go build -o mgc # basic build without embedded openapi
go build -tags "embed" -ldflags "-s -w" -o mgc # stripped build, with embedded openapi
```

Alternatively, during development one may run without building:

```sh
go run main.go
```

> **NOTE:**
> consider using `scripts/build_release.sh` to build with recommended
> flags for all supported platforms.

## How to add OpenAPI Descriptions (OAD)

Run [add_specs.sh](../../scripts/README.md#add_specssh):

```sh
scripts/add_specs.sh <API_NAME> <API_URL> <CANONICAL_URL>
```
It will start a process that uses the following scripts:

- [sync_oapi.py](../../scripts/README.md#sync_oapipy)
- [remove_tenant_id.py](../../scripts/README.md#remove_tenant_idpy)
- [yaml_merge.py](../../scripts/README.md#yaml_mergepy)

This process uses a file called `<API_NAME>.openapi.yaml` stored in `$CUSTOM_DIR` (default `openapi-customizations`) which is used to add interface-specific modifications that make CLI and TF usage cleaner.

In short, what it does is:

1. Get the internal and external (the public one) OAD;
1. See if there is a difference between the current external requestBody and the internal requestBody, if there is, update the external with the internal requestBody schema;
1. Replace the Error object by simplifying the error response with an object containing `message` and `slug` (The Magalu Kong gateway already does this for the external OAD but as the internal one does not go through Kong we have to make this replacement here too);
1. Remove `x-tenant-id` param from OAD actions;
1. Merge the OAD with the custom file mentioned above.

The result is a YAML file called `<API_NAME>.openapi.yaml` stored in `$OAPIDIR` (default `mgc/cli/openapis`) which is ready to be used by the CLI.

> **NOTE:**
>For the resulting yaml to actually be seen it needs to be indexed, this is easily done by running [oapi_index_gen.py](../../scripts/README.md#oapi_index_genpy):
> ```sh
> python3 ./scripts/oapi_index_gen.py "--embed=mgc/sdk/openapi/embed_loader.go" mgc/cli/openapis
>```

Alternatively, use [add_all_specs.sh](../../scripts/README.md#add_all_specssh) by editing it (add a line executing [add_specs.sh](../../scripts/README.md#add_specssh)) and runing it.
