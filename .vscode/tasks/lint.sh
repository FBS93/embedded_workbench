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
ruff_bin=""
ruff_config="${workspace_root}/pyproject.toml"

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

if command -v ruff >/dev/null 2>&1; then
    ruff_bin="$(command -v ruff)"
fi

if [ -z "${ruff_bin}" ] || [ ! -x "${ruff_bin}" ]; then
    echo "❌ Error: Ruff linter not found."
    exit 1
fi

if [ ! -f "${ruff_config}" ]; then
    echo "❌ Error: Ruff linter config file not found: ${ruff_config}"
    exit 1
fi

echo "CMake preset: ${preset}"
echo "Ruff linter: ${ruff_bin}"

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

echo "Running Python lint with Ruff ..."
if "${ruff_bin}" check --config "${ruff_config}" "${workspace_root}"; then
    echo "✅ Python lint successful."
else
    echo "❌ Python lint failed."
    exit 1
fi
