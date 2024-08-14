#!/bin/bash

set -e

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
LIBDIR=${LIBDIR:-"$ROOTDIR/mgc/lib"}

export MGC_SDK_OPENAPI_DIR=${MGC_SDK_OPENAPI_DIR:-"$ROOTDIR/mgc/cli/openapis"}
export MGC_SDK_BLUEPRINTS_DIR=${MGC_SDK_BLUEPRINTS_DIR:-"$ROOTDIR/mgc/cli/blueprints"}
BRANCH="$(git rev-parse --abbrev-ref HEAD)"

# if [[ "$BRANCH" != "main" ]]; then
#   echo 'Aborting script';
#   exit 0;
# fi

## Only generate if we're on a Git tag
# if ! git describe --tags --exact-match HEAD >/dev/null 2>&1; then
#     if [[ "$VERSION" != "" ]]; then
        (rm -rf $LIBDIR/products)
        (cd $ROOTDIR/mgc/codegen; go build -tags "embed release" -o codegen; ./codegen $LIBDIR)
        (cd $LIBDIR; go mod tidy)
#     fi
# fi
