# CI workflow

## Overview

EW provides customizable CI through GitHub Actions workflow template files.

The CI infrastructure is split into two workflows:
- [ci.yml](../../.github/workflows/ci.yml) defines the `CI workflow`.
- [ci_hil.yml](../../.github/workflows/ci_hil.yml) defines the `HIL CI workflow`.

The workflow files are the source of truth for exact CI steps.

## Cloud CI

Cloud CI runs only the `CI workflow` on:
- `push` to `main`.
- `pull_request` targeting `main`.
- Manual `workflow_dispatch`.

Cloud CI cannot run HIL tests because they depend on the embedded target hardware.

## Local CI

Local CI runs both CI workflows, `CI workflow` and `HIL CI workflow`, with `act` through the corresponding VS Code entry point. See [vscode_workflow.md](vscode_workflow.md).

## Customization

For a derived project, customize CI by editing only the GitHub Actions workflow template files defined in the [Overview](#overview) chapter.
