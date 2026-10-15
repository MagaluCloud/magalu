import argparse

import yaml
from spec_add_project_scope import add_project_scope


def load_yaml(path: str) -> dict:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.FullLoader)


def save_yaml(spec: dict, path: str):
    with open(path, "w") as fd:
        yaml.dump(spec, fd, sort_keys=False, indent=4, allow_unicode=True)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="Apply project scope",
        description="Expands the project scope marker into the header "
        "parameter that carries it, on every path of a scoped spec.",
    )
    parser.add_argument(
        "spec_file",
        type=str,
        nargs="+",
        help="Merged spec files to expand, in place",
    )
    args = parser.parse_args()

    for path in args.spec_file:
        spec = load_yaml(path)
        added = add_project_scope(spec)
        if not added:
            continue

        save_yaml(spec, path)
        print(f"{path}: project scope added to {added} paths")
