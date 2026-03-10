#!/usr/bin/env bash
set -e

# Parse command-line arguments
# Input:
#  - SOURCE_DIR: directory where recursive formatting starts.
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 [SOURCE_DIR]"
    exit 1
fi

source_dir="$1"
clang_format_bin=""
script_path="$(realpath "${BASH_SOURCE[0]}")"
script_dir="$(dirname "$script_path")"
workspace_root="$(realpath "${script_dir}/../..")"
asm_formatter="${workspace_root}/tools/scripts/asm_format.py"
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

echo "🧼 Formatting..."
echo "clang-format: ${clang_format_bin}"
echo "ASM formatter: ${asm_formatter}"

cd "${source_dir}"

# Format C/C++ source files recursively.
find . -type f \( -name '*.c' -o -name '*.h' -o -name '*.cpp' -o -name '*.hpp' \) \
    -exec "${clang_format_bin}" -i {} +

# Format assembler source files recursively.
find . -type f -name '*.S' -exec "${asm_formatter}" {} +

echo "✅ Formatting completed."
