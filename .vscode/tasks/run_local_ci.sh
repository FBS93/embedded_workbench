#!/usr/bin/env bash
set -euo pipefail

echo "🚦 Run Local CI"

# Parse command-line arguments.
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 [WORKSPACE_ROOT] [PRESET...]"
    exit 1
fi

workspace_root="${1%/}"
shift
presets=("$@")

if [ ! -d "${workspace_root}" ]; then
    echo "❌ Error: workspace directory not found: ${workspace_root}"
    exit 1
fi

for preset in "${presets[@]}"; do
    if [ -z "${preset}" ]; then
        echo "❌ Error: no CMake preset selected."
        exit 1
    fi
done

cd "${workspace_root}"

total_steps=$((2 + ${#presets[@]}))
step=1

echo "[${step}/${total_steps}] Checking formatting..."
./.vscode/tasks/format.sh "${workspace_root}/sw" "${workspace_root}" --check

step=$((step + 1))
echo "[${step}/${total_steps}] Running lint/build for ${#presets[@]} preset(s)..."
./.vscode/tasks/lint.sh "${workspace_root}" "${presets[@]}"

for preset in "${presets[@]}"; do
    step=$((step + 1))
    echo "[${step}/${total_steps}] Running CTest for preset '${preset}'..."
    ctest --test-dir "build/${preset}" --output-on-failure \
        --test-output-size-passed 0 \
        --test-output-size-failed 0
done

echo "✅ Local CI passed."
