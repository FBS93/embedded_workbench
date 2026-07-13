#!/usr/bin/env bash
set -e

echo "🚦 Run local CI"

# Parse command-line arguments.
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 [WORKSPACE_ROOT] [ACT_ARGS...]"
    exit 1
fi

workspace_root="${1%/}"
shift

ci_workflow_file="${workspace_root}/.github/workflows/ci.yml"
hil_workflow_file="${workspace_root}/.github/workflows/ci_hil.yml"
job_name="ci"
act_platform_image="${ACT_PLATFORM_IMAGE:-catthehacker/ubuntu:act-24.04}"
hil_env_args=(
    --env "RPI_USER=${RPI_USER-}"
    --env "RPI_HOST=${RPI_HOST-}"
    --env "GDB_PORT=${GDB_PORT-}"
    --env "LOG_PORT=${LOG_PORT-}"
    --env "LOG_BAUD_RATE=${LOG_BAUD_RATE-}"
    --env "LOG_SERIAL_DEVICE=${LOG_SERIAL_DEVICE-}"
)
act_workspace_args=()

docker_workspace_contains_repo() {
    local docker_workspace_root="$1"

    docker run --rm \
        --mount "type=bind,source=${docker_workspace_root},target=/mnt/workspace,readonly" \
        "${act_platform_image}" \
        test -f /mnt/workspace/.devcontainer/devcontainer.json >/dev/null 2>&1
}

resolve_docker_workspace_root() {
    local docker_workspace_root="${LOCAL_CI_DOCKER_WORKSPACE_ROOT:-${workspace_root}}"
    local container_id
    local mounted_workspace_root

    if docker_workspace_contains_repo "${docker_workspace_root}"; then
        printf '%s\n' "${docker_workspace_root}"
        return 0
    fi

    container_id="$(hostname)"
    mounted_workspace_root="$({
        docker inspect \
            --format "{{range .Mounts}}{{if eq .Destination \"${workspace_root}\"}}{{.Source}}{{end}}{{end}}" \
            "${container_id}" 2>/dev/null || true
    })"

    if [ -n "${mounted_workspace_root}" ] && docker_workspace_contains_repo "${mounted_workspace_root}"; then
        printf '%s\n' "${mounted_workspace_root}"
        return 0
    fi

    return 1
}

# Validate required inputs.
if [ ! -d "${workspace_root}" ]; then
    echo "❌ Error: workspace directory not found: ${workspace_root}"
    exit 1
fi

if [ ! -f "${ci_workflow_file}" ]; then
    echo "❌ Error: CI workflow not found: ${ci_workflow_file}"
    exit 1
fi

if [ ! -f "${hil_workflow_file}" ]; then
    echo "❌ Error: HIL CI workflow not found: ${hil_workflow_file}"
    exit 1
fi

# Validate required commands.
if ! command -v act >/dev/null 2>&1; then
    echo "❌ Error: act not found. Rebuild the devcontainer to install it."
    exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
    echo "❌ Error: docker not found. act requires Docker access."
    exit 1
fi

if ! docker info >/dev/null 2>&1; then
    echo "❌ Error: Docker daemon not reachable. act requires access to a running Docker daemon."
    exit 1
fi

if ! docker_workspace_root="$(resolve_docker_workspace_root)"; then
    echo "❌ Error: Docker cannot bind-mount the repository workspace."
    echo "   Reopen/rebuild the devcontainer with Docker-outside-of-Docker enabled, or set"
    echo "   LOCAL_CI_DOCKER_WORKSPACE_ROOT to the Docker-host-visible repository path."
    exit 1
fi

if [ "${docker_workspace_root}" != "${workspace_root}" ]; then
    act_workspace_args=(
        --container-options "--mount type=bind,source=${docker_workspace_root},target=${docker_workspace_root}"
        --env "LOCAL_CI_CHECKOUT_PATH=${docker_workspace_root}"
    )
fi

cd "${workspace_root}"

echo "CI workflow: ${ci_workflow_file}"
echo "HIL CI workflow: ${hil_workflow_file}"
echo "Job: ${job_name}"
echo "Act platform image: ${act_platform_image}"
echo "Docker workspace root: ${docker_workspace_root}"
act --version

echo "Running local CI ..."
act workflow_dispatch -W "${ci_workflow_file}" -j "${job_name}" -P "ubuntu-24.04=${act_platform_image}" "${act_workspace_args[@]}" "$@"

echo "Running local HIL CI ..."
act workflow_dispatch -W "${hil_workflow_file}" -j "${job_name}" -P "self-hosted=${act_platform_image}" "${act_workspace_args[@]}" "${hil_env_args[@]}" "$@"
