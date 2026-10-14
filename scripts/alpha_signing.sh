#!/bin/bash

set -e

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
KEYS=$ROOTDIR/mgc/sdk/preview/alphasync/trusted_keys.json

BIN_DIR=$(mktemp -d)
trap 'rm -rf "$BIN_DIR"' EXIT

(builtin cd $ROOTDIR/mgc/sdk && go build -o "$BIN_DIR/alpha-publishing" ./preview/alphasync/cmd/alpha-publishing)

if [ $# -eq 0 ]; then
    "$BIN_DIR/alpha-publishing"
fi

"$BIN_DIR/alpha-publishing" "$1" --keys "$KEYS" "${@:2}"
