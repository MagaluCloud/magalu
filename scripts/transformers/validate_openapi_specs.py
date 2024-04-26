#!/usr/bin/env python3

from typing import Any, Dict
import yaml
import os.path
import jsonschema

OAPISchema = Dict[str, Any]

base_dir = os.path.dirname(os.path.dirname(os.path.dirname(__file__)))

# NOTE: this file was downloaded from:
# https://www.schemastore.org/json/ (registry) as:
# https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/schemas/v3.0/schema.json
OAPI_JSON_SCHEMA_FILE = os.path.join(base_dir, "jsonschemas", "openapis.json")

with open(OAPI_JSON_SCHEMA_FILE) as fd:
    schema = yaml.load(fd, Loader=yaml.FullLoader)

# we extend the schema with non-standard "$id",
# but without an actual "x-" prefix. Declare we know that variable:
schema["properties"]["$id"] = {"type": "string"}

# format_validator must be explicitly given, so get the validator first:
validator = jsonschema.validators.validator_for(schema)


def validate_oapi(spec: OAPISchema) -> None:
    # Descomentar para encontrar o erro...
    # print("SCHEMA:\n\n")
    # print(schema)
    # print("SPEC:\n\n")
    # print(spec)
    jsonschema.validate(
        spec,
        schema,
        cls=validator,
        format_checker=validator.FORMAT_CHECKER,
    )


if __name__ == "__main__":
    import argparse
    import re

    parser = argparse.ArgumentParser(
        description="Validate OpenAPI Specification (YAML or JSON) Files",
    )
    parser.add_argument(
        "files",
        nargs="+",
        type=argparse.FileType("r"),
        help="file to be validated",
    )
    args = parser.parse_args()

    name_pattern = re.compile("(?!index).*[.]openapi[.](yaml|yml|json)$")

    for fd in args.files:
        name = os.path.basename(fd.name)
        # ignore garbage that may be received by shell's "*""
        if not name_pattern.match(name):
            print(f"{fd.name}: ignored")
            continue

        spec = yaml.load(fd, Loader=yaml.FullLoader)
        try:
            validate_oapi(spec)
            print(f"{fd.name}: ok")
        except Exception as e:
            raise SystemExit(f"failed {fd.name!r}: {e}") from e
