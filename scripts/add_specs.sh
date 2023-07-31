#!/bin/bash

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
OAPIDIR=$ROOTDIR/$1

if ! test -d $OAPIDIR;
then
    echo "Directory provided does not exist or empty: $OAPIDIR"
    exit
else
    echo "Saving OpenAPI Specs on $OAPIDIR"
fi


# TODO: make it as a list and derive the name from the URL (regex, maybe)
API_LIST=(
    # Add new APIs in this format: key|url
    "dbaas|https://dbaas.br-ne-1.jaxyendy.com/openapi.json"
    "mke|https://mke.br-ne-1.jaxyendy.com/docs/openapi-with-snippets.json"
    "object-storage|https://object-storage.br-ne-1.jaxyendy.com/openapi.json"
)

if ! command -v python3 &> /dev/null
then
    echo "python3 dependency not found. Install before running the script"
    exit
fi

for api in "${API_LIST[@]}"; do
    KEY="${api%%|*}"
    VALUE="${api##*|}"
    FILENAME=$KEY.openapi.yaml
    echo "Working on >> $VALUE"
    if python3 $BASEDIR/sync_oapi.py $VALUE -o $OAPIDIR/$FILENAME;
    then
        # Remove tenant id
        echo "... Removing any x-tenant-id references"
        python3 $BASEDIR/remove_tenant_id.py $OAPIDIR/$FILENAME
        # Check for customizations
        if test -f $ROOTDIR/openapi-customizations/$FILENAME;
        then
            echo "... Applying customizations"
            python3 $BASEDIR/yaml_merge.py $OAPIDIR/$FILENAME $ROOTDIR/openapi-customizations/$FILENAME
        fi
        echo "... Saved spec in $OAPIDIR/$FILENAME"
    fi
done
