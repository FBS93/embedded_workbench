#!/usr/bin/env bash
set -e

echo "🔎 Run linter"

# Parse command-line arguments.
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 [SOURCE_DIR] [WORKSPACE_ROOT] [PRESET]"
    exit 1
fi

source_dir="$1"
workspace_root="$2"
preset="$3"
build_dir="${workspace_root}/build/${preset}"
log_file="${build_dir}/lint_report.txt"
compile_commands="${build_dir}/compile_commands.json"
header_filter="^${source_dir}/.*"
lint_status=0

# Validate required inputs.
if [ ! -d "${source_dir}" ]; then
    echo "❌ Error: source directory not found: ${source_dir}"
    exit 1
fi

if [ ! -d "${workspace_root}" ]; then
    echo "❌ Error: workspace directory not found: ${workspace_root}"
    exit 1
fi

if [ ! -f "${workspace_root}/.clang-tidy" ]; then
    echo "❌ Error: .clang-tidy file not found: ${workspace_root}/.clang-tidy"
    exit 1
fi

# Validate required commands.
if ! command -v clang-tidy >/dev/null 2>&1; then
    echo "❌ Error: clang-tidy not found."
    exit 1
fi

# Generate compile_commands.json for the selected preset.
cmake --preset "${preset}" --no-warn-unused-cli

if [ ! -f "${compile_commands}" ]; then
    echo "❌ Error: compile_commands.json not found: ${compile_commands}"
    exit 1
fi

# Clean previous log.
: "${preset:?}"
mkdir -p "${build_dir}"
: > "${log_file}"

found_files=0

while IFS= read -r -d '' file; do
    found_files=1
    echo "🔍 Analyzing: ${file}"
    printf '🔍 Analyzing: %s\n' "${file}" >> "${log_file}"
    clang-tidy "${file}" -p "${build_dir}" --config-file="${workspace_root}/.clang-tidy" --header-filter="${header_filter}" >> "${log_file}" 2>&1 || lint_status=$?
done < <(find "${source_dir}" -type f -name '*.c' -print0)

# Check if there are files to analyze.
if [ "${found_files}" -eq 0 ]; then
    echo "⚠️ No files found to analyze." | tee "${log_file}"
    exit 1
fi

echo "Linting completed."
echo "Report: ${log_file}"

exit "${lint_status}"
