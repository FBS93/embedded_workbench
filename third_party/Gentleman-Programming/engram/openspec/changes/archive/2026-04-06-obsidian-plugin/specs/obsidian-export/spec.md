# obsidian-export Specification

## Purpose

CLI command that reads the Engram SQLite store and writes a structured Obsidian-compatible markdown vault with YAML frontmatter, wikilinks, hub notes, and incremental sync state.

---

## Requirements

### Requirement: REQ-EXPORT-01 — Command Parsing

The `engram obsidian-export` subcommand MUST accept the following flags:
- `--vault <path>` (required): absolute or relative path to the Obsidian vault root
- `--project <name>` (optional): filter export to a single project
- `--limit <n>` (optional): cap exported observations at n (default: unlimited)
- `--since <date>` (optional): export only observations created after this ISO-8601 date

Missing `--vault` MUST produce a descriptive error and exit code 1.

#### Scenario: Happy path — all flags provided

- GIVEN the user runs `engram obsidian-export --vault ~/vault --project eng --limit 100 --since 2026-01-01`
- WHEN the command is parsed
- THEN export runs scoped to project "eng", max 100 observations, created after 2026-01-01

#### Scenario: Missing required --vault flag

- GIVEN the user runs `engram obsidian-export --project eng`
- WHEN the command is parsed
- THEN an error message "flag --vault is required" is printed to stderr
- AND the process exits with code 1

#### Scenario: Unknown flag

- GIVEN the user runs `engram obsidian-export --unknown`
- WHEN the command is parsed
- THEN an error is returned and the process exits with code 1

---

### Requirement: REQ-EXPORT-02 — Vault Directory Structure

The exporter MUST create the following directory hierarchy inside `{vault}/engram/`:

```
{vault}/engram/
├── {project}/{type}/{slug}-{id}.md   ← observation notes
├── _sessions/{session-id}.md          ← session hub notes
├── _topics/{topic-prefix}.md          ← topic cluster hub notes
└── .engram-sync-state.json            ← incremental sync state
```

All directories MUST be created if they do not exist. The exporter MUST NOT write outside the `{vault}/engram/` subdirectory.

#### Scenario: Fresh vault — directories don't exist

- GIVEN `--vault ~/vault` points to a directory without an `engram/` subfolder
- WHEN export runs
- THEN `~/vault/engram/` and all needed subdirectories are created
- AND observation files are written inside them

#### Scenario: Vault directory not writable

- GIVEN the vault path exists but the user has no write permission
- WHEN export runs
- THEN an error is printed to stderr and the process exits with code 1

---

### Requirement: REQ-EXPORT-03 — Observation → Markdown Conversion

Each observation MUST be written as a markdown file with:
- A YAML frontmatter block containing: `type`, `project`, `scope`, `topic_key`, `session_id`, `created_at`, `tags`
- The observation `content` field as the markdown body
- A `## Wikilinks` section listing all generated `[[links]]`

`topic_key` and `session_id` in frontmatter MAY be empty strings when the observation has no value.

#### Scenario: Observation with all fields populated

- GIVEN an observation with type="bugfix", project="eng", scope="project", topic_key="auth/jwt", session_id="abc123", created_at="2026-01-01T10:00:00Z", content="Fixed the bug"
- WHEN the markdown file is generated
- THEN the file starts with a YAML frontmatter block containing all six fields
- AND the body contains "Fixed the bug"
- AND a `## Wikilinks` section lists `[[session-abc123]]` and `[[topic-auth]]`

#### Scenario: Observation with no topic_key

- GIVEN an observation with no topic_key set
- WHEN the markdown file is generated
- THEN `topic_key: ""` appears in frontmatter
- AND no `[[topic-*]]` wikilink is emitted

---

### Requirement: REQ-EXPORT-04 — Session Hub Notes

For each unique `session_id` among exported observations, the exporter MUST generate a hub note at `_sessions/{session-id}.md` listing all observations in that session as wikilinks.

#### Scenario: Session with multiple observations

- GIVEN observations obs-1 and obs-2 both have session_id="sess-42"
- WHEN export runs
- THEN `_sessions/sess-42.md` is created
- AND it contains `[[obs-1-slug-1]]` and `[[obs-2-slug-2]]` as wikilinks

#### Scenario: Observation with no session_id

- GIVEN an observation with no session_id
- WHEN export runs
- THEN no session hub note is generated for that observation
- AND the observation note is still exported normally

---

### Requirement: REQ-EXPORT-05 — Topic Cluster Hub Notes

A topic hub note at `_topics/{prefix}.md` MUST be generated **only when ≥2 exported observations share the same topic_key prefix** (the first segment of a `/`-delimited topic_key).

Observations with an empty or absent topic_key MUST be excluded from prefix grouping.

#### Scenario: Two observations share a prefix → hub created

- GIVEN obs-A has topic_key="sdd/spec" and obs-B has topic_key="sdd/design"
- WHEN export runs
- THEN `_topics/sdd.md` is created listing both observations as wikilinks

#### Scenario: Only one observation with a given prefix → hub NOT created

- GIVEN only obs-A has topic_key="auth/jwt" and no other observation shares prefix "auth"
- WHEN export runs
- THEN `_topics/auth.md` is NOT created

#### Scenario: Observation with no topic_key → excluded from hub grouping

- GIVEN obs-C has no topic_key
- WHEN export runs
- THEN obs-C does not contribute to any topic hub prefix count

---

### Requirement: REQ-EXPORT-06 — Incremental Sync

The exporter MUST persist the last-export timestamp in `{vault}/engram/.engram-sync-state.json`. On subsequent runs, it MUST only write observations created or updated after that timestamp.

#### Scenario: First export — no state file

- GIVEN `.engram-sync-state.json` does not exist
- WHEN export runs
- THEN all matching observations are exported
- AND `.engram-sync-state.json` is created with `last_export_at` set to the current time

#### Scenario: Incremental export — state file exists

- GIVEN `.engram-sync-state.json` has `last_export_at: "2026-03-01T00:00:00Z"`
- AND a new observation was created at "2026-04-01T00:00:00Z"
- WHEN export runs without `--since`
- THEN only the new observation is written
- AND `.engram-sync-state.json` is updated to the new timestamp

#### Scenario: --since flag overrides state file

- GIVEN `.engram-sync-state.json` has `last_export_at: "2026-03-01T00:00:00Z"`
- AND the user passes `--since 2026-01-01`
- WHEN export runs
- THEN observations since 2026-01-01 are exported (overriding the state file timestamp)

---

### Requirement: REQ-EXPORT-07 — Collision-Safe File Naming

Observation files MUST be named `{slug}-{id}.md` where `{slug}` is a lowercase, hyphen-separated truncation of the observation title or first 40 chars of content, and `{id}` is the numeric observation ID.

#### Scenario: Two observations with identical content previews

- GIVEN obs-1 (id=1) and obs-2 (id=2) both start with "Fixed authentication"
- WHEN filenames are generated
- THEN obs-1 → `fixed-authentication-1.md` and obs-2 → `fixed-authentication-2.md`

---

### Requirement: REQ-EXPORT-08 — Idempotency

Re-running export MUST NOT create duplicate files. If a file already exists and the observation has not changed, the file MUST be left unchanged. If the observation content has changed, the file MUST be overwritten.

#### Scenario: Re-export with no changes

- GIVEN a vault with previously exported files
- AND no observations have changed
- WHEN export runs again
- THEN no files are written (no disk I/O for unchanged observations)

#### Scenario: Re-export after observation updated

- GIVEN obs-5 was exported previously
- AND obs-5's content has since been updated in the store
- WHEN export runs again
- THEN obs-5's markdown file is overwritten with the new content

---

### Requirement: REQ-EXPORT-09 — Deleted Observations Handling

Observations that are soft-deleted in the store (deleted_at IS NOT NULL) MUST NOT be exported. If a previously exported file corresponds to a now-deleted observation, it MUST be removed from the vault.

#### Scenario: Observation deleted after first export

- GIVEN obs-3 was exported to `eng/bugfix/some-fix-3.md`
- AND obs-3 is later soft-deleted in the store
- WHEN export runs again
- THEN `eng/bugfix/some-fix-3.md` is deleted from the vault
