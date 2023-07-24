from typing import Any, Dict
import argparse
import yaml


OAPISchema = Dict[str, Any]


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def save_external(spec: OAPISchema, path: str):
    with open(path, "w") as fd:
        yaml.dump(spec, fd, sort_keys=False, indent=4, allow_unicode=True)


def remove_param(spec: OAPISchema, param_name: str):
    for path in spec.get("paths", {}).values():
        for action in path.values():
            if "parameters" not in action:
                continue

            filtered_params = [
                p for p in action.get("parameters", [{}]) if p.get("name") != param_name
            ]

            if not filtered_params:
                del action["parameters"]
            else:
                action["parameters"] = filtered_params


def remove_tenant_id(spec: OAPISchema):
    return remove_param(spec, param_name="x-tenant-id")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Remove `x-tenant-id` param from OpenAPI spec actions"
    )
    # External = Viveiro in MGC context, intermediate between product and Kong
    parser.add_argument(
        "path",
        type=str,
        help="File path to an OpenAPI spec to be parsed",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        help="Path to save the modified YAML. Defaults to overwrite",
    )
    args = parser.parse_args()

    spec = load_yaml(args.path)

    remove_tenant_id(spec)

    save_external(spec, args.output or args.path)
