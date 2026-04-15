# KiCad state in devcontainer

## Purpose

Define how KiCad user state is versioned and restored inside the devcontainer so the hardware environment remains reproducible across rebuilds.

## Scope

The versioned KiCad state is stored in:

```text
.devcontainer/kicad-state/<version>/
```

Each version directory may contain only the following:

- `config/kicad_common.json`
- `config/sym-lib-table`
- `config/fp-lib-table`
- `config/user.hotkeys`
- `config/design-block-lib-table`
- `config/colors/`
- `data/plugins/`
- `data/template/`
- `custom-libraries/`

`colors/`, `plugins/`, `template/`, and `custom-libraries/` shall contain a `.gitkeep` file only when the directory would otherwise be empty.

Caches, logs, backups, and temporary KiCad files shall not be versioned.

## Workflow

- `.devcontainer/setup-kicad.py` restores the versioned KiCad state into `~/.config/kicad/<version>/` and `~/.local/share/KiCad/<version>/`.
- `.devcontainer/save-kicad-state.py` captures the current KiCad user state back into `.devcontainer/kicad-state/<version>/`.
- `.vscode/tasks/save_kicad_state.sh` is the VS Code task wrapper used by the hardware toolbox button.
- `.devcontainer/devcontainer.json` runs `python3 .devcontainer/setup-kicad.py` through `postCreateCommand` so the state is restored automatically after container creation.

## Custom libraries

Custom KiCad libraries that are not already located inside the repository shall be moved or copied into:

```text
.devcontainer/kicad-state/<version>/custom-libraries/
```

After that, the KiCad library tables shall be updated to point to the repo-managed paths.

The save script warns when the KiCad tables still reference custom libraries outside the repository.
