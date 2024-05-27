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

### [transform.py](./transformers/transform.py)

Loads a product's openapi spec file or fetch from an URL and transforms or
removes internal structures (before the Gateway) to prepare it for
"external/public" usage.

The CLI and TF solutions are based on public openapi specifications, this
scripts allows internal specs to be used by them.

Usage:

```shell
python3 ./transformers/transform.py $SPEC_FILE $SPEC_UID -o $SPEC_OUTPUT_FILE
```

```shell
python3 ./transformers/transform.py $SPEC_FILE_URL $SPEC_UID -o $SPEC_OUTPUT_FILE
```

### [yaml_merge.py](./yaml_merge.py)

Merges spec customizations created for the CLI/TF interfaces into an already
existing spec, the existing spec can be overriden or a new spec can be generated
from the output.

> Make sure a raw product spec goes through the transformations before merging
the customizations.

Usage:

```shell
python3 ./scripts/yaml_merge.py [-h] [--override] [-o OUTPUT] <BASE> <CUSTOMIZATIONS>
```

Example:

```shell
python3 ./scripts/yaml_merge.py --override ./mgc/cli/openapis/block-storage.openapi.yaml ./openapi-customizations/block-storage.openapi.yaml
```

### [oapi_index_gen.py](./oapi_index_gen.py)

Generates an index file for the MGC_SDK consumption, the SDK used by the
interfaces will only load specs defined in the `index.openapi.yaml` file, specs will
reference each other based on their SPEC_UID (see `transforms.py` for more info).

The `--embed` options will insert the specs into the Go binary, allowing the
user to use the SDK without having `.openapi.yaml` files in a specific folder.

Usage:

```shell
python3 ./scripts/oapi_index_gen.py [-h] [-o OUTPUT] [--embed EMBED] dir
```

Example:

```shell
python3 ./scripts/oapi_index_gen.py "--embed=mgc/sdk/openapi/embed_loader.go" mgc/cli/openapis
```

### [add_specs.sh](./add_specs.sh)

Shell script to add new OpenAPI specifications into the command line. It will
receive a yaml file, parse, transform, customize and insert into the index file
for the SDK usage.

If no customization file exists for the current spec, a new one will be created.

Usage:

```shell
./scripts/add_specs.sh $API_NAME $API_ENDPOINT_NAME $API_SPEC_FILE $SPEC_UID
```

Example:

```shell
./scripts/add_specs.sh block-storage volume ./block-storage.openapi.json https://block-storage.jaxyendy.com/openapi.json
```

### [sync_oapi.py](./sync_oapi.py)

Sync external OAPI schema with the internal schema by fixing any mismatch of
requestBody between external and internal implementation.

Usage:

```shell
python3 ./scripts/sync_oapi.py  [-h] [--ext EXT] [-o OUTPUT] <INTERNAL_SPEC_URL> <SPEC_UID>
```
Example:

```shell
python3 ./scripts/sync_oapi.py https://block-storage.br-ne-1.jaxyendy.com/openapi.json https://block-storage.jaxyendy.com/openapi.json --ext ./mgc/cli/openapis/block-storage.openapi.yaml
```

### [spec_stats.py](./spec_stats.py)

Shows general statistics, information that could generate problems with CLI or TF interfaces, or wrong REST definitions. It is good for validating that there are no unwanted interfaces from the endpoints. Missing crud, for example, is a useful statistic for us, it doesn't necessarily mean the API is broken.

Usage:

```shell
python3 ./scripts/spec_stats.py [-h] [--filter FILTER] [--filter-out FILTER_OUT] [-o OUTPUT] [--ignore-disabled IGNORE_DISABLED] dir_or_file
```

Example:

```shell
python3 ./scripts/spec_stats.py ./mgc/cli/openapis
```

### [tf_generate_docs.sh](./tf_generate_docs.sh)

Uses [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs#terraform-plugin-docs) to generate documentation about Terraform providers and resources.

Usage:

```shell
python3 ./scripts/tf_generate_docs.sh
```

### [gen_expected_cli_help_output.py](./gen_expected_cli_help_output.py)

This script will loop through all the CLI commands and generate output corresponding to each of them. When run, it will iterate through the commands available in the CLI, producing the output associated with each command in sequence.

```shell
python3 ./scripts/gen_expected_cli_help_output.py ./mgc
```

### [gen_expected_output_file.sh](./gen_expected_output_file.sh)

This script has the purpose of creating a file based on the output of the `gen_expected_cli_help_output.py`. When executed, it uses the information produced by the script as input. This approach aims to automate and simplify the file creation process and help on test CLI output.

Usage:

```shell
./scripts/gen_expected_output_file.sh
```
