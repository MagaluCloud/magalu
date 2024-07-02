#!/usr/bin/env python3

from typing import Iterator, TextIO

import argparse
import json
import subprocess
import sys


def gen_cli_dump_tree(cli: str) -> Iterator[list[str]]:
    args = [cli, "dump-tree", "-o", "json", "--raw"]
    with subprocess.Popen(args, stdout=subprocess.PIPE, encoding="utf-8") as p:
        tree = json.load(p.stdout)
        assert isinstance(tree, list)
        return tree


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate expected CLI dump-tree output",
    )
    parser.add_argument(
        "cli",
        type=str,
        help="the binary to use during executions",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=argparse.FileType("w"),
        default=sys.stdout,
    )
    args = parser.parse_args()

    out_file: TextIO = args.output
    tree = gen_cli_dump_tree(args.cli)
    json.dump(
        tree,
        out_file,
        ensure_ascii=False,
        indent=True,
        sort_keys=True,
    )
