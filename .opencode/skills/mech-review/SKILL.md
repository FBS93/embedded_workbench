---
name: mech-review
description: >
  Review mechanical process work products before closure.
  Trigger: When a mechanical process is ready to be closed and requires validation.
---

## When to Use

- Before closing any mechanical process activity.
- When a full review of mechanical process work products is required.

## Critical Patterns

- Shall strictly follow the provided resources and recursively review all referenced documents required to understand the current task.
- Shall verify the following points only for mechanical process work products with status `not approved`:
  - Ensure that all mechanical process work products are complete and consistent.
  - Identify missing traceability, inconsistencies, or deviations from the defined mechanical process.
  - Verify that documentation is DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall update the status of compliant reviewed mechanical process work products to `approved`; work products that fail the review shall remain `not approved`.
- Shall clearly state whether the mechanical process work products are approved and the process can be closed, or whether further work is needed.

## Resources

- [development_methodology_overview.md](../../../doc/development_methodology/development_methodology_overview.md)
- [mechanical_domain_overview.md](../../../doc/development_methodology/mechanical_domain/mechanical_domain_overview.md)
