---
name: hw-qualification-test
description: >
  Define and implement hardware qualification tests against hardware requirements to validate that the hardware architecture and hardware implementation correctly realize them.
  Trigger: When hardware integration test is completed and hardware qualification tests shall be defined and implemented.
---

## When to Use

- After hardware integration test is completed.
- When any hardware qualification test activity is required.

## Critical Patterns

- Shall strictly follow the provided resources and recursively review all referenced documents required to understand the current process, inputs, outputs, rules and directly referenced work products.
- Shall align with hardware requirements (no inconsistencies).
- Shall ensure all hardware qualification tests are derived from hardware requirements by adding explicit traceability between them.
- Shall define and implement the hardware qualification tests against the hardware requirements to validate that the hardware architecture and hardware implementation correctly realize them.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall ask for approval on key hardware qualification test decisions.
- Shall set produced/modified work products status to `not approved`.

## Resources

- [development_methodology_overview.md](../../../doc/development_methodology/development_methodology_overview.md)
- [hardware_domain_overview.md](../../../doc/development_methodology/hardware_domain/hardware_domain_overview.md)
- [hw-qualification-test.md](../../../doc/development_methodology/hardware_domain/processes/hw-qualification-test.md)
