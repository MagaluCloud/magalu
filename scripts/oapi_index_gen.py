# noqa: T201

import os
from typing import Dict, Any, List, TypedDict
import argparse
import yaml
import re

OAPISchema = Dict[str, Any]


class IndexModule(TypedDict):
    name: str
    path: str
    version: str
    description: str


IndexModules = List[IndexModule]


class IndexFile(TypedDict):
    version: str
    modules: IndexModules


modname_re = re.compile("^(?P<name>[a-z0-9-]+)[.]openapi[.]yaml$")
index_filename = "index.yaml"
index_version = "1.0.0"


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def save_index(mods: IndexModules, path: str):
    with open(os.path.join(path, index_filename), "w") as fd:
        idx_file = IndexFile(version=index_version, modules=mods)
        yaml.dump(idx_file, fd, indent=4, allow_unicode=True)


def load_mods(oapiDir: str, outDir: str = None):
    if outDir is None:
        outDir = oapiDir

    mods = []
    for filename in sorted(os.listdir(oapiDir)):
        match = modname_re.match(filename)
        if not match:
            if filename != index_filename:
                print("ignored file:", filename)
            continue

        filepath = os.path.join(oapiDir, filename)
        relpath = os.path.relpath(filepath, outDir)

        data = load_yaml(filepath)
        info = data["info"]
        mods.append(
            IndexModule(
                name=match.group("name"),
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
