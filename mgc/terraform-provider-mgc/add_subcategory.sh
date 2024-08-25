#!/bin/bash

script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

resources="$script_dir/docs/resources"
datasources="$script_dir/docs/data-sources"
jsonfile="$script_dir/subcategory.json"

if [ ! -d "$resources" ]; then
    echo "Directory does not exist"
    exit 1
fi

if [ ! -d "$datasources" ]; then
    echo "Directory does not exist"
    exit 1
fi

if [ ! -f "$jsonfile" ]; then
    echo "JSON file does not exist"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "jq is not installed. Please install it to parse JSON files."
    exit 1
fi

update_subcategory() {
    local file="$1"
    local category="$2"
    sed -i "s/subcategory: .*/subcategory: \"$category\"/" "$file"
    echo "Updated subcategory in $file to '$category'"
}

process_directory() {
    local directory="$1"
    local pattern="$2"
    local category="$3"
    
    for file in "$directory"/$pattern; do
        if [ -f "$file" ]; then
            update_subcategory "$file" "$category"
        fi
    done
}

jq -r 'to_entries | .[] | [.key, (.value | join(" "))] | @tsv' "$jsonfile" | while IFS=$'\t' read -r category patterns; do
    patterns=$(echo "$patterns" | tr -d '"')
    echo "Processing category: $category, patterns: $patterns"
    for pattern in $patterns; do
        process_directory "$resources" "$pattern" "$category"
        process_directory "$datasources" "$pattern" "$category"
    done
done