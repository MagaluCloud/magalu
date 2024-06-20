#!/bin/sh

BASEDIR=$(dirname $0)
# ROOTDIR=$(builtin cd $BASEDIR/..; pwd)

OAPIDIR=${OAPIDIR:-"mgc/cli/openapis"}
OAPIEMBED=${OAPIEMBED:-"mgc/sdk/openapi/embed_loader.go"}
OAPI_PATH=$ROOTDIR/$OAPIDIR

set -xe

$BASEDIR/add_specs.sh block-storage volume mgc/spec_manipulator/cli_specs/conv.block-storage.jaxyendy.openapi.json https://block-storage.jaxyendy.com/openapi.json
echo "BLOCK-STORAGE"

$BASEDIR/add_specs.sh dbaas database mgc/spec_manipulator/cli_specs/database.jaxyendy.openapi.json https://dbaas.jaxyendy.com/openapi.json
echo "DATABASE"

$BASEDIR/add_specs.sh kubernetes kubernetes mgc/spec_manipulator/cli_specs/kubernetes.jaxyendy.openapi.json https://mke.br-ne-1.com/docs/openapi-with-snippets.json
echo "KUBERNETES"

# # This file is NOT being used, the API is not recommended and we should follow with their S3 compatible API
# # $BASEDIR/add_specs.sh object-storage objects https://object-storage.br-ne-1.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh virtual-machine compute mgc/spec_manipulator/cli_specs/conv.virtual-machine.jaxyendy.openapi.json https://virtual-machine.jaxyendy.com/openapi.json
echo "VIRTUAL MACHINE"

# # $BASEDIR/add_specs.sh virtual-machine-xaas compute-xaas mgc/spec_manipulator/cli_specs/virtual-machine-xaas.jaxyendy.openapi.json https://virtual-machine.jaxyendy.com/internal/v1/openapi.json
# # echo "VIRTUAL MACHINE XAAS"

$BASEDIR/add_specs.sh network network mgc/spec_manipulator/cli_specs/conv.network.jaxyendy.openapi.json https://network.jaxyendy.com/openapi.json
echo "NETWORK"

$BASEDIR/add_specs.sh container-registry container-registry mgc/spec_manipulator/cli_specs/container-registry.openapi.yaml https://container-registry.jaxyendy.com/openapi.json
echo "REGISTRY"
