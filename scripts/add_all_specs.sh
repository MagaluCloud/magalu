#!/bin/sh

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
OAPIEMBED=${OAPIEMBED:-"mgc/sdk/openapi/embed_loader.go"}
OAPI_PATH=$ROOTDIR/$OAPIDIR

set -xe

$BASEDIR/add_specs.sh block-storage https://block-storage.br-ne-1.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh dbaas https://dbaas.br-ne-1.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh mke https://mke.br-ne-1.jaxyendy.com/docs/openapi-with-snippets.json

# This file is NOT being used, the API is not recommended and we should follow with their S3 compatible API
# $BASEDIR/add_specs.sh object-storage https://object-storage.br-ne-1.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh virtual-machine https://virtual-machine.br-ne-1.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh vpc https://vpc.br-ne-1.jaxyendy.com/openapi.json

python3 $BASEDIR/oapi_index_gen.py "--embed=$OAPIEMBED" $OAPI_PATH
