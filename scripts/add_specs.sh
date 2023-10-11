#!/bin/bash

set -e

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIEMBED=${OAPIEMBED:-"mgc/sdk/openapi/embed_loader.go"}
OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
CUSTOM_DIR=${CUSTOM_DIR:-"openapi-customizations"}

OAPI_PATH=$ROOTDIR/$OAPIDIR
CUSTOM_PATH=$ROOTDIR/$CUSTOM_DIR

API_NAME=$1
API_SPEC_FILE=$2
SPEC_UID=$3
SPEC_FILE="$API_NAME.openapi.yaml"

if ! test -f $CUSTOM_PATH/$SPEC_FILE; then
    cat > $CUSTOM_PATH/$SPEC_FILE << EOF
# This file is to be merged on top of $OAPIDIR/$SPEC_FILE
# using yaml_merge.py
# NOTE: Lists are merged by their indexes, be careful with parameters, tags and such!
# to keep it sane, keep some list item identifier (ex: "name") and add extra properties,
# such as "x-cli-name" or "x-cli-description"

servers:
-   url: https://api-$API_NAME{env}.{region}.jaxyendy.com
    variables:
        region:
            description: Region to reach the service
            default: br-ne-1
            enum:
            - br-ne-1
            - br-ne-2
            - br-se-1
        env:
            description: Environment to use
            default: ''
            enum:
            - ''
            - .pre-prod
            x-cli-transforms:
            -   type: translate
                translations:
                -   from: prod
                    to: ''
                -   from: pre-prod
                    to: .pre-prod

EOF
fi

python3 $BASEDIR/transformers/transform.py $API_SPEC_FILE $SPEC_UID -o $OAPI_PATH/$SPEC_FILE
python3 $BASEDIR/yaml_merge.py --override $OAPI_PATH/$SPEC_FILE $CUSTOM_PATH/$SPEC_FILE
python3 $BASEDIR/oapi_index_gen.py "--embed=$OAPIEMBED" $OAPI_PATH
