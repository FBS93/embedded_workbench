# hw-implementation

## Purpose

Implement the hardware design from hardware requirements and hardware architecture.

## Input work products

- Hardware requirements
- Hardware architecture
- Tools

## Output work products

- Hardware implementation
- Hardware-mechanical interface
- Hardware-software interface

## Steps

1. Review the hardware requirements and hardware architecture.
2. Review available tools required for the implementation.
3. Implement each hardware assembly according to its hardware architecture definition.
4. Ensure consistency between the implementation and the hardware architecture and hardware requirements.
5. Check the hardware implementation for completeness, consistency, and correctness.
6. Define the hardware-mechanical interface according to the hardware architecture and hardware implementation.
7. Ensure consistency between hardware-mechanical interface, hardware architecture and hardware implementation.
8. Check the hardware-mechanical interface for completeness, consistency, and correctness.
9. Define the hardware-software interface according to the hardware architecture and hardware implementation.
10. Ensure consistency between hardware-software interface, hardware architecture and hardware implementation.
11. Check the hardware-software interface for completeness, consistency, and correctness.

## Guidelines

### Hardware implementation work product

Each hardware assembly shall be implemented strictly according to its hardware architecture definition and shall not introduce additional constraints.

Hardware requirements shall be used only as complementary input when the hardware architecture does not fully define the hardware assembly.

The implementation of each hardware assembly shall:
- Be organized in a dedicated root folder whose name matches the hardware assembly defined in the architecture. This establishes an implicit 1:1 relationship between the hardware assembly and its implementation.

#### KiCad library integration

The hardware implementation work product shall use the KiCad libraries defined in [KiCad libraries for JLCPCB and LCSC](../resources/kicad_libraries_for_jlcpcb_lcsc.md).

#### Schematic implementation

The schematic implementation of each hardware assembly shall:
- Be written in Python using [SKiDL](https://skidl.readthedocs.io/).
- Follow the rules defined in [skidl_guidelines.md](../resources/skidl_guidelines.md).
- Be implemented in the following source files:
  - A single top-level hardware assembly Python source file, named after the hardware assembly in lower_snake_case.
  - One hardware component Python source file for each hardware component belonging to the hardware assembly, named after the corresponding hardware component in lower_snake_case.
- Realize the referenced hardware interfaces and hardware parameters only through the corresponding hardware component implementation to which they are linked in the hardware architecture.
- Generate all required exported files from the source files and never edit exported files manually.
- Use the top-level assembly source file to generate the netlist file (`.net`).
- Use the top-level assembly source file to generate the whole assembly schematic file (`.svg`).
- Use each hardware component source file to generate the corresponding component schematic file (`.svg`).

#### PCB implementation

When a hardware assembly requires a PCB, the PCB implementation shall:
- Be implemented manually by a human in [KiCad](https://www.kicad.org/).
- Not be modified directly by AI-assisted tooling.
- Be derived from the generated top-level assembly source file netlist.
- Be derived from PCB outline constraints defined in upstream traceable elements, if any.
- Preserve the connectivity defined by the generated netlist. If connectivity changes are required, the netlist shall never be modified manually; the change shall always come from the source files used to generate the netlist.
- Contain manufacturing-oriented configuration required for the hardware assembly.

The hardware implementation shall not be considered complete until the human responsible for the PCB implementation explicitly confirms that the PCB implementation is complete.

AI-assisted tooling may continue with other tasks that does not depend on the final PCB implementation. However, no activity depending on the final PCB implementation shall be considered complete until the required human confirmation is provided.

### Hardware-mechanical interface work product

The hardware-mechanical interface work product shall follow the [Cross-domain interface definition](../../resources/cross_domain_interface_definition.md).

The hardware-mechanical interface work product shall be defined using as input the traceable elements of the hardware architecture and hardware implementation that are relevant to the mechanical domain.

The hardware-mechanical interface work product may include, for example, when applicable:
- PCB outline constraints.
- Mounting hole definitions and keep-out constraints.
- Connector locations and accessibility constraints.
- A 3D model of the PCB including all assembled components.

### Hardware-software interface work product

The hardware-software interface work product shall follow the [Cross-domain interface definition](../../resources/cross_domain_interface_definition.md).

The hardware-software interface work product shall be defined using as input the traceable elements of the hardware architecture and hardware implementation that are relevant to the software domain.

The hardware-software interface work product may include, for example, when applicable:
- Programmable device interfaces and peripherals relevant to firmware implementation.
- Reference manuals and datasheets for hardware components relevant to firmware implementation.
- External interrupt sources and the conditions that trigger them.
- Boot, reset, programming, and debug interfaces.
