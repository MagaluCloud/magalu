#!/bin/sh

set -xe

ENTRYPOINT=${ENTRYPOINT:-mgc/cli/main.go}
BUILDDIR=${BUILDDIR:-build}
NAME=${NAME:-mgc}
VERSION=${VERSION}
TAGS=${TAGS:-embed}
LDFLAGS=${LDFLAGS:-"-s -w"}
DESIRED_DIST_REGEXP=${DESIRED_DIST_REGEXP:-"^\(linux\|darwin\|windows\)/\(amd64\|arm64\)"}

if [ -z "$VERSION" ]; then
    VERSION=`git log -1 --pretty=format:%h`
fi

mkdir -p "$BUILDDIR"

for D in `go tool dist list | grep "$DESIRED_DIST_REGEXP"`; do
    OS=`echo "$D" | cut -d/ -f1`
    ARCH=`echo "$D" | cut -d/ -f2`
    GOOS="$OS" GOARCH="$ARCH" go build -tags "$TAGS" -ldflags "$LDFLAGS" -o "$BUILDDIR/$NAME-$OS-$ARCH-$VERSION" "$ENTRYPOINT"
done
