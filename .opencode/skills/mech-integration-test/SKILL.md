---
name: mech-integration-test
description: >
  Define and implement mechanical integration tests against mechanical architecture to validate that the mechanical implementation correctly realizes it.
  Trigger: When mechanical implementation is completed and mechanical integration tests shall be defined and implemented.
---

## When to Use

- After mechanical implementation is completed.
- When any mechanical integration test activity is required.

## Critical Patterns

- Shall strictly follow the provided resources and recursively review all referenced documents required to understand the current process, inputs, outputs, rules and directly referenced work products.
- Shall align with mechanical architecture (no inconsistencies).
- Shall ensure all mechanical integration tests are derived from mechanical architecture by adding explicit traceability between them.
- Shall define and implement the mechanical integration tests against the mechanical architecture to validate that the mechanical implementation correctly realizes it.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall ask for approval on key mechanical integration test decisions.
- Shall set produced/modified work products status to `not approved`.

## Resources

- [development_methodology_overview.md](../../../doc/development_methodology/development_methodology_overview.md)
- [mechanical_domain_overview.md](../../../doc/development_methodology/mechanical_domain/mechanical_domain_overview.md)
- [mech-integration-test.md](../../../doc/development_methodology/mechanical_domain/processes/mech-integration-test.md)
