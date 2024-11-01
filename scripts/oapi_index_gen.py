# noqa: T201

import os
from typing import Dict, Any, List, TypedDict, Tuple
import argparse
import yaml
import re

OAPISchema = Dict[str, Any]


class IndexModule(TypedDict):
    name: str
    url: str
    path: str
    version: str
    description: str


IndexModules = List[IndexModule]


class IndexFile(TypedDict):
    version: str
    modules: IndexModules


modname_re = re.compile("^(?P<name>[a-z0-9-]+)[.]openapi[.]yaml$")
index_filename = "index.openapi.yaml"
index_version = "1.0.0"


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.FullLoader)


def save_index(mods: IndexModules, path: str):
    with open(os.path.join(path, index_filename), "w") as fd:
        idx_file = IndexFile(version=index_version, modules=mods)
        yaml.dump(idx_file, fd, indent=4, allow_unicode=True)
        return idx_file


def load_mods(
    oapiDir: str, outDir: str | None = None
) -> Tuple[Dict[str, OAPISchema], IndexModules]:
    if outDir is None:
        outDir = oapiDir

    full_mods = {}
    mods = []
    for filename in sorted(os.listdir(oapiDir)):
        if filename == index_filename:
            continue
        match = modname_re.match(filename)
        if not match:
            if filename != index_filename:
                print("ignored file:", filename)
            continue
        filepath = os.path.join(oapiDir, filename)
        relpath = os.path.relpath(filepath, outDir)
        data = load_yaml(filepath)
        info = data["info"]
        url = data["$id"]
        name = match.group("name")
        full_mods[filename] = data
        description = info.get("x-mgc-description", info.get("description", ""))
        mods.append(
            IndexModule(
                name=name,
                url=url,
                path=relpath,
                description=description,
                version=info.get("version", ""),
                summary=info.get("summary", description),
            )
        )
    return full_mods, mods


embed_json_opts = {
    "separators": (",", ":"),
    "ensure_ascii": False,
    "sort_keys": True,
}


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

    full_mods, mods = load_mods(args.dir, args.output)

    idx_file = save_index(mods, args.output or args.dir)
