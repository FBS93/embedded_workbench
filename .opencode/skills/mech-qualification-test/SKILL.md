---
name: mech-qualification-test
description: >
  Define and implement mechanical qualification tests against mechanical requirements to validate that the mechanical architecture and mechanical implementation correctly realize them.
  Trigger: When mechanical integration test is completed and mechanical qualification tests shall be defined and implemented.
---

## When to Use

- After mechanical integration test is completed.
- When any mechanical qualification test activity is required.

## Critical Patterns

- Shall strictly follow the provided resources and recursively review all referenced documents required to understand the current process, inputs, outputs, rules and directly referenced work products.
- Shall align with mechanical requirements (no inconsistencies).
- Shall ensure all mechanical qualification tests are derived from mechanical requirements by adding explicit traceability between them.
- Shall define and implement the mechanical qualification tests against the mechanical requirements to validate that the mechanical architecture and mechanical implementation correctly realize them.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall ask for approval on key mechanical qualification test decisions.
- Shall set produced/modified work products status to `not approved`.

## Resources

- [development_methodology_overview.md](../../../doc/development_methodology/development_methodology_overview.md)
- [mechanical_domain_overview.md](../../../doc/development_methodology/mechanical_domain/mechanical_domain_overview.md)
- [mech-qualification-test.md](../../../doc/development_methodology/mechanical_domain/processes/mech-qualification-test.md)
