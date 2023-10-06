#!/bin/bash

set -xe

ENTRYPOINT=${ENTRYPOINT:-mgc/cli/main.go}
BUILDDIR=${BUILDDIR:-build}
NAME=${NAME:-mgc}
VERSION=${VERSION}
TAGS=${TAGS:-embed}
CGO_ENABLED=${CGO_ENABLED:-0}
LDFLAGS=${LDFLAGS:-"-s -w"}
DESIRED_DIST_REGEXP=${DESIRED_DIST_REGEXP:-"^\(linux\|darwin\|windows\)/\(amd64\|arm64\)"}


if [ -z "$VERSION" ]; then
    VERSION=`git log -1 --pretty=format:%h`
fi

source ./scripts/tf_generate_docs.sh

mkdir -p "$BUILDDIR"

for D in `go tool dist list | grep "$DESIRED_DIST_REGEXP"`; do
    OS=`echo "$D" | cut -d/ -f1`
    ARCH=`echo "$D" | cut -d/ -f2`
    EXT=`if [ "$OS" = "windows" ]; then echo ".exe"; fi`
    GOOS="$OS" GOARCH="$ARCH" go build -tags "$TAGS" -ldflags "$LDFLAGS" -o "$BUILDDIR/$NAME-$OS-$ARCH-$VERSION$EXT" "$ENTRYPOINT"
done

cp mgc/cli/RUNNING.md "$BUILDDIR/README.md"
cp -a mgc/cli/openapis "$BUILDDIR"
cp -a mgc/cli/examples "$BUILDDIR"
cp mgc/sdk/openapi/README.md "$BUILDDIR/OPENAPI.md"
cp -r mgc/terraform-provider-mgc/docs "$BUILDDIR"
