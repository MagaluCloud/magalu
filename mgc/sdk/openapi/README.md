# OpenAPI

The MGC SDK is heavily based on the
[OpenAPI](https://www.openapis.org/),
most of the commands are auto-generated in runtime from a
[schema file](https://spec.openapis.org/oas/latest.html)
where each schema file is mapped to a module, each module is composed of
resources (OpenAPI tag), each composed by actions (OpenAPI operation).

```
index.openapi.yaml   -> module-name.openapi.yaml -> tag        -> operation
[entrypoint]         [module: module-name]       [resource]    [action]
```

Currently, the spec of container-registry is the best "example" to follow: https://mcr.br-ne-1.jaxyendy.com/docs/openapi.yaml

## Reading

The SDK may contain embedded OpenAPI files if built with `-tags "embed"`.
In addition to that, it will look into the directory defined by the
environment variable `$MGC_SDK_OPENAPI_DIR` or `./openapis` if not set.

> **NOTE:**
> if using a binary with embedded files, one may still provide overrides
> by using a file `./openapis/file-to-be-overridden.openapi.yaml`.
> In order to add a new file, one must create the `index.openapi.yaml`
> including that file.


## Adding new spec

In the `scripts/add_all_specs.sh` you can add a new spec, like this:
```
$BASEDIR/add_specs_without_region.sh profile profile mgc/spec_manipulator/cli_specs/conv.globaldb.openapi.yaml https://globaldb.jaxyendy.com/openapi-cli.json
echo "SSH"

# EXAMPLE
# $BASEDIR/SCRIPT.sh NOME_NO_MENU URL_PATH LOCAL_DA_SPEC HTTPS://LOCAL_DA_SPEC
```

After this, just run `./scripts/add_all_specs.sh` and BUILD all.

## Entry Point (index.openapi.yaml)

The entrypoint listing all the desired schema files using the format below.
This file is a JSON or YAML file with the name `index.openapi.yaml`
(regardless if it's JSON or YAML, always use this **exact** filename).

```yaml
version: 1.0.0 # index version, must be "1.0.0"
modules:
-   description: Your Module Description
    name: module-name
    path: module-name.openapi.yaml
    version: 1.2.3 # your module version
```


## Extensions

Some extensions may be added in the OpenAPI spec in order to control the runtime generation. They are prefixed with `x-mgc`
and can be used to edit names and more complex behaviors. Be aware that some extensions can only be used in certain OpenAPI spec
elements. The following list shows which extensions can be used in the spec:

- Parameter
    - `x-mgc-name`
    - `x-mgc-description`
    - `x-mgc-hidden`
- Tag
    - `x-mgc-name`
    - `x-mgc-description`
    - `x-mgc-hidden`
- Operation
    - `x-mgc-name`
    - `x-mgc-description`
    - `x-mgc-hidden`
    - `x-mgc-confirmable`
    - `x-mgc-confirmPrompt`
    - `x-mgc-wait-termination`
    - `x-mgc-output-flag`
- Link
    - `x-mgc-wait-termination`
    - `x-mgc-extra-parameters`
    - `x-mgc-hidden`
- Schema
    - `x-mgc-name`
    - `x-mgc-description`
    - `x-mgc-hidden`

### `x-mgc-name`

Use this extension to rename a tag, operation, parameter or a JSON Schema property (i.e: request body).

```yaml
tags:
  - name: tag_key
    x-mgc-name: new_tag_name
```

### `x-mgc-description`

Use this extension to edit a description in the OpenAPI spec.

```yaml
paths:
   /v0/some/path:
        get:
            description: operation description
            x-mgc-description: edited description. This one will be used
```

### `x-mgc-hidden`

Use this extension to hide a tag or path in the OpenAPI spec. Anything marked with this tag will be invisible to the
autocomplete and help outputs, unless `--cli.show-internal` is passed. The operations can still be accessed by passing
their explicit names.


```yaml
paths:
   /v0/some/path:
        get:
            x-mgc-hidden: true
```

### `x-mgc-confirmable`

Add this extension to an operation to require user confirmation before execution in the CLI (all `delete` operations apply
this extension by default). This extension is an object with a single string property called `message` used to define the
confirmation message to be shown to the user.


```yaml
paths:
   /v0/some/path:
        patch:
            x-mgc-confirmable:
                message: "This action requires confirmation. Are you sure you wish to continue?"
```

### `x-mgc-promptInput`

Add this extension to an operation in the CLI to demand stronger user confirmation before execution.
It consists of an object with two string properties: `message` and `confirmValue`.
The `message` can incorporate `{{.confirmationValue}}`, replaced with the content of `confirmValue`.
Furthermore, `message` has access to `{{.parameters}}` and `{{.configs}}`.
Essentially, `message` serves as the confirmation prompt shown to the user, while `confirmValue`
specifies the input required from the user to confirm the operation.


```yaml
paths:
   /v0/some/path:
        patch:
            x-mgc-promptInput:
                message: "This action requires stronger confirmation. Please retype `{{.confirmationValue}} to confirm"
                confirmValue: "I agree with this operation"
```

### `x-mgc-wait-termination`

Add this extension to an operation to add a termination conditon. The operation will be executed until the condition
is satisfied, or until a maximum number of attempts. The `x-mgc-wait-termination` extension is an object with three properties:

- `maxRetries`: an integer defining the max number of attempts
- `interval`: interval in seconds between each attempt
- `jsonPathQuery`: the termination condition expressed in jsonpath syntax
- `templateQuery`: the termination condition expressed in Go Template syntax

```yaml
paths:
   /v0/some/path:
        post:
            x-mgc-wait-termination:
                maxRetries: 10
                interval: 1s
                jsonPathQuery: $.result.status == "completed"
```

### `x-mgc-output-flag`

Defines the default output format. Accepted formats: json, yaml, table, template, jsonpath, template-file and jsonpath-file.
Similar to the `--cli.output`/`-o` flag in the CLI.

When choosing `table` you can specify columns and rows using jsonpath syntax. The following example outputs a table with
three columns (ID, NAME and VERSION) where each row is defined by a jsonpath expression.

```yaml
paths:
   /v0/some/path:
        post:
            x-mgc-output-flag: yaml # Will show output as yaml
            x-mgc-output-flag: table=ID:$.images[*].id,NAME:$.images[*].name,VERSION:$.images[*].version
            x-mgc-output-flag: jsonpath=$.id
            x-mgc-output-flag: template={{.id}}
            x-mgc-output-flag: remove=$.machine_types[*].sku|$.machine_types[*].status
```


### `x-mgc-extra-parameters`

Add extra parameters to a link. This extension is an array of objects of the form:

- name: parameter name
- required: a bool indicating whether the parameter is required or not
- schema: parameter schema

> NOTE: if a extra parameter has the same name of an existing parameter in the target request, it will not be added

```yaml
paths:
   /v0/some/path:
        post:
            links:
                delete:
                    x-mgc-extra-parameters:
                        - name: id
                          required: true
                          schema:
                            type: string
                            format: uuid
                            title: Id
```

## Parameters x Config

Server variables, header and cookie parameters are handled as **Config**.

Query and path parameters, as well as request body properties are
handled as **Parameters**.

## Links

Links are an effective strategy for chaining operations. This is because links
make it easier to use a command response to call another executor,
automatically mapping the parameters.

When you add a new OpenAPI YAML file using the `scripts/add_specs.sh` script,
the links between endpoints are automatically generated and included in the
specification.

Example: Suppose you have two endpoints in your API:

```yaml
/product (POST) : Creates a new product
/product/{product_id} (GET): Retrieves details about a specific product
```

The script will generate a link from `/product` to `/product/{product_id}`.
This allows users to easily create a new product using the POST method and then
immediately retrieve the details of the newly created product using the GET
method with the ID in POST response
