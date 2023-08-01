#!/bin/bash

set -e

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
CUSTOM_DIR=${CUSTOM_DIR:-"openapi-customizations"}

OAPI_PATH=$ROOTDIR/$OAPIDIR
CUSTOM_PATH=$ROOTDIR/$CUSTOM_DIR

API_NAME=$1
API_URL=$2
SPEC_FILE="$API_NAME.openapi.yaml"

if ! test -f $CUSTOM_PATH/$SPEC_FILE; then
    cat > $CUSTOM_PATH/$SPEC_FILE << EOF
# This file is to be merged on top of $OAPIDIR/$SPEC_FILE
# using yaml_merge.py
# NOTE: Lists are merged by their indexes, be careful with parameters, tags and such!
# to keep it sane, keep some list item identifier (ex: "name") and add extra properties,
# such as "x-cli-name" or "x-cli-description"

servers:
-   url: https://api-$API_NAME.{region}.jaxyendy.com
    variables:
        region:
            description: Region to reach the service
            default: br-ne-1
            enum:
            - br-ne-1
            - br-ne-2
            - br-se-1
EOF
fi

python3 $BASEDIR/sync_oapi.py $API_URL --ext $OAPI_PATH/$SPEC_FILE
python3 $BASEDIR/remove_tenant_id.py $OAPI_PATH/$SPEC_FILE
python3 $BASEDIR/yaml_merge.py $OAPI_PATH/$SPEC_FILE $CUSTOM_PATH/$SPEC_FILE
