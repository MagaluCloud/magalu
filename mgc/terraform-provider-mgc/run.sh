#!/bin/bash

# Build newest provider
go build || exit

TFLOG=info

case $2 in
    plan)
        (cd examples/$1 && terraform plan)
        ;;
    apply)
        (cd examples/$1 && terraform apply)
        ;;
    clear)
        (cd examples/$1 && terraform state rm "$3")
        ;;
    *)
        echo "Invalid arguments."
        echo "Usage: $0 <example_folder> <command>"
        echo ""
        echo "Possible commands:"
        echo -e "apply:\tExecutes \`terraform apply\`"
        echo -e "plan:\tExecutes \`terraform plan\`"
        echo -e "clear:\tExecutes \`terraform clear <module>\`"
        ;;
esac
