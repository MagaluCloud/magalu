#!/bin/bash

set -xe

script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")
echo $script_path
PROVIDER_DIR=$script_dir
TF_PLUGIN_DOCS="github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4"
ABSOLUTE_PATH_OAPI=$(realpath $script_dir/../cli/openapis)
ABSOLUTE_PATH_BLUEPRINTS=$(realpath $script_dir/../cli/blueprints)

mkdir -p $script_dir/docs

MGC_SDK_OPENAPI_DIR=$ABSOLUTE_PATH_OAPI MGC_SDK_BLUEPRINTS_DIR=$ABSOLUTE_PATH_BLUEPRINTS go run $TF_PLUGIN_DOCS generate --provider-dir="$PROVIDER_DIR"

sleep 5

echo "Adding subcategories"
. $script_dir/add_subcategory.sh

echo "Moving extra docs"
cp -r $script_dir/docs-extra/. $script_dir/docs
