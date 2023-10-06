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

Usage:

```shell
./scripts/add_specs.sh <API_NAME> <API_URL> <CANONICAL_URL>
```

Example:

```shell
./scripts/add_specs.sh block-storage https://block-storage.br-ne-1.jaxyendy.com/openapi.json https://block-storage.jaxyendy.com/openapi.json
```

Shell script to add OpenAPI specifications from remote. It will fetch, parse, create
customizations and leave ready for usage of CLI.

It also creates a new customization file if it doesn't exist.


### [sync_oapi.py](./sync_oapi.py)

Usage:

```shell
python3 ./scripts/sync_oapi.py  [-h] [--ext EXT] [-o OUTPUT] <INTERNAL_SPEC_URL> <CANONICAL_URL>
```
Example:

```shell
python3 ./scripts/sync_oapi.py https://block-storage.br-ne-1.jaxyendy.com/openapi.json https://block-storage.jaxyendy.com/openapi.json --ext ./mgc/cli/openapis/block-storage.openapi.yaml
```

Sync external OAPI schema with the internal schema by fixing any mismatch of requestBody between external and internal implementation. After that, we change the server URL to Kong and adjust schema of error returns.
