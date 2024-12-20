#!/bin/bash

set -xe

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
SDKDIR=${SDKDIR:-"mgc/sdk"}
OAPIDIR=${OAPIDIR:-"$SDKDIR/openapi/openapis"}
OAPICUSTOMDIR=${OAPICUSTOMDIR:-"openapi-customizations"}

python3 $ROOTDIR/scripts/oapi_apply_customizations.py $OAPIDIR $OAPICUSTOMDIR
