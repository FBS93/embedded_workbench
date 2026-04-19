# KiCad workflow

## Overview

This document describes the KiCad state persistence strategy used by the hardware domain to keep the KiCad user environment reproducible inside the devcontainer.

## State persistence strategy

KiCad uses environment variables to configure the paths where it reads and writes user state. These paths are configured in [.devcontainer/devcontainer.json](../../../../../.devcontainer/devcontainer.json):

`KICAD_CONFIG_HOME` and `KICAD_DOCUMENTS_HOME` store KiCad configuration and user documents respectively. These paths point to a Docker volume that is not versioned with the workspace. To preserve the KiCad environment across fresh clones of the repository and devcontainer rebuilds, the project stores a base state containing the initial KiCad configuration, installed plugins, and required libraries. This base state is restored only when the runtime volume is empty. Once restored, the user continues working against the persistent runtime volume. This avoids versioning volatile metadata and other files that change frequently during normal KiCad usage and would otherwise add noise to the repository. When the base state needs to be updated, the scripts documented in [tools/kicad_state/kicad_state.md](../../../../tools/kicad_state/kicad_state.md) can capture the current runtime state or restore the stored base state.

`KICAD10_3RD_PARTY` stores third-party libraries, plugins, and related resources. This directory points inside the workspace so it can be versioned. It is important to keep this directory under version control because it contains the third-party KiCad libraries required by the project, and those libraries shall evolve together with the hardware design.

The persistence flow used to keep the KiCad environment stable across fresh clones and devcontainer rebuilds is:
1. The devcontainer `postStartCommand` runs [`post_start.sh`](../../../../.devcontainer/post_start.sh).
2. [`post_start.sh`](../../../../.devcontainer/post_start.sh) prepares the persistent volume mounted at `/kicad-runtime`.
3. [`post_start.sh`](../../../../.devcontainer/post_start.sh) runs [`restore_kicad_state.py`](../../../../tools/kicad_state/restore_kicad_state.py).
4. [`restore_kicad_state.py`](../../../../tools/kicad_state/restore_kicad_state.py) restores the versioned state only when the runtime is still empty.
5. KiCad then works directly against the persistent runtime paths.

## VS Code tasks

The following [VS Code tasks](../../../../.vscode/tasks.json) are available to execute manual actions of the KiCad state workflow:

- `Open KiCad`: opens KiCad against the already prepared runtime.
- `Save KiCad state`: captures the current runtime state to be used in subsequent restores.
- `Restore KiCad state`: force-resets the runtime from the stored state.
- `Clean KiCad runtime state`: removes the current runtime contents from `KICAD_CONFIG_HOME` and `KICAD_DOCUMENTS_HOME` without touching `KICAD10_3RD_PARTY`.

## Notes

- The automatic restore happens only on container start.
- The runtime is not overwritten automatically once it already contains content.
