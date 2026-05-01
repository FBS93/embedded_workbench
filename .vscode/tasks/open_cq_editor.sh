#!/usr/bin/env bash
set -e

echo "🧊 Open CQ-Editor"

# Validate required commands.
if ! command -v cq-editor >/dev/null 2>&1; then
    echo "❌ Error: cq-editor not found."
    exit 1
fi

exec cq-editor
