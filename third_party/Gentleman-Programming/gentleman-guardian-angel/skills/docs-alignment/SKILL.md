# Skill: gga-docs-alignment

## Purpose
Keep documentation in sync with every code and workflow change in GGA.

## When to Use
Any code or workflow change that affects user or contributor behavior.

## Docs Map

| File | Covers |
|------|--------|
| `README.md` | Installation, usage, CLI flags, providers, examples |
| `CONTRIBUTING.md` | Dev setup, contribution workflow, labels, testing, PR process |

## Change → Doc Rules

| If you change... | Update... |
|-----------------|-----------|
| A CLI flag (add, remove, rename) | `README.md` — flags table + examples |
| A provider (add, remove, change behavior) | `README.md` + `CONTRIBUTING.md` |
| The contribution workflow | `CONTRIBUTING.md` |
| The label system | `CONTRIBUTING.md` |
| Test structure or test commands | `CONTRIBUTING.md` |
| Install/uninstall behavior | `README.md` |
| The `.gga` config format | `README.md` |
| Hook behavior or marker format | `README.md` + `CONTRIBUTING.md` |
| CI mode behavior (`--ci`) | `README.md` |

## Critical Rules
- NEVER reference CLI flags, commands, or scripts that don't exist in code
- NEVER let a PR that changes behavior ship without updating docs
- Docs must describe the FINAL state, not the journey ("now supports X" is wrong)
- Code examples in docs must be tested and accurate
- PRs that change behavior without updating docs should be rejected in review

## Verification Checklist
Before submitting a PR with code changes:

```
[ ] Does this change affect any CLI flags?      → update README.md
[ ] Does this change affect provider behavior?  → update README.md + CONTRIBUTING.md
[ ] Does this change affect the contribution
    workflow or labels?                         → update CONTRIBUTING.md
[ ] Does this add/remove a config option?       → update README.md
[ ] Are all code examples in the docs still
    accurate after this change?                 → verify manually
```

## Anti-patterns

```
# BAD — README still shows old flag name
README: --cache-dir
Code: --cache-path  ← was renamed

# BAD — docs reference non-existent script
README: "run setup.sh to configure"  ← setup.sh doesn't exist

# BAD — missing provider in list
README: "Supports: Claude, Gemini, Codex"  ← Ollama was just added but not listed
```

## Cookbook

| If... | Then... |
|-------|---------|
| Adding a new provider | Add it to README providers table + CONTRIBUTING dev setup |
| Renaming a CLI flag | Search README for old flag name, update all occurrences |
| Changing PR workflow | Update CONTRIBUTING.md contribution steps section |
| Adding a new make target | Add to CONTRIBUTING.md development commands |
