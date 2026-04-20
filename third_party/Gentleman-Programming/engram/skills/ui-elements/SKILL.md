---
name: engram-ui-elements
description: >
  Creation rules for Engram UI elements, pages, cards, metrics, and detail flows.
  Trigger: Adding or changing dashboard UI components or connected browsing flows.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Adding a new page or partial to the dashboard
- Creating cards, metrics, tables, lists, or detail views
- Designing connected navigation between related entities

---

## UX Rules

1. Every list item should lead somewhere useful when domain relationships exist.
2. Prefer connected flows: project -> session -> observation -> full detail.
3. Empty states must explain what is missing and what unlocks data.
4. Metrics must reflect real system state, not decorative counters.
5. Detail pages should show metadata, content, and the next relevant links.

---

## Composition Rules

- Use metrics for orientation, not for replacing core content.
- Use cards for browsable entities and tables for dense comparative admin data.
- Avoid nested framed boxes unless they communicate a hierarchy the user needs.
- Keep action controls close to the entity they affect.
