#!/usr/bin/env bash
set -e

# Parse command-line arguments
# Input:
#  - SOURCE_DIR: directory where recursive formatting starts.
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 [SOURCE_DIR] [WORKSPACE_ROOT]"
    exit 1
fi

source_dir="$1"
workspace_root="$2"
clang_format_bin=""
ruff_bin=""
asm_formatter="${workspace_root}/tools/scripts/asm_format.py"
ruff_config="${workspace_root}/pyproject.toml"
# Resolve clang-format binary:
# 1) Use clang-format from PATH.
# 2) Fallback to VS Code extension bundled binary.
if command -v clang-format >/dev/null 2>&1; then
    clang_format_bin="$(command -v clang-format)"
else
    ext_bin="$(
        find "$HOME/.vscode-server/extensions" "$HOME/.vscode/extensions" \
            -type f -name 'clang-format' 2>/dev/null \
            | sort -V | tail -n 1
    )"
    clang_format_bin="${ext_bin}"
fi

if [ ! -d "${source_dir}" ]; then
    echo "❌ Error: source directory not found: ${source_dir}"
    exit 1
fi

if [ -z "${clang_format_bin}" ] || [ ! -x "${clang_format_bin}" ]; then
    echo "❌ Error: clang-format not found."
    exit 1
fi

if [ ! -x "${asm_formatter}" ]; then
    echo "❌ Error: ASM formatter script not found or not executable: ${asm_formatter}"
    exit 1
fi

if command -v ruff >/dev/null 2>&1; then
    ruff_bin="$(command -v ruff)"
fi

if [ -z "${ruff_bin}" ] || [ ! -x "${ruff_bin}" ]; then
    echo "❌ Error: ruff not found."
    exit 1
fi

if [ ! -f "${ruff_config}" ]; then
    echo "❌ Error: Ruff config file not found: ${ruff_config}"
    exit 1
fi

echo "🧼 Formatting..."
echo "clang-format: ${clang_format_bin}"
echo "ASM formatter: ${asm_formatter}"
echo "ruff: ${ruff_bin}"

cd "${source_dir}"

# Format C/C++ source files recursively.
find . -type f \( -name '*.c' -o -name '*.h' -o -name '*.cpp' -o -name '*.hpp' \) \
    -exec "${clang_format_bin}" -i {} +

# Format assembler source files recursively.
find . -type f -name '*.S' -exec "${asm_formatter}" {} +

# Format Python files across the workspace using the repository Ruff config.
"${ruff_bin}" check --fix --config "${ruff_config}" "${workspace_root}"
"${ruff_bin}" format --config "${ruff_config}" "${workspace_root}"

echo "✅ Formatting completed."
