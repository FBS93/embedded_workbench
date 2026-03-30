---
name: sw-integration-test
description: >
  Define and implement software integration tests against software architecture to validate that the software detailed design and software implementation correctly realizes it.
  Trigger: When software unit test is completed and software integration tests shall be defined and implemented.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- After software unit test is completed.
- When any software integration test activity is required.

## Critical Patterns

- Shall follow the provided resources.
- Shall align with software architecture (no inconsistencies).
- Shall ensure all software integration tests are derived from software architecture by adding explicit traceability between them.
- Shall define and implement the software integration tests against the software architecture to validate that the software detailed design and software implementation correctly realizes it.
- Shall keep documentation DRY (Don't Repeat Yourself) by adding references instead of duplicating content.
- Shall ask for approval on key software integration test decisions.
- Shall request final validation before closing the software integration test to ensure no further iterations are needed.

## Resources

- **Documentation**: See [references/docs.md](references/docs.md) for resources needed to perform this skill.