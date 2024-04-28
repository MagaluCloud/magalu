#!/bin/bash

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
MGCDIR=${MGCDIR:-"mgc/cli/"}
DUMP_TREE="script-qa/cli-dump-tree.json"
OUT_DIR="script-qa/cli-doc"

set -xe
cd $MGCDIR
go build

echo "generating $OUT_DIR..."
python3 ../../scripts/gen_expected_cli_doc_output.py \
    ./mgc \
    "../../$DUMP_TREE" \
    "../../$OUT_DIR"
echo "generating $OUT_DIR: done"
