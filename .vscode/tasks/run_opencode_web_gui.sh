#!/usr/bin/env bash
set -e

HOSTNAME="0.0.0.0"
PORT="4096"

if ! command -v opencode >/dev/null 2>&1; then
    echo "❌ Error: opencode command not found."
    exit 1
fi

count_opencode_credentials() {
    local auth_output
    local auth_output_clean
    local credential_count

    auth_output="$(opencode auth list 2>/dev/null || true)"
    auth_output_clean="$(printf '%s\n' "${auth_output}" | sed -E 's/\x1B\[[0-9;]*[[:alpha:]]//g')"
    credential_count="$(printf '%s\n' "${auth_output_clean}" | grep -Eo '[0-9]+ credentials' | tail -n 1 | cut -d' ' -f1)"

    if [[ "${credential_count}" =~ ^[0-9]+$ ]]; then
        printf '%s\n' "${credential_count}"
        return 0
    fi

    printf '0\n'
}

ask_user() {
    local message="$1"
    local reply

    while true; do
        read -r -p "${message} [y/N]: " reply
        case "${reply}" in
            [yY]|[yY][eE][sS])
                return 0
                ;;
            ""|[nN]|[nN][oO])
                return 1
                ;;
            *)
                echo "Invalid answer. Enter 'y' or 'n'."
                ;;
        esac
    done
}

# Check whether OpenCode already has at least one logged-in provider.
CREDENTIAL_COUNT="$(count_opencode_credentials)"

if [[ "${CREDENTIAL_COUNT}" -eq 0 ]]; then
    echo "No active OpenCode provider session found."

    while ask_user "Open a provider login before starting the web GUI?"; do
        echo "Starting OpenCode provider login..."
        opencode auth login

        CREDENTIAL_COUNT="$(count_opencode_credentials)"
        if [[ "${CREDENTIAL_COUNT}" -gt 0 ]]; then
            echo "✅ Detected ${CREDENTIAL_COUNT} active OpenCode credential(s)."
        else
            echo "No new credentials detected."
        fi

        if [[ "${CREDENTIAL_COUNT}" -gt 0 ]]; then
            if ! ask_user "Open another provider login?"; then
                break
            fi
        fi
    done
else
    echo "✅ Detected ${CREDENTIAL_COUNT} active OpenCode credential(s)."
fi

echo "Starting OpenCode web on ${HOSTNAME}:${PORT} ..."

exec opencode web --hostname "${HOSTNAME}" --port "${PORT}"
