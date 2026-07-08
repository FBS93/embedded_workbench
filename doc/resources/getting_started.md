# Getting started

## Overview

This document describes the shortest path to start working with the repository.

The expected development entry point is VS Code running the provided devcontainer.

## Host prerequisites

The host environment shall provide the following:
- VS Code with the `ms-vscode-remote.vscode-remote-extensionpack` extension installed.
- OpenCode installed. The devcontainer persists OpenCode state (including credentials and sessions) through the bind mount defined in [devcontainer.json](../../.devcontainer/devcontainer.json).

## Initial setup

1. Open the repository in VS Code.
2. Run `Dev Containers: Reopen in Container` from the VS Code command palette.
3. Wait until the devcontainer build and startup process completes.

The devcontainer definition is provided in [../../.devcontainer](../../.devcontainer/).
