# AI-assisted development

## Overview

This document describes the repository-specific integration of OpenCode and Gentle AI Stack in the ECW devcontainer.

The general usage model, agent configuration, and provider setup are documented by the official OpenCode and Gentle AI Stack documentation and should be used as the primary reference for tool-specific behavior.

- OpenCode documentation: <https://opencode.ai/docs>
- Gentle AI Stack documentation: refer to the official project documentation for agent, preset, component, skill, and persona details.

## Repository integration

ECW provides a preconfigured AI-assisted workflow inside the devcontainer.

The repository-specific integration consists of the following:

- OpenCode is available inside the devcontainer.
- Gentle AI Stack is preconfigured for the `opencode` agent.
- Engram data is stored inside the workspace through `ENGRAM_DATA_DIR=${containerWorkspaceFolder}/.engram` in [devcontainer.json](../../.devcontainer/devcontainer.json).
- Local AI state directories `.engram`, `.atl`, and `.cache` are ignored by Git and shall not be committed.

## OpenCode web GUI

The OpenCode web interface can be started from the preconfigured VS Code task `run_opencode_web_gui` or from the toolbox entry `OpenCode web GUI`.

The launcher script [run_opencode_web_gui.sh](../../.vscode/tasks/run_opencode_web_gui.sh) performs the following repository-specific steps before starting the web interface:

- Verifies that the `opencode` command is available in the container.
- Checks `opencode auth list` to determine whether credentials are already configured.
- Can prompt the user to authenticate with one or more providers before starting the web interface.
- Starts the web interface with `opencode web --hostname 0.0.0.0 --port 4096`.

Once started, the web GUI listens on port `4096` inside the devcontainer.

## Notes

The devcontainer setup currently installs `gentle-ai` from the project release archive instead of the upstream `install.sh` helper because the helper currently fails during Docker builds with `tmpdir: unbound variable`.

This detail is an implementation note for the container image and should only matter when maintaining the devcontainer setup itself.
