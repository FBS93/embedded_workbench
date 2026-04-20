<div align="center">

# Agent Teams Lite

**Give any AI coding agent a team of specialized sub-agents.**

One install. Works with Claude Code, Gemini CLI, OpenCode, Cursor, Windsurf, Codex, and VS Code Copilot.

[Quick Start](#quick-start) · [How It Works](#how-it-works) · [Commands](#commands) · [Supported Tools](#supported-tools) · [Docs](#learn-more) · [Contributing](CONTRIBUTING.md)

</div>

---

## The Problem

AI coding assistants are powerful, but they struggle with complex features:

- **Context overload** — Long conversations lead to compression, lost details, hallucinations
- **No structure** — "Build me dark mode" produces unpredictable results
- **No review gate** — Code gets written before anyone agrees on what to build
- **No memory** — Specs live in chat history that vanishes

## The Solution

**Agent Teams Lite** is an agent-team orchestration pattern where a lightweight coordinator delegates all real work to specialized sub-agents. Each sub-agent starts with fresh context, executes one focused task, and returns a structured result.

```
YOU: "I want to add CSV export to the app"

ORCHESTRATOR (delegate-only, minimal context):
  → launches EXPLORER sub-agent     → returns: codebase analysis
  → shows you summary, you approve
  → launches PROPOSER sub-agent     → returns: proposal artifact
  → launches SPEC WRITER sub-agent  → returns: spec artifact
  → launches DESIGNER sub-agent     → returns: design artifact
  → launches TASK PLANNER sub-agent → returns: tasks artifact
  → shows you everything, you approve
  → launches IMPLEMENTER sub-agent  → returns: code written, tasks checked off
  → launches VERIFIER sub-agent     → returns: verification artifact
  → launches ARCHIVER sub-agent     → returns: change closed
```

**The key insight**: the orchestrator NEVER does real work directly. It delegates everything to sub-agents, tracks state, and synthesizes summaries. This keeps the main thread small and stable.

Supports multiple persistence modes — see [persistence docs](docs/persistence.md).

## How It Works

Three concepts:

1. **Delegate-first architecture** — Your main agent becomes an orchestrator that delegates all real work to specialized sub-agents. Each sub-agent gets a fresh context, does focused work, and returns only a summary. [Learn more →](docs/architecture.md)

2. **Spec-Driven Development (SDD)** — A DAG of phases: `explore → propose → spec + design → tasks → apply → verify → archive`. Each phase is a skill that any AI agent can run. [See the phases →](docs/sub-agents.md)

3. **Your skills, pre-loaded** — A skill registry catalogs your coding standards (React, TDD, Tailwind, etc.), generates compact rules, and the orchestrator injects the relevant ones into every sub-agent automatically. [Details →](docs/sub-agents.md#skill-registry)

## Quick Start

```bash
git clone https://github.com/Gentleman-Programming/agent-teams-lite.git
cd agent-teams-lite
./scripts/setup.sh --all
```

That's it. The script detects which tools you have installed, installs the full 15-skill set, and configures them.

> **Windows?** Run `.\scripts\setup.ps1 -All` in PowerShell.

For manual installation or per-tool setup, see [docs/installation.md](docs/installation.md).

## Commands

| Command | What It Does |
|---------|-------------|
| `/sdd-init` | Initialize SDD context. Detects stack, bootstraps persistence, builds skill registry. |
| `/sdd-explore <topic>` | Investigate an idea. Reads codebase, compares approaches. No files created. |
| `/sdd-new <name>` | Start a new change by delegating exploration + proposal to sub-agents. |
| `/sdd-continue` | Run the next dependency-ready phase via sub-agent(s). |
| `/sdd-ff <name>` | Fast-forward planning with sub-agents (proposal → specs → design → tasks). |
| `/sdd-apply` | Implement tasks in batches. Checks off items as it goes. |
| `/sdd-verify` | Validate implementation against specs. Reports CRITICAL / WARNING / SUGGESTION. |
| `/sdd-archive` | Close a change and persist final state in the active artifact store. |
| `/skill-registry` | Create or update the skill registry for the current project. |

## Supported Tools

| Tool | Sub-agent support | Setup |
|------|------------------|-------|
| Claude Code | Full (Task tool) | `./scripts/setup.sh --claude` |
| OpenCode | Full (delegate/task) | `./scripts/setup.sh --opencode` |
| Gemini CLI | Full (Task tool) | `./scripts/setup.sh --gemini` |
| Codex | Full (Task tool) | `./scripts/setup.sh --codex` |
| Cursor | Inline only | `./scripts/setup.sh --cursor` |
| VS Code Copilot | Inline only | `./scripts/setup.sh --vscode` |
| Windsurf | Inline only | `./scripts/setup.sh --windsurf` |

> **Full** = orchestrator delegates to sub-agents. **Inline** = skills loaded directly into main agent.

## Why Not Just Use OpenSpec?

| | OpenSpec | Agent Teams Lite |
|---|---|---|
| **Dependencies** | Requires `npm install -g @fission-ai/openspec` | Zero. Pure Markdown files. |
| **Sub-agents** | Runs inline (one context window) | True sub-agent delegation (fresh context per phase) |
| **Context usage** | Everything in one conversation | Orchestrator stays lightweight, sub-agents get fresh context |
| **Customization** | Edit YAML schemas + rebuild | Edit Markdown files, instant effect |
| **Tool support** | 20+ tools via CLI | Any tool that can read Markdown (infinite) |
| **Setup** | CLI init + slash commands | Copy files + go |

## Token Economics

We measured the real token costs of delegation across 3 independent analyses. Delegation saves 50-70% of tokens for medium-to-large features by avoiding context pollution and compaction cascades.

→ **[Full analysis: docs/token-economics.md](docs/token-economics.md)**

## Learn More

| Topic | Description |
|-------|-------------|
| [Architecture](docs/architecture.md) | System diagrams, DAG, capability comparison, project structure |
| [Installation](docs/installation.md) | Per-tool setup, manual install, verification steps |
| [Sub-Agents](docs/sub-agents.md) | Phase descriptions, skill registry, shared conventions |
| [Persistence](docs/persistence.md) | Storage modes (engram, openspec, hybrid, none) |
| [Concepts](docs/concepts.md) | Delta specs, RFC 2119 keywords, archive cycle |
| [Token Economics](docs/token-economics.md) | Real-world token usage analysis and optimizations |
| [Changelog](docs/changelog.md) | Version history and notable upgrades |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full workflow. TL;DR:

1. Open an issue using a template (bug or feature)
2. Wait for `status:approved` label
3. Open a PR linking the approved issue

## License

Apache 2.0 — see [LICENSE](LICENSE).

---

<div align="center">
  <sub>Built by <a href="https://github.com/Gentleman-Programming">Gentleman Programming</a></sub>
</div>
