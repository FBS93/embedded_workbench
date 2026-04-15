#!/usr/bin/env bash
set -e

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 [WORKSPACE_ROOT]"
    exit 1
fi

workspace_root="$1"

cd "${workspace_root}"

echo "💾 Saving KiCad state ..."
python ".devcontainer/save-kicad-state.py"
echo "✅ KiCad state saved."
