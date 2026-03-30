---
name: sw-detailed-design
description: >
  Define software detailed design from software requirements and software architecture.
  Trigger: When software requirements and software architecture are approved and software detailed design shall be defined.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- After software requirements and software architecture are approved.
- When any software detailed design activity is required.

## Critical Patterns

- Shall follow the provided resources.
- Shall align with software requirements and software architecture (no inconsistencies).
- Shall ensure all software detailed design is derived from software architecture by adding explicit traceability between them.
- Shall prioritize reuse over creation.
- Shall avoid embedding implementation details beyond what is required for software detailed design.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall identify and list all opportunities for reuse and ask for user decision before classifying them as reusable or project-specific.
- Shall ask for approval on key software detailed design decisions.
- Shall request final validation before closing the software detailed design to ensure no further iterations are needed.

## Resources

- **Documentation**: See [references/docs.md](references/docs.md) for resources needed to perform this skill.