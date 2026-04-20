<!--
  IMPORTANT: PRs without a linked issue will be CLOSED automatically.
  This project follows a strict issue-first workflow:

  1. Open an issue describing the bug or feature
  2. Wait for a maintainer to add `status:approved`
  3. THEN open this PR linking that issue

  Skipping any of these steps = PR closed without review.
-->

## Linked Issue

Closes #<!-- issue number here -->

## PR Type

<!-- Check the one that applies. PRs without a type:* label will fail the automated check. -->

- [ ] `type:bug` — Bug fix
- [ ] `type:feature` — New feature or enhancement
- [ ] `type:docs` — Documentation changes only
- [ ] `type:refactor` — Code refactor (no behavior change)
- [ ] `type:chore` — Maintenance (deps, CI, tooling)
- [ ] `type:breaking-change` — Breaking change (add `!` to commit type)

## Summary

<!-- 2-4 sentences: what does this PR do, and why? -->

## Changes

| File | Change |
|------|--------|
| `bin/gga` | |
| `lib/` | |
| `spec/` | |

## Test Plan

<!-- Check all that apply and describe what you tested -->

- [ ] `make lint` (ShellCheck) passes locally
- [ ] `make test` passes locally (all unit tests)
- [ ] `shellspec spec/unit` passes
- [ ] `shellspec spec/integration/commands_spec.sh` passes
- [ ] Manual testing: <!-- describe what you tested manually -->

## Automated Checks

The following checks run automatically on every PR:

| Check | Validates |
|-------|-----------|
| Check Issue Reference | PR body contains `Closes/Fixes/Resolves #N` |
| Check Issue Has `status:approved` | Linked issue was approved by a maintainer |
| Check PR Has `type:*` Label | PR has exactly one `type:*` label |
| Lint | ShellCheck passes on `bin/gga` and `lib/*.sh` |
| Unit Tests | `shellspec spec/unit` passes |
| Integration Tests | `shellspec spec/integration/commands_spec.sh` passes |

## Contributor Checklist

- [ ] I have linked the issue above with `Closes #N`
- [ ] The linked issue has `status:approved`
- [ ] I have added a `type:*` label to this PR
- [ ] `make lint` passes (ShellCheck)
- [ ] `make test` passes (all unit tests)
- [ ] Integration tests pass
- [ ] New functions/commands have tests in `spec/`
- [ ] I followed conventional commits (`feat:`, `fix:`, `docs:`, etc.)
- [ ] I have NOT added "Co-Authored-By" or AI attribution to commits
- [ ] Breaking changes are documented and use `!` in the commit type (`feat!:`, `fix!:`)

## Notes for Reviewers

<!-- Optional: anything reviewers should pay special attention to, edge cases, or known limitations -->
