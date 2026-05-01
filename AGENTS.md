# Embedded Workbench skills

Skills for the Embedded Workbench. Load the smallest skill that fits the task.

The development methodology adopted by the Embedded Workbench follows a V-model approach. The specific methodology and processes to be followed are defined in [development_methodology_overview.md](doc/development_methodology/development_methodology_overview.md).

Some skills represent V-model processes, while others are supporting skills.

## V-Model process skills

These skills map directly to development processes and should be used following the V-model methodology.

### Specifications

| Skill | Purpose | Path |
|---|---|---|
| `specs-definition` | Define project specifications from stakeholder specifications. | skills/specs-definition/SKILL.md |

### System domain

| Skill | Purpose | Path |
|---|---|---|
| `sys-requirements` | Define system requirements from specifications. | skills/sys-requirements/SKILL.md |
| `sys-architecture` | Define and implement system qualification tests against system requirements to validate that the system implementation correctly realizes them. | skills/sys-architecture/SKILL.md |
| `sys-integration-test` | Define and implement system integration tests against system architecture to validate that the system implementation correctly realizes it. | skills/sys-integration-test/SKILL.md |
| `sys-qualification-test` | Validate system against system requirements | skills/sys-qualification-test/SKILL.md |

### Software domain

| Skill | Purpose | Path |
|---|---|---|
| `sw-requirements` | Define software requirements from system requirements, system architecture and hardware-software interface. | skills/sw-requirements/SKILL.md |
| `sw-architecture` | Define software architecture from software requirements. | skills/sw-architecture/SKILL.md |
| `sw-detailed-design` | Define software detailed design from software requirements and software architecture. | skills/sw-detailed-design/SKILL.md |
| `sw-implementation` | Implement the software from software requirements, software architecture and software detailed design. | skills/sw-implementation/SKILL.md |
| `sw-unit-test` | Define and implement software unit tests against software detailed design to validate that the software implementation correctly realizes it. | skills/sw-unit-test/SKILL.md |
| `sw-integration-test` | Define and implement software integration tests against software architecture to validate that the software detailed design and software implementation correctly realizes it. | skills/sw-integration-test/SKILL.md |
| `sw-qualification-test` | Define and implement software qualification tests against software requirements to validate that the software architecture, software detailed design, and software implementation correctly realizes them. | skills/sw-qualification-test/SKILL.md |

### Hardware domain

| Skill | Purpose | Path |
|---|---|---|
| `hw-requirements` | Define hardware requirements from system requirements, system architecture and mechanical-hardware interface. | skills/hw-requirements/SKILL.md |
| `hw-architecture` | Define hardware architecture from hardware requirements. | skills/hw-architecture/SKILL.md |
| `hw-implementation` | Implement the hardware design from hardware requirements and hardware architecture. | skills/hw-implementation/SKILL.md |
| `hw-integration-test` | Define and implement hardware integration tests against hardware architecture to validate that the hardware implementation correctly realizes it. | skills/hw-integration-test/SKILL.md |
| `hw-qualification-test` | Define and implement hardware qualification tests against hardware requirements to validate that the hardware architecture and hardware implementation correctly realize them. | skills/hw-qualification-test/SKILL.md |

### Mechanical domain

| Skill | Purpose | Path |
|---|---|---|
| `mech-requirements` | Define mechanical requirements from system requirements and system architecture. | skills/mech-requirements/SKILL.md |
| `mech-architecture` | Define mechanical architecture from mechanical requirements. | skills/mech-architecture/SKILL.md |
| `mech-implementation` | Implement the mechanical design from mechanical requirements and mechanical architecture. | skills/mech-implementation/SKILL.md |
| `mech-integration-test` | Define and implement mechanical integration tests against mechanical architecture to validate that the mechanical implementation correctly realizes it. | skills/mech-integration-test/SKILL.md |
| `mech-qualification-test` | Define and implement mechanical qualification tests against mechanical requirements to validate that the mechanical architecture and mechanical implementation correctly realize them. | skills/mech-qualification-test/SKILL.md |

## Supporting skills

These skills support the development process but are not tied to a specific V-model process.

| Skill | Purpose | Path |
|---|---|---|
| `sys-review` | Review system process work products before closure. | skills/sys-review/SKILL.md |
| `sw-review` | Review software process work products before closure. | skills/sw-review/SKILL.md |
| `hw-review` | Review hardware process work products before closure. | skills/hw-review/SKILL.md |
| `mech-review` | Review mechanical process work products before closure. | skills/mech-review/SKILL.md |
| `sw-debug` | Debug software issues by investigating problems and deriving accurate conclusions based on evidence to enable correct resolution. | skills/sw-debug/SKILL.md |
| `sw-fuzz-corpus` | Analyze EDF fuzzable events to generate an AFL++ corpus. | skills/sw-fuzz-corpus/SKILL.md |
| `workbench-debug` | Debug and resolve development environment issues by investigating problems and deriving accurate conclusions based on evidence. | skills/workbench-debug/SKILL.md |
| `technical-note` | Create technical notes for relevant embedded engineering findings. | skills/technical-note/SKILL.md |
