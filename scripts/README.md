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
