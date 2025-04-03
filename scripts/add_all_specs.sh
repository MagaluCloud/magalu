#!/bin/sh

BASEDIR=$(dirname $0)

OAPIDIR=${OAPIDIR:-"mgc/sdk/openapi/openapis"}
OAPI_PATH=$ROOTDIR/$OAPIDIR

set -ue

$BASEDIR/add_specs.sh block-storage volume specs/conv.block-storage.jaxyendy.openapi.json https://block-storage.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh dbaas database specs/database.jaxyendy.openapi.yaml https://dbaas.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh kubernetes kubernetes specs/kubernetes.jaxyendy.openapi.json https://mke.br-ne-1.com/docs/openapi-with-snippets.json

$BASEDIR/add_specs.sh virtual-machine compute specs/conv.virtual-machine.jaxyendy.openapi.json https://virtual-machine.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh network network specs/conv.network.jaxyendy.openapi.json https://network.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh container-registry container-registry specs/container-registry.openapi.yaml https://container-registry.jaxyendy.com/openapi.json

$BASEDIR/add_specs.sh audit audit specs/conv.events-consult.openapi.yaml https://events-consult.jaxyendy.com/openapi-cli.json

$BASEDIR/add_specs_without_region.sh profile profile specs/conv.globaldb.openapi.yaml https://globaldb.jaxyendy.com/openapi-cli.json

$BASEDIR/add_specs.sh load-balancer load-balancer specs/lbaas.openapi.yaml https://lbaas.jaxyendy.com/openapi-cli.json

# EXAMPLE
# $BASEDIR/SCRIPT.sh NOME_NO_MENU URL_PATH LOCAL_DA_SPEC SPEC_UID
