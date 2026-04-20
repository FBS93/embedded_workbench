# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment
- This project uses `pip` for dependency management.
- Tests are managed using `pytest`.
- `tox` is used to run tests across multiple environments (different Python/KiCad versions).

## Common Commands
- **Run all tests (default environment)**: `pytest tests`
- **Run tests across all supported test environments**: `tox`
- **Run a specific test**: Navigate to the `tests` directory and use `pytest path/to/test.py`
- **Build the package**: `python setup.py sdist`
- **Clean build artifacts**: `rm dist/*`

## Project Architecture
SKiDL acts as an infrastructure-as-code tool for circuit design, converting Python-based
circuit descriptions into netlists for PCB layout tools (primarily KiCad).

### Key Modules
- `src/skidl/`: Contains the core package logic.
  - Core circuitry elements like `Part`, `Net`, and `Circuit` definitions.
  - Netlist generation and ERC (Electrical Rules Checking) logic.
  - `schematics`: Logic for generating schematics from SKiDL.
  - `scripts`: User-facing CLI utilities.
  - `tools`: Backend interfaces from SKiDL to various EDA packages.
- `tests/`: Extensive test suite for functionalities ranging from basic circuit construction to hierarchical schematic generation and SPICE integration.
  - `unit_tests`: Unit tests, both manually and AI-generated.
    - `ai_tests`: AI-generated unit tests.
  - `test_data`: Data (mostly part libraries) needed to run unit tests. 
  - `examples/`: Examples to test various features.

### Schematic Generation Integration
For KiCad integration (specifically KiCad 6-9), the process involves:
- Symbol definition extraction from draw commands.
- Hierarchical UUID generation and multi-file schematic output.
- Force-directed placement and routing algorithms to handle component positioning and connectivity.
- Coordinate system handling (KiCad is Y-down, requiring transformations).
