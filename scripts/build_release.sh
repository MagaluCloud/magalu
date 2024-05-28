#!/bin/bash

set -xe

ENTRYPOINT=${ENTRYPOINT:-mgc/cli/main.go}
BUILDDIR=${BUILDDIR:-build}
NAME=${NAME:-mgc}
VERSION=${VERSION}
TAGS=${TAGS:-embed release}
CGO_ENABLED=${CGO_ENABLED:-0}
LDFLAGS=${LDFLAGS:-"-s -w"}
# DESIRED_DIST_REGEXP=${DESIRED_DIST_REGEXP:-"^\(linux\|darwin\|windows\)/\(amd64\|arm64\)"}
DESIRED_DIST_REGEXP=${DESIRED_DIST_REGEXP:-"^\(linux\)/\(amd64\)"}
CHECK_RELEASE=${CHECK_RELEASE:-1}
TERRAFORM_ENTRYPOINT=${TERRAFORM_ENTRYPOINT:-mgc/terraform-provider-mgc/main.go}
TERRAFORM_NAME=${TERRAFORM_NAME:-terraform-provider-mgc}
TERRAFORM_SRCDIR=${TERRAFORM_SRCDIR:-$(dirname $TERRAFORM_ENTRYPOINT)}

if [ -z "$VERSION" ]; then
    VERSION=`git log -1 '--pretty=format:%(describe:tags)'`
fi

if [ $CHECK_RELEASE -eq 1 ]; then
    ./scripts/check_release.sh
fi

source ./scripts/tf_generate_docs.sh

rm -rf "$BUILDDIR"
mkdir -p "$BUILDDIR"

for D in `go tool dist list | grep "$DESIRED_DIST_REGEXP"`; do
    OS=`echo "$D" | cut -d/ -f1`
    ARCH=`echo "$D" | cut -d/ -f2`
    EXT=`if [ "$OS" = "windows" ]; then echo ".exe"; fi`

    SUBDIR="$BUILDDIR/$NAME-cli-$OS-$ARCH-$VERSION"
    mkdir -p "$SUBDIR"

    # BUILD CLI
    GOOS="$OS" GOARCH="$ARCH" go build -buildvcs=false -tags "$TAGS" -ldflags "$LDFLAGS -X magalu.cloud/sdk.Version=$VERSION" -o "$SUBDIR/$NAME$EXT" "$ENTRYPOINT"

    # Copy cli additional files
    cp mgc/cli/RUNNING.md "$SUBDIR/README.md"
    cp -a mgc/cli/examples "$SUBDIR"
    cp mgc/sdk/openapi/README.md "$SUBDIR/OPENAPI.md"
    cp -r share "$SUBDIR"

    SUBDIR="$BUILDDIR/$NAME-terraform-$OS-$ARCH-$VERSION"
    mkdir -p "$SUBDIR"

    # Build Terraform provider
    GOOS="$OS" GOARCH="$ARCH" go build -buildvcs=false -tags "$TAGS" -ldflags "$LDFLAGS -X magalu.cloud/sdk.Version=$VERSION" -o "$SUBDIR/$TERRAFORM_NAME$EXT" "$TERRAFORM_ENTRYPOINT"

    # Copy terraform additional files
    cp -r $TERRAFORM_SRCDIR/docs "$SUBDIR"
    cp -r $TERRAFORM_SRCDIR/user-guide.md "$SUBDIR"
    cp -r $TERRAFORM_SRCDIR/install.sh "$SUBDIR"
done
