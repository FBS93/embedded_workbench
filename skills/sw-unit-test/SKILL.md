---
name: sw-unit-test
description: >
  Define and implement software unit tests against software detailed design to validate that the software implementation correctly realizes it.
  Trigger: When software implementation is completed and software unit tests shall be defined and implemented.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- After software implementation is completed.
- When any software unit test activity is required.

## Critical Patterns

- Shall follow the provided resources.
- Shall align with software detailed design (no inconsistencies).
- Shall ensure all software unit tests are derived from software detailed design by adding explicit traceability between them.
- Shall define and implement the software unit tests against the software detailed design to validate that the software implementation correctly realizes it.
- Shall never define and implement the software unit tests against the software implementation.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall ask for approval on key software unit test decisions.
- Shall request final validation before closing the software unit test to ensure no further iterations are needed.

## Resources

- **Documentation**: See [references/docs.md](references/docs.md) for resources needed to perform this skill.