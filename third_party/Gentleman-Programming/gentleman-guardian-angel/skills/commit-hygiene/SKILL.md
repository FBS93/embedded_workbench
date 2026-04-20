# Skill: gga-commit-hygiene

## Purpose
Enforce conventional commits and a clean, readable history for GGA.

## When to Use
Any commit creation, review, or branch cleanup.

## Format
```
type(scope): description
```

Validated by regex:
```
^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9\._-]+\))?!?: .+
```

## Valid Types

| Type | Use For |
|------|---------|
| feat | New feature or behavior |
| fix | Bug fix |
| docs | Documentation only |
| chore | Maintenance, dependencies, tooling |
| refactor | Restructuring without behavior change |
| test | Adding or fixing tests |
| ci | CI/CD pipeline changes |
| style | Formatting, whitespace, no logic change |
| perf | Performance improvement |
| build | Build system, Makefile changes |
| revert | Reverting a previous commit |

## Valid Scopes for GGA

| Scope | Covers |
|-------|--------|
| providers | `lib/providers.sh`, any AI provider integration |
| hooks | Hook system, marker injection, pre-commit logic |
| cache | `lib/cache.sh`, hash logic |
| cli | `bin/gga` flags, argument parsing, output |
| ci | GitHub Actions workflows |
| config | `.gga` config file parsing |
| install | `install.sh`, `uninstall.sh` |

## Critical Rules
- NEVER add `Co-Authored-By` or any AI attribution trailers
- One logical change per commit — no "add X and fix Y and update docs"
- Use imperative mood: "add" not "added", "fix" not "fixed"
- Reference issue numbers when relevant: `feat(providers): add GitHub Models support (#12)`
- Keep subject line ≤ 72 characters
- If a commit needs explanation, add a body separated by a blank line

## Examples

```
feat(providers): add GitHub Models provider support
fix(hooks): resolve marker injection in existing hooks
fix(cache): handle hash collision on empty files
test(cache): add unit tests for hash invalidation
test(hooks): add integration tests for marker injection
refactor(cli): extract provider dispatch to lib/providers.sh
ci: split unit and integration test jobs
docs(contributing): add label system documentation
chore: update shellspec to 0.28.1
```

## Anti-patterns

```
# BAD — vague
update stuff

# BAD — AI attribution
feat: add feature

Co-Authored-By: Claude <claude@anthropic.com>

# BAD — multiple concerns
feat(providers): add ollama and fix hook bug and update readme

# BAD — past tense
fix(cache): fixed the hash collision
```
