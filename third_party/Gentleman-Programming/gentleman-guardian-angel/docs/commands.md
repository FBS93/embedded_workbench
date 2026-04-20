# Commands

> 📖 Back to [README](../README.md)

Full command reference for Gentleman Guardian Angel.

---

## Commands Table

| Command                     | Description                                               | Example                         |
| --------------------------- | --------------------------------------------------------- | ------------------------------- |
| `init`                      | Create sample `.gga` config file                          | `gga init`                      |
| `install`                   | Install git pre-commit hook (default)                     | `gga install`                   |
| `install --commit-msg`      | Install git commit-msg hook (for commit message validation) | `gga install --commit-msg`    |
| `uninstall`                 | Remove git hooks from current repo                        | `gga uninstall`                 |
| `run`                       | Run code review on staged files                           | `gga run`                       |
| `run --ci`                  | Run code review on last commit (for CI/CD)                | `gga run --ci`                  |
| `run --pr-mode`             | Review all files changed in the full PR                   | `gga run --pr-mode`             |
| `run --pr-mode --diff-only` | PR review with diffs only (faster, cheaper)               | `gga run --pr-mode --diff-only` |
| `run --no-cache`            | Run review ignoring cache                                 | `gga run --no-cache`            |
| `config`                    | Display current configuration and status                  | `gga config`                    |
| `cache status`              | Show cache status for current project                     | `gga cache status`              |
| `cache clear`               | Clear cache for current project                           | `gga cache clear`               |
| `cache clear-all`           | Clear all cached data                                     | `gga cache clear-all`           |
| `help`                      | Show help message with all commands                       | `gga help`                      |
| `version`                   | Show installed version                                    | `gga version`                   |

---

## Command Details

### `gga init`

Creates a sample `.gga` configuration file in your project root with sensible defaults.

```bash
$ gga init
✅ Created config file: .gga
```

---

### `gga install`

Installs a git hook that automatically runs code review on every commit.

**Default (pre-commit hook):**

```bash
$ gga install
✅ Installed pre-commit hook: .git/hooks/pre-commit
```

**With commit message validation (commit-msg hook):**

```bash
$ gga install --commit-msg
✅ Installed commit-msg hook: .git/hooks/commit-msg
```

The `--commit-msg` flag installs a commit-msg hook instead of pre-commit. This allows GGA to also validate your commit message (e.g., conventional commits format, issue references, etc.). The commit message is automatically included in the AI review.

If a hook already exists, GGA will append to it rather than replacing it.

---

### `gga uninstall`

Removes the git pre-commit hook from your repository.

```bash
$ gga uninstall
✅ Removed pre-commit hook
```

---

### `gga run [--no-cache]`

Runs code review on currently staged files. Uses intelligent caching by default to skip unchanged files.

```bash
$ git add src/components/Button.tsx
$ gga run
# Reviews the staged file (uses cache)

$ gga run --no-cache
# Forces review of all files, ignoring cache
```

---

### `gga config`

Shows the current configuration, including where config files are loaded from and all settings.

```bash
$ gga config

Current Configuration:

Config Files:
  Global:  Not found
  Project: .gga

Values:
  PROVIDER:          claude
  FILE_PATTERNS:     *.ts,*.tsx,*.js,*.jsx
  EXCLUDE_PATTERNS:  *.test.ts,*.spec.ts
  RULES_FILE:        AGENTS.md
  STRICT_MODE:       true
  TIMEOUT:           300s
  PR_BASE_BRANCH:    auto-detect

Rules File: Found
```

---

## 🚫 Bypass Review

Sometimes you need to commit without review:

```bash
# Skip pre-commit hook entirely
git commit --no-verify -m "wip: work in progress"

# Short form
git commit -n -m "hotfix: urgent fix"
```
