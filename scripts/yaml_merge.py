from typing import Any, Dict, Hashable
import argparse
import yaml

OAPISchema = Dict[str, Any]


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def save_external(spec: OAPISchema, path: str):
    with open(path, "w") as fd:
        yaml.dump(spec, fd, sort_keys=False, indent=4, allow_unicode=True)


def merge_item(value: Any, new: Any, override: bool, path: list[str]) -> bool:
    if value == new:
        return False

    value_type = type(value)
    new_type = type(new)
    hasheable_types = (
        isinstance(value_type, Hashable),
        isinstance(new_type, Hashable),
    )
    if not all(hasheable_types) and value_type != new_type:
        raise NotImplementedError(
            f"{path}: cannot merge type {value_type!r} with {new_type!r}"
        )

    if value_type == dict:
        merge_dict(value, new, override, path)
    elif value_type == list:
        merge_list(value, new, override, path)
    elif not override:
        raise ValueError(f"{path}: not overriding {value!r} with {new!r}")
    else:
        return True


def merge_dict(dst: dict, extra: dict, override: bool, path: list[str]) -> dict:
    for k, new in extra.items():
        value = dst.setdefault(k, new)
        if value is not new:
            path.append(k)
            replace = merge_item(value, new, override, path)
            path.pop()
            if replace:
                dst[k] = new

    return dst


def merge_list(dst: list, extra: list, override: bool, path: list[str]) -> list:
    for i, new in enumerate(extra):
        if len(dst) <= i:
            dst.append(new)
        else:
            value = dst[i]
            path.append(i)
            replace = merge_item(value, new, override, path)
            path.pop()
            if replace:
                dst[i] = new

    return dst


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Merge two YAML files",
    )
    parser.add_argument(
        "base",
        type=str,
        help="the base file to open",
    )
    # External = Viveiro in MGC context, intermediate between product and Kong
    parser.add_argument(
        "extra",
        type=str,
        help="the extra file to merge on top of base",
    )
    parser.add_argument(
        "--override",
        action="store_true",
        default=False,
        help="Override existing scalars",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        help="Path to save the new external YAML. Defaults to overwrite base",
    )
    args = parser.parse_args()

    base = load_yaml(args.base)
    extra = load_yaml(args.extra)

    merge_dict(base, extra, args.override, [])

    save_external(base, args.output or args.base)
