#!/bin/bash
me=$(basename "$0")
echo "RUNNING $me"

MGCDIR=${MGCDIR:-"mgc/cli/"}
OUT_FILE="mgc/cli/cli-dump-tree.json"

set -xe
cd $MGCDIR

go build -tags \"embed\" -o mgc

echo "generating $OUT_FILE..."
python3 ../../scripts/gen_expected_cli_dump_tree.py ./mgc -o "../../$OUT_FILE"
echo "generating $OUT_FILE: done"

echo "ENDING $me"
