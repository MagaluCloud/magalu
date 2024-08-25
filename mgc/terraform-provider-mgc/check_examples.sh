#!/bin/bash

script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

search_dir1=$script_dir/docs/data-sources
search_dir2=$script_dir/docs/resources

# Counter for files with errors
error_count=0

# Function to check a file
check_file() {
    local file="$1"
    if ! grep -q '^## Example Usage$' "$file"; then
        echo "Error: Missing '## Example Usage' in $file"
        ((error_count++))
    fi
}

# Function to process a directory
process_directory() {
    local dir="$1"
    echo "Checking files in $dir"
    while IFS= read -r -d '' file; do
        check_file "$file"
    done < <(find "$dir" -type f \( -name "*.md" -o -name "*.markdown" \) -print0)
}

# Process both directories
process_directory "$search_dir1"
process_directory "$search_dir2"

# Print summary
echo "Finished checking files in both directories."
echo "Found $error_count file(s) missing '## Example Usage'."

# Exit with error if any files were missing the required line
if [ $error_count -gt 0 ]; then
    exit 1
else
    exit 0
fi
