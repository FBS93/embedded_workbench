#!/usr/bin/env bash
set -e

# Check if the correct number of arguments is provided
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 [PRESET] [SOURCE_DIR]"
    exit 1
fi

# Parse command-line arguments
preset=$1
source_dir=$2

# Clean build artifacts
build_dir="${source_dir}/build/${preset}"
rm -rf "$build_dir"

echo "✅ Clean completed."
