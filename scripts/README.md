# Utility Scripts

The scripts in this folder are utilities for the MGC SDKs. Most of them are written in
Python following PEP convention enforced by `flake8` and `black`.

One can run most of them with:

```shell
python3 <script> -h
```

## Scripts

### [build_release.sh](./build_release.sh)

Usage:

```shell
VERSION=v0.0.0 ./scripts/build_release.sh
```
> **NOTE:**
>`$VERSION` is used to set the correct version on build, the next version can be found by looking at git tags.

Creates the `build` directory with:

- Binaries for all supported platforms;
- Markdown documentations;
- `examples` directory with scripts that serve as examples;
- `openapis` directory with the required OpenAPI Descriptions;
- `docs` directory with TF documentation.

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

### [remove_tenant_id.py](./remove_tenant_id.py)

Usage:

```shell
python3 ./scripts/remove_tenant_id.py [-h] [-o OUTPUT] <PATH>
```

Example:

```shell
python3 ./scripts/remove_tenant_id.py ./mgc/cli/openapis/block-storage.openapi.yaml
```

Remove `x-tenant-id` param from OpenAPI spec actions.

### [yaml_merge.py](./yaml_merge.py)

Usage:

```shell
python3 ./scripts/yaml_merge.py [-h] [--override] [-o OUTPUT] <BASE> <EXTRA>
```

Example:

```shell
python3 ./scripts/yaml_merge.py --override ./mgc/cli/openapis/block-storage.openapi.yaml ./openapi-customizations/block-storage.openapi.yaml
```

Merge `EXTRA` YAML file on top of `BASE` YAML file.