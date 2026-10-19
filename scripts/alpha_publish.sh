#!/bin/bash

set -euo pipefail

BASEDIR=$(dirname $0)
ROOTDIR=$(builtin cd $BASEDIR/..; pwd)
KEYS=$ROOTDIR/mgc/sdk/preview/alphasync/trusted_keys.json
MGC=${MGC:-mgc}
BUCKET=${BUCKET:-mgc-cli-preview}
BASE_URL=${BASE_URL:-https://$BUCKET.br-se1.magaluobjects.com}

if [ $# -eq 0 ] || { [ "$1" != grant ] && [ "$1" != revoke ]; }; then
    echo "usage: $0 grant|revoke [flags]" >&2
    exit 2
fi
COMMAND=$1
shift

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
(builtin cd $ROOTDIR/mgc/sdk && go build -o "$WORK/alpha-publishing" ./preview/alphasync/cmd/alpha-publishing)

DRY_RUN=false
YES=false
ARGS=()
for arg in "$@"; do
    case $arg in
        --dry-run) DRY_RUN=true ;;
        --yes) YES=true ;;
        *) ARGS+=("$arg") ;;
    esac
done


status=$(curl -s -o /dev/null -w '%{http_code}' --max-time 15 "$BASE_URL/") || true
if [ "$status" != 403 ]; then
    echo "error: an anonymous listing of $BUCKET answered $status, not 403." >&2
    echo "If it lists its keys, make it private first:" >&2
    echo "  $MGC object-storage buckets acl set --dst $BUCKET --private" >&2
    exit 1
fi

"$WORK/alpha-publishing" "$COMMAND" --keys "$KEYS" --base-url "$BASE_URL" --out "$WORK/out" "${ARGS[@]}"

PLAN=$WORK/out/plan.tsv
if [ ! -s "$PLAN" ]; then
    echo "Nothing to upload."
    exit 0
fi

echo
echo "To upload to $BUCKET, in this order:"
cut -f2 "$PLAN" | sed 's/^/  /'

if $DRY_RUN; then
    echo "Dry run: nothing uploaded."
    exit 0
fi
if ! $YES; then
    read -r -p "Upload? [y/N] " answer
    [ "$answer" = y ] || { echo "Nothing uploaded."; exit 1; }
fi

while IFS=$'\t' read -r src dst; do
    "$MGC" object-storage objects upload --src "$src" --dst "$BUCKET/$dst" --no-confirm >/dev/null
    "$MGC" object-storage objects acl set --dst "$BUCKET/$dst" --public-read --no-confirm >/dev/null

    want=$(sha256sum < "$src" | cut -d' ' -f1)
    got=$(curl -sf --max-time 30 "$BASE_URL/$dst" | sha256sum | cut -d' ' -f1)
    if [ "$want" != "$got" ]; then
        echo "error: $BASE_URL/$dst is not public, or differs from what was uploaded" >&2
        exit 1
    fi
    echo "published $dst"
done < "$PLAN"
