#!/bin/bash

me=$(basename "$0")
echo "RUNNING $me"

MGCDIR=${MGCDIR:-"mgc/cli/"}
DUMP_TREE="mgc/cli/cli-dump-tree.json"
OUT_DIR="mgc/cli/docs"

set -xe
cd $MGCDIR

go build -tags \"embed\" -o mgc

echo "generating $OUT_DIR..."
python3 ../../scripts/gen_expected_cli_doc_output.py \
    ./mgc \
    "../../$DUMP_TREE" \
    "../../$OUT_DIR"
echo "generating $OUT_DIR: done"

echo "ENDING $me"
