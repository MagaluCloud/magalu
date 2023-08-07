# OpenAPI

The MGC SDK is heavily based on the
[OpenAPI](https://www.openapis.org/),
most of the commands are auto-generated in runtime from a
[schema file](https://spec.openapis.org/oas/latest.html)
where each schema file is mapped to a module, each module is composed of
resources (OpenAPI tag), each composed by actions (OpenAPI operation).

```
index.yaml   -> module-name.openapi.yaml -> tag        -> operation
[entrypoint]    [module: module-name]       [resource]    [action]
```

## Reading

The SDK may contain embedded OpenAPI files if built with `-tags "embed"`.
In addition to that, it will look into the directory defined by the
environment variable `$MGC_SDK_OPENAPI_DIR` or `./openapis` it not set.

> **NOTE:**
> if using a binary with embedded files, one may still provide overrides
> by using a file `./openapis/file-to-be-overridden.openapi.yaml`.
> In order to add a new file, one must create the `index.yaml`
> including that file.


## Entry Point (index.yaml)

The entrypoint listing all the desired schema files using the format below.
This file is a JSON or YAML file with the name `index.yaml`
(regardless if it's JSON or YAML, always use this **exact** filename).

```yaml
version: 1.0.0 # index version, must be "1.0.0"
modules:
-   description: Your Module Description
    name: module-name
    path: module-name.openapi.yaml
    version: 1.2.3 # your module version
```

> **NOTE:**
> it's easier to generate the index using `scripts/oapi_index_gen.py`


## Extensions

Some extensions may be added in the OpenAPI in order to control the
runtime generation:

* `x-cli-name` may be present to change the tag, operation, parameter
  or JSON Schema properties (ie: request body) to control its name.
* `x-cli-name` like `x-cli-name`, but affects the description.
* `x-cli-hidden: true` like `x-cli-name`, but if `true` will skip
  using such entry.

> **NOTE:**
> it's easier to keep the customizations in another YAML file and use
> `scripts/yaml_merge.py` to merge the original file with the
> desired customizations, producing the final file to be used.

### Example

The following snippet show how to customize `visible-tag`, giving it
another name `my-resource-name` and description `my resource description`.

We're hiding `POST /v0/some/path` using `x-cli-hidden`, so
`my-resource-name` will have a single action `GET /v0/some/path` that
will be named `retrieve` and description `this will be used`.

Note that `hidden-tag` is not used since it's hidden.

```yaml
paths:
   /v0/some/path:
        get:
            tags:
            - visible-tag
            description: this won't be used
            x-cli-name: retrieve
            x-cli-description: this will be used
        post:
            tags:
            - visible-tag
            x-cli-hidden: true

   /v0/other/path:
        get:
            tags:
            - hidden-tag
tags:
-   name: visible-tag
    description: this won't be used
    x-cli-name: my-resource-name
    x-cli-description: my resource description
-   name: hidden-tag
    x-cli-hidden: true # no resources using this tag will be visible
```

## Parameters x Config

Server variables, header and cookie parameters are handled as **Config**.

Query and path parameters, as well as request body properties are
handled as **Parameters**.
