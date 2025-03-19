#!/bin/bash

set -ue

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
OAPIDIR=${OAPIDIR:-"mgc/sdk/openapi/openapis"}

python3 $BASEDIR/oapi_index_gen.py $OAPIDIR 
