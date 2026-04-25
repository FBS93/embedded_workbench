# VS Code workflow

## Overview

VS Code running the provided devcontainer is the official development entry point for the repository.

The available VS Code functions are organized by domain. Each function is listed in this document, while the corresponding domain-specific documentation remains the source of truth for its detailed usage.

The repository exposes these functions through three VS Code entry point types:
- Tasks defined in [tasks.json](../../.vscode/tasks.json).
- Launch configurations defined in [launch.json](../../.vscode/launch.json).
- Task buttons configured in [settings.json](../../.vscode/settings.json), grouped as workbench, software, hardware, and mechanical toolboxes.

## Workbench

| Function | VS Code entry point | Purpose | Reference |
|---|---|---|---|
| OpenCode web GUI | `run_opencode_web_gui` task, `🌐 OpenCode web GUI` task button | Opens the OpenCode web interface for repository AI workflows on port `4096`. | [AI-assisted SDD](ai_assisted_sdd.md) |
| Raspberry Pi SSH access | `connect_rpi_ssh` task, `🔗 Connect Raspberry Pi (SSH)` task button | Opens an SSH session to the Raspberry Pi used by the repository workflow. | [Raspberry Pi setup](raspberry_pi_setup.md) |
| Third-party sync | `run_third_party_sync` task, `🔄 Sync third-party repositories` task button | Runs the local third-party synchronization workflow. | [Third-party sync workflow](third_party_sync_workflow.md), [third_party_sync tool](../../tools/third_party_sync/third_party_sync.md) |

## Software domain

| Function | VS Code entry point | Purpose | Reference |
|---|---|---|---|
| Build | `build` task, `🔨 Build` task button | Builds the selected software target using the active CMake preset and active CMake build target. | |
| Clean | `clean` task, `🧹 Clean` task button | Cleans the build artifacts generated for the active CMake preset. | |
| Format | `format` task, `🧼 Format all source files` task button | Formats repository C, ASM, CMake, and Python sources using the configured repository tooling. ASM files are also formatted on save through VS Code settings. | [Embedded C coding guidelines](../development_methodology/software_domain/resources/embedded_c_coding_guidelines.md), [Python coding guidelines](../development_methodology/resources/python_coding_guidelines.md), [ASM formatter tool](../development_methodology/software_domain/resources/asm_format.md) |
| Lint | `lint` task, `🔎 Run linter` task button | Runs the configured software static analysis workflow using the active CMake preset. | |
| Host debug | `Host` launch configuration | Launches host-side debugging for the active CMake launch target. | |
| Target debug | `Target STM32F103C8Tx (launch/attach)` launch configurations, `run_target_gdb_server` task | Launches or attaches a remote debug session on the target platform. | [Embedded target remote debugging](../development_methodology/software_domain/resources/embedded_target_remote_debugging.md) |
| Target logging | `run_target_logging_server` task, `📃 Enable target logging` task button | Starts the remote logging workflow for the target platform. | [Embedded target remote logging](../development_methodology/software_domain/resources/embedded_target_remote_logging.md) |

## Hardware domain

| Function | VS Code entry point | Purpose | Reference |
|---|---|---|---|
| KiCad | `open_kicad` task, `⚡ Open KiCad` task button | Opens KiCad using the repository workflow and runtime configuration. | [KiCad workflow](../development_methodology/hardware_domain/resources/kicad_workflow.md) |
| KiCad state save | `save_kicad_state` task, `💾 Save KiCad state` task button | Saves the tracked KiCad runtime state used by the repository workflow. | [KiCad workflow](../development_methodology/hardware_domain/resources/kicad_workflow.md), [KiCad state](../../tools/kicad_state/kicad_state.md) |
| KiCad state restore | `restore_kicad_state` task, `♻️ Restore KiCad state` task button | Restores the tracked KiCad runtime state used by the repository workflow. | [KiCad workflow](../development_methodology/hardware_domain/resources/kicad_workflow.md), [KiCad state](../../tools/kicad_state/kicad_state.md) |
| KiCad runtime cleanup | `clean_kicad_runtime_state` task, `🧹 Clean KiCad runtime state` task button | Cleans the KiCad runtime state used by the repository workflow. | [KiCad workflow](../development_methodology/hardware_domain/resources/kicad_workflow.md), [KiCad state](../../tools/kicad_state/kicad_state.md) |

## Mechanical domain

| Function | VS Code entry point | Purpose | Reference |
|---|---|---|---|
| CQ-Editor | `open_cq_editor` task, `🧊 Open CQ-Editor` task button | Opens CQ-Editor for interactive CadQuery model inspection and editing. | |

## Expected extensions

The following VS Code extensions are declared for the devcontainer in [devcontainer.json](../../.devcontainer/devcontainer.json).

The repository also expects VS Code Remote - Containers support on the host machine in order to open the workspace through the provided devcontainer.

### Host prerequisite

- `ms-vscode-remote.vscode-remote-extensionpack`

### Devcontainer - Workbench
- `spencerwmiles.vscode-task-buttons`
- `mhutchie.git-graph`
- `bierner.markdown-mermaid`
- `bierner.markdown-preview-github-styles`

### Devcontainer - Software domain
- `ms-vscode.cmake-tools`
- `ms-vscode.cpptools`
- `dan-c-underwood.arm`
- `ms-python.python`
- `marus25.cortex-debug`
- `ms-vscode.vscode-serial-monitor`
- `emeraldwalk.runonsave`

### Devcontainer - Hardware and mechanical domains

No additional VS Code extension is currently declared specifically for the repository hardware or mechanical workflows beyond the common extensions and the tools provided by the devcontainer.
