# noqa: T201

import os
from typing import Dict, Any, List, TypedDict
import argparse
import yaml

OAPISchema = Dict[str, Any]


class IndexModule(TypedDict):
    name: str
    path: str
    version: str
    description: str


IndexModules = List[IndexModule]


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def save_index(mods: IndexModules, path: str):
    with open(os.path.join(path, "index.yaml"), "w") as fd:
        yaml.dump(mods, fd, indent=4, allow_unicode=True)


def load_mods(oapiDir: str, outDir: str = None):
    if outDir is None:
        outDir = oapiDir

    mods = []
    for filename in sorted(os.listdir(oapiDir)):
        name, ext = os.path.splitext(filename)
        if name == "index" or ext != ".yaml":
            print("ignored file:", filename)
            continue

        filepath = os.path.join(oapiDir, filename)
        relpath = os.path.relpath(filepath, outDir)

        data = load_yaml(filepath)
        info = data["info"]
        mods.append(
            IndexModule(
                name=name.split(".")[0],
                path=relpath,
                description=info.get("x-cli-description", info.get("description", "")),
                version=info.get("version", ""),
            )
        )
    return mods


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate index file for all OAPI YAML files in directory",
    )
    parser.add_argument(
        "dir",
        type=str,
        help="Directory of openapi files",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        help="Directory to save the new index YAML. Defaults to openapi directory",
    )
    args = parser.parse_args()

    mods = load_mods(args.dir, args.output)
    print("indexed modules:")
    for mod in mods:
        print(mod)

    save_index(mods, args.output or args.dir)
