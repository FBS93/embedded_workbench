#!/usr/bin/env bash
set -euo pipefail

echo "🔎 Run linter"

# Parse command-line arguments.
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 [WORKSPACE_ROOT] [PRESET]"
    exit 1
fi

workspace_root="${1%/}"
preset="$2"
build_dir="${workspace_root}/build/${preset}"

# Validate required inputs.
if [ -z "${preset}" ]; then
    echo "❌ Error: no CMake preset selected."
    exit 1
fi

if [ ! -d "${workspace_root}" ]; then
    echo "❌ Error: workspace directory not found: ${workspace_root}"
    exit 1
fi

# Validate required commands.
if ! command -v cmake >/dev/null 2>&1; then
    echo "❌ Error: cmake not found."
    exit 1
fi

echo "CMake preset: ${preset}"

cd "${workspace_root}"

echo "Configuring with clang-tidy enabled ..."
# Using --no-warn-unused-cli flag because this: https://gitlab.kitware.com/cmake/cmake/-/issues/17261
if cmake --preset "${preset}" -DEW_ENABLE_CLANG_TIDY=ON --no-warn-unused-cli; then
    echo "✅ CMake configuration successful."
else
    echo "❌ CMake configuration failed."
    exit 1
fi

echo "Building preset '${preset}' with clang-tidy enabled ..."
if cmake --build "${build_dir}" --clean-first; then
    echo "✅ Lint build successful."
else
    echo "❌ Lint build failed."
    exit 1
fi
