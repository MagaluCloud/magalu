#!/bin/bash
me=$(basename "$0")
echo "RUNNING $me"

MGCDIR=${MGCDIR:-"mgc/cli/"}
SPECS_DIR=${SPECS_DIR:-"mgc/spec_manipulator/"}
OUT_FILE="mgc/cli/cli-dump-tree.json"

#set -xe

cd $SPECS_DIR
go build -o cicd
cd ..
cd ..

cd $MGCDIR
go build -tags \"embed\" -o mgc
cd ..
cd ..

echo "generating $OUT_FILE..."
$SPECS_DIR"cicd" pipeline dumptree -c $MGCDIR"mgc" -o "$OUT_FILE"
echo "generating $OUT_FILE: done"

echo "ENDING $me"
