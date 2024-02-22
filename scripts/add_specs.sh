#!/bin/bash

set -e

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
CUSTOM_DIR=${CUSTOM_DIR:-"openapi-customizations"}

OAPI_PATH=$ROOTDIR/$OAPIDIR
CUSTOM_PATH=$ROOTDIR/$CUSTOM_DIR

API_NAME=$1
API_ENDPOINT_NAME=$2
API_SPEC_FILE=$3
SPEC_UID=$4
SPEC_FILE="$API_NAME.openapi.yaml"

if ! test -f $CUSTOM_PATH/$SPEC_FILE; then
    cat > $CUSTOM_PATH/$SPEC_FILE << EOF
# This file is to be merged on top of $OAPIDIR/$SPEC_FILE
# using yaml_merge.py
# NOTE: Lists are merged by their indexes, be careful with parameters, tags and such!
# to keep it sane, keep some list item identifier (ex: "name") and add extra properties,
# such as "x-mgc-name" or "x-mgc-description"

servers:
-   url: https://{env}/{region}/$API_ENDPOINT_NAME
    variables:
        region:
            description: Region to reach the service
            default: br-ne-1
            enum:
            - br-ne-1
            - br-se-1
            - br-mgl-1
            x-mgc-transforms:
            -   type: translate
                translations:
                -   from: br-ne1
                    to: br-ne-1
                -   from: br-se1
                    to: br-se-1
                -   from: br-mgl1
                    to: br-mgl-1
        env:
            description: Environment to use
            default: ''
            enum:
            - api.magalu.cloud
            - api.pre-prod.jaxyendy.com
            x-mgc-transforms:
            -   type: translate
                translations:
                -   from: prod
                    to: api.magalu.cloud
                -   from: pre-prod
                    to: api.pre-prod.jaxyendy.com

EOF
fi

python3 $BASEDIR/transformers/transform.py $API_NAME $API_SPEC_FILE $SPEC_UID -o $OAPI_PATH/$SPEC_FILE
python3 $BASEDIR/yaml_merge.py --override $OAPI_PATH/$SPEC_FILE $CUSTOM_PATH/$SPEC_FILE
$BASEDIR/oapi_index_gen.sh
