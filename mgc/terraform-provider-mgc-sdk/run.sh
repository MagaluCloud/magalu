#!/bin/bash

# Build newest provider
go build || exit

export TF_LOG="${TF_LOG:-info}"

tf_exec=terraform
tf_args="${@:2}"

if [[ "$MGC_OPENTF" ]]; then
    tf_exec=tofu
fi

if [ -z $tf_args ]; then
    echo "Invalid arguments."
    echo "Usage: $0 <example_folder> <command>"
    echo ""
    echo "Possible commands:"
    echo -e "apply:\tExecutes \`$tf_exec apply\`"
    echo -e "plan:\tExecutes \`$tf_exec plan\`"
    echo -e "clear:\tExecutes \`$tf_exec clear <module>\`"
else
    (cd examples/$1 && $tf_exec "$tf_args")
fi
