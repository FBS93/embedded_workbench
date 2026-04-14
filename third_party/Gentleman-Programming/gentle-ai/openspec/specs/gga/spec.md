# GGA Specification

## Purpose

Defines the install and runtime behavior of the GGA component, covering the PowerShell shim asset and its Windows-specific install step.

## Requirements

### Requirement: PowerShell Shim Asset

The system MUST embed a `gga.ps1` file as a Go asset under `internal/assets/gga/`. The shim MUST delegate execution to the Git Bash binary resolved by `gitBashPath()`, forwarding all arguments verbatim and propagating the exit code.

#### Scenario: Shim delegates to Git Bash

- GIVEN the embedded `gga.ps1` is installed on a Windows machine with Git Bash present
- WHEN the user runs `gga <subcommand>` from PowerShell
- THEN the shim invokes Git Bash with the resolved bash binary path and all supplied arguments
- AND the process exits with the same code returned by the underlying GGA bash command

#### Scenario: Arguments containing spaces are forwarded correctly

- GIVEN `gga.ps1` is installed
- WHEN the user runs `gga commit -m "my message"` from PowerShell
- THEN the argument `"my message"` reaches GGA as a single token (not split)

#### Scenario: Exit code propagation on error

- GIVEN `gga.ps1` is installed
- WHEN the underlying GGA command exits with a non-zero code
- THEN PowerShell's `$LASTEXITCODE` reflects that exact non-zero value

---

### Requirement: Windows Install Step

On Windows, the installer MUST write `gga.ps1` to the same directory as the GGA bash script after GGA's own `install.sh` completes. The write MUST use an atomic no-op pattern: if the file already exists with identical content, the installer MUST NOT overwrite it.

#### Scenario: First-time install on Windows

- GIVEN GGA has completed its own install
- AND `gga.ps1` does not yet exist in the install directory
- WHEN the Windows install step runs
- THEN `gga.ps1` is written to the install directory with correct content

#### Scenario: Idempotent re-install (content unchanged)

- GIVEN `gga.ps1` already exists with content matching the current embedded asset
- WHEN the installer runs again
- THEN the file is NOT overwritten (no write I/O occurs)

#### Scenario: Stale shim is updated

- GIVEN `gga.ps1` exists but its content differs from the current embedded asset
- WHEN the installer runs
- THEN the file is atomically replaced with the new content

#### Scenario: Git Bash not found at install time

- GIVEN Git Bash is not installed on the target Windows machine
- WHEN the install step attempts to resolve `gitBashPath()`
- THEN the installer surfaces a clear, actionable error message
- AND installation halts without writing a broken shim

---

### Requirement: Non-Windows Systems Unaffected

On non-Windows platforms, the installer MUST NOT attempt to write `gga.ps1` or invoke the PowerShell shim step.

#### Scenario: Linux/macOS install flow unchanged

- GIVEN a Linux or macOS host
- WHEN GGA install runs
- THEN no `.ps1` file is created and no Windows-specific code path executes

---

## Doc Note

`docs/platforms.md` MUST remove any Windows limitation note that states PowerShell is unsupported once this change ships. This is a documentation-only update with no behavioral requirement beyond keeping the doc accurate.
