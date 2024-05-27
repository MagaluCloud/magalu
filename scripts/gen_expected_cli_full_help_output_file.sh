#!/bin/bash
me=$(basename "$0")
echo "RUNNING $me"


BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
MGCDIR=${MGCDIR:-"mgc/cli/"}
DUMP_TREE="script-qa/cli-dump-tree.json"
OUT_DIR="script-qa/cli-help"

set -xe
cd $MGCDIR

go build -tags \"embed\" -o mgc

echo "generating $OUT_DIR..."
python3 ../../scripts/gen_expected_cli_help_output.py \
    ./mgc \
    "../../$DUMP_TREE" \
    "../../$OUT_DIR"
echo "generating $OUT_DIR: done"

echo "ENDING $me"
