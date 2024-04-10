#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Version not found"
    exit 1
fi

version="$1"

if ! [[ $version =~ ^v?[0-9]+(\.[0-9]+){2,3}(-[a-zA-Z0-9]+)?$ ]];
then
  echo "Version not available for release. have=$version"
  exit 1
fi
