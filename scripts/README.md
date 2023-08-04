# Utility Scripts

The scripts in this folder are utilities for the MGC SDKs. Most of them are written in
Python following PEP convention enforced by `flake8` and `black`.

One can run most of them with:

```shell
python3 <script> -h
```

## Scripts

### [add_all_specs.sh](./add_all_specs.sh)

Run [add_specs.sh](./add_specs.sh) with all supported specifications.

### [add_specs.sh](./add_specs.sh)

Shell script to add OpenAPI specifications from remote. It will fetch, parse, create
customizations and leave ready for usage of CLI. Example:

```shell
./scripts/add_specs.sh mke https://mke.br-ne-1.jaxyendy.com/docs/openapi-with-snippets.json
```

This will also create a new customization file if not present with the following content:

```yaml
# This file is to be merged on top of $OAPI_PATH/$SPEC_FILE
# using yaml_merge.py
# NOTE: Lists are merged by their indexes, be careful with parameters, tags and such!
# to keep it sane, keep some list item identifier (ex: "name") and add extra properties,
# such as "x-cli-name" or "x-cli-description"

servers:
-   url: https://api-$API_NAME.{region}.jaxyendy.com
    variables:
        region:
            description: Region to reach the service
            default: br-ne-1
            enum:
            - br-ne-1
            - br-ne-2
            - br-se-1
```
