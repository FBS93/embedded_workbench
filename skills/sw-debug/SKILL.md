---
name: sw-debug
description: >
  Debug software issues by investigating problems and deriving accurate conclusions based on evidence to enable correct resolution.
  Trigger: When a software defect or unexpected behavior needs to be investigated and understood.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- When a software defect needs to be investigated or reproduced.
- When software behavior is incorrect or not as expected.

## Critical Patterns

- Shall follow the provided resources.
- Shall use deterministic and minimal reproduction cases to isolate defects.
- Shall prioritize debugging on host for fast iteration and reproducibility.
- Shall escalate to target debugging only when the issue depends on hardware behavior, timing, or cannot be reproduced on host.
- Shall ensure any changes introduced for debugging purposes are controlled and removed once the debugging session is completed.
- Shall base all conclusions on verified and reproducible evidence to ensure correct problem resolution.
- Shall avoid making assumptions and escalate when the root cause cannot be determined with confidence.

## Resources

- **Documentation**: See [references/docs.md](references/docs.md) for resources needed to perform this skill.