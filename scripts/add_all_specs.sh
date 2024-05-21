#!/bin/sh

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
OAPIEMBED=${OAPIEMBED:-"mgc/sdk/openapi/embed_loader.go"}
OAPI_PATH=$ROOTDIR/$OAPIDIR

set -xe

$BASEDIR/add_specs.sh block-storage volume mgc/api-specs/cli_specs/block-storage.viveiro.openapi.yaml https://block-storage.jaxyendy.com/openapi.json
echo "BLOCK-STORAGE"

$BASEDIR/add_specs.sh dbaas database mgc/api-specs/cli_specs/database.viveiro.openapi.yaml https://dbaas.jaxyendy.com/openapi.json
echo "DATABASE"

$BASEDIR/add_specs.sh kubernetes kubernetes mgc/magalu/mgc/api-specs/cli_specs/kubernetes.viveiro.openapi.yaml https://mke.br-ne-1.com/docs/openapi-with-snippets.json
echo "KUBERNETES"

$BASEDIR/add_specs.sh virtual-machine compute mgc/magalu/mgc/api-specs/cli_specs/virtual-machine.viveiro.openapi.yaml https://virtual-machine.jaxyendy.com/openapi.json
echo "VIRTUAL MACHINE"

$BASEDIR/add_specs.sh network network mgc/magalu/mgc/api-specs/cli_specs/network.viveiro.openapi.yaml https://network.jaxyendy.com/openapi.json
echo "NETWORK"

# # This file is NOT being used, the API is not recommended and we should follow with their S3 compatible API
# # $BASEDIR/add_specs.sh object-storage objects https://object-storage.br-ne-1.jaxyendy.com/openapi.json
