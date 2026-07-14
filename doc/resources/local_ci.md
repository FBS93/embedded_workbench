# Local CI

## Overview

This document describes how Local CI validates changes in the devcontainer before they are pushed to configured remote branches.

Local CI can also validate the active CMake preset on demand through the corresponding VS Code entry point. See [VS Code workflow](vscode_workflow.md).

## Architecture

- The pre-push hook runs Local CI for non-deletion updates to configured remote branches, validating each configured CMake preset.
- The VS Code task validates the active CMake configure preset on demand.

Git push --> Pre-push hook --> Local CI runner --> Configured CMake presets

VS Code entry point --> Local CI runner --> Active CMake configure preset

## Local CI Setup

Local CI is executed by [run_local_ci.sh](../../.vscode/tasks/run_local_ci.sh).

This script:

- Checks repository formatting once without modifying source files.
- Runs clang-tidy and a clean build for each requested preset in order.
- Runs Ruff once.
- Runs CTest for each requested preset in order.

The [pre-push hook](../../.githooks/pre-push) executes the script with the configured presets. The hook is a local safeguard, not a server-side access control mechanism, and can be bypassed with `git push --no-verify`.

Presets containing Hardware-in-the-Loop (HiL) tests require the hardware setup to be configured and operational. See [Raspberry Pi setup](raspberry_pi_setup.md) and [dual targeting setup](dual_targeting_setup.md). Software-in-the-Loop (SiL) presets do not require embedded target hardware.

## Dependencies

### Devcontainer

Make sure that the environment variables in [devcontainer.json](../../.devcontainer/devcontainer.json) are configured for Local CI:

- `LOCAL_CI_BRANCHES` lists the exact remote branch names, separated by spaces, that trigger pre-push Local CI.
- `LOCAL_CI_PRESETS` lists the CMake presets, separated by spaces, that pre-push Local CI validates.

The devcontainer configures the [pre-push hook](../../.githooks/pre-push) through `core.hooksPath` when the container is created.
