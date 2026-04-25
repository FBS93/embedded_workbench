# Embedded Workbench (EW) overview

EW is a template for embedded systems development, providing a fully containerized environment with all dependencies preinstalled for fast setup and reproducibility.

EW includes:
- V-model development methodology. See [development_methodology_overview.md](development_methodology/development_methodology_overview.md).
- AI-assisted SDD (Spec Driven Development). See [ai_assisted_sdd.md](resources/ai_assisted_sdd.md).
- Dual targeting setup. See [dual_targeting_setup.md](resources/dual_targeting_setup.md).
- Technical notes related to embedded systems. See [technical_notes](technical_notes/) folder.
- VS Code workflow. See [vscode_workflow.md](resources/vscode_workflow.md).
- Template repository synchronization support. See [template_sync.md](resources/template_sync.md).
- Third-party repository synchronization support. See [third_party_sync_workflow.md](resources/third_party_sync_workflow.md).

All user-configurable settings of the template are exposed through environment variables defined in the [devcontainer.json](../.devcontainer/devcontainer.json) file.

# Glossary

| Term | Definition |
|------|------------|
| | |

# Usage example

The following resources provide practical entry points for using EW:
- Getting started. See [getting_started.md](resources/getting_started.md).

@todo Add automated fuzz testing based on EDF event-sequence injection:
- As a standalone robustness-testing tool, separate from u_test / i_test / q_test.
- As a host-only EBF PORT specialized for fuzzing, instead of a new EBF CORE or an additional EDF AO.
- Integrating AFL++ **persistent mode** into the runtime entry point of that host fuzzing port (in the EDF kernel for this dedicated port).
- Injecting fuzzed inputs through the real EDF event APIs.
- Requiring all application events to be defined in a single dedicated library with a fixed convention, so the fuzz harness can be auto-generated.
- Following a single-root-event injection pattern: inject one fuzzed event, let EDF dispatch it and all derived events until quiescence, then continue with the next fuzzed event.
- Auto-generating the fuzz harness from the central event definitions.
- Keeping the fuzz harness deterministic: the order and number of injected events are defined by the AFL++ input, while the harness itself performs no internal randomization. Use the same input to randomize the event payload when needed.
- Redefining `EAF_onError()` during fuzz campaigns to classify controlled asserts separately from critical failures.
- Treating crashes, hangs, memory corruption, and uncontrolled framework failures as critical findings.
