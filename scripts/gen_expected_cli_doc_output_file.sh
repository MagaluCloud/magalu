#!/bin/bash

me=$(basename "$0")
echo "RUNNING $me"

MGCDIR=${MGCDIR:-"mgc/cli/"}
SPECS_DIR=${SPECS_DIR:-"mgc/spec_manipulator/"}
DUMP_TREE="mgc/cli/cli-dump-tree.json"
OUT_DIR="mgc/cli/docs"

set -xe

cd $SPECS_DIR
go build -o cicd
cd ..
cd ..

cd $MGCDIR
go build -tags \"embed\" -o mgc
cd ..
cd ..


echo "generating $OUT_DIR..."
$SPECS_DIR"cicd" pipeline cligendoc -c $MGCDIR"mgc" -d "$DUMP_TREE" -o "$OUT_DIR"
echo "generating $OUT_DIR: done"

echo "ENDING $me"
