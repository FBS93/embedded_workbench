# Contributing to Gentleman Guardian Angel (GGA)

This project follows a **strict issue-first workflow**. No PR is accepted without a linked, approved issue. Read this document before opening any issue or PR.

---

## Contribution Workflow

```
Open Issue → Get status:approved → Open PR → Add type:* label → Review & Merge
```

Every contribution follows these exact steps — no exceptions.

---

## Step 1: Open an Issue

Use the appropriate issue template:

- **[Bug Report](https://github.com/Gentleman-Programming/gentleman-guardian-angel/issues/new?template=bug_report.yml)** — Something is broken
- **[Feature Request](https://github.com/Gentleman-Programming/gentleman-guardian-angel/issues/new?template=feature_request.yml)** — New feature or improvement

> Blank issues are disabled. You must use a template.

Describe the problem or proposal clearly. The more context you provide, the faster the review.

---

## Step 2: Wait for Approval

A maintainer will review your issue and add one of these labels:

| Label | Meaning |
|-------|---------|
| `status:needs-review` | Added automatically when you open the issue |
| `status:approved` | Maintainer approved — you may now open a PR |

**Do not open a PR until your issue has `status:approved`.** PRs without an approved issue will be closed automatically.

---

## Step 3: Open a Pull Request

Once your issue has `status:approved`:

1. **Fork** the repository
2. **Create a branch** from `main` with a descriptive name:
   ```bash
   git checkout -b feat/add-github-models-provider
   git checkout -b fix/hooks-marker-injection
   ```
3. **Make your changes** — see [Development Setup](#development-setup) below
4. **Open a PR** on GitHub
5. **Link your issue** in the PR body:
   ```
   Closes #42
   ```
6. **Add a `type:*` label** to your PR (see [Label System](#label-system))

---

## Step 4: Automated PR Checks

Two sets of checks run automatically on every PR:

### PR Validation (`.github/workflows/pr-check.yml`)

| Check | What It Validates |
|-------|------------------|
| Check Issue Reference | PR body contains `Closes/Fixes/Resolves #N` |
| Check Issue Has `status:approved` | Linked issue was approved by a maintainer |
| Check PR Has `type:*` Label | PR has exactly one `type:*` label |

### CI Tests (`.github/workflows/ci.yml`)

| Check | Command |
|-------|---------|
| Lint | ShellCheck on `bin/gga` and `lib/*.sh` |
| Unit Tests | `shellspec spec/unit` |
| Integration Tests | `shellspec spec/integration/commands_spec.sh` |

**All checks must pass before a PR can be merged.**

---

## Label System

### Type Labels

Applied to PRs to categorize the change:

| Label | When to Use |
|-------|-------------|
| `type:bug` | Fixes a bug |
| `type:feature` | Adds a new feature or enhancement |
| `type:docs` | Documentation changes only |
| `type:refactor` | Code refactor with no behavior change |
| `type:chore` | Maintenance: deps, CI, tooling |
| `type:breaking-change` | Breaking change (also use `!` in commit type) |

### Status Labels

Applied to issues by maintainers:

| Label | Meaning |
|-------|---------|
| `status:needs-review` | Waiting for maintainer review |
| `status:approved` | Approved — a PR may be opened |

### Priority Labels

| Label | Meaning |
|-------|---------|
| `priority:high` | Blocking or critical |
| `priority:medium` | Important but not urgent |
| `priority:low` | Nice to have |

---

## PR Rules

- **Focused scope** — one issue per PR, no unrelated changes
- **Conventional commits** — see [Conventional Commit Format](#conventional-commit-format) below
- **Run tests locally** before pushing — do not open a PR that breaks CI
- **Update documentation** if your change affects behavior or adds commands
- **No Co-Authored-By** — do not add AI attribution to commits

---

## Conventional Commit Format

```
<type>(<scope>): <short description>

[optional body]

Closes #<issue>
```

### Examples

```bash
feat(providers): add GitHub Models provider support
fix(hooks): resolve marker injection in existing hooks
docs(contributing): add label system documentation
refactor(cache): extract hash computation to lib function
chore(deps): bump ShellSpec version in CI
fix!: change hook marker format (breaking change)
```

Scopes: `providers`, `hooks`, `cache`, `cli`, `ci`, `config`, `docs`

For breaking changes, add `!` after the type: `feat!:`, `fix!:`

---

## Development Setup

```bash
# Clone the repo
git clone git@github.com:Gentleman-Programming/gentleman-guardian-angel.git
cd gentleman-guardian-angel

# Install ShellSpec (macOS)
brew install shellspec

# Install ShellSpec (Linux)
curl -fsSL https://git.io/shellspec | sh -s -- --yes

# Run all tests
make test

# Run linter
make lint
```

---

## Testing

```bash
# Run all tests
make test

# Lint only (ShellCheck)
make lint

# Unit tests only
shellspec spec/unit

# Integration tests
shellspec spec/integration/commands_spec.sh

# Specific test file
shellspec spec/unit/cache_spec.sh

# Specific test by line number
shellspec spec/unit/cache_spec.sh:42

# Verbose output
shellspec --format documentation
```

Tests live in:
- `spec/unit/` — fast tests with mocks, no external dependencies
- `spec/integration/` — tests requiring a real git repo or environment

All new features and bug fixes must include tests.
