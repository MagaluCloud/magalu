#!/bin/sh

set -xe

PROVIDER_DIR="./mgc/terraform-provider-mgc/"
TF_PLUGIN_DOCS="github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4"
ABSOLUTE_PATH_OAPI=$(realpath ${MGC_SDK_OPENAPI_DIR:-./mgc/cli/openapis})
ABSOLUTE_PATH_BLUEPRINTS=$(realpath ${ABSOLUTE_PATH_BLUEPRINTS:-./mgc/cli/blueprints})

mkdir -p $PROVIDER_DIR/docs

MGC_SDK_OPENAPI_DIR=$ABSOLUTE_PATH_OAPI MGC_SDK_BLUEPRINTS_DIR=$ABSOLUTE_PATH_BLUEPRINTS go run $TF_PLUGIN_DOCS generate --provider-dir=$PROVIDER_DIR
