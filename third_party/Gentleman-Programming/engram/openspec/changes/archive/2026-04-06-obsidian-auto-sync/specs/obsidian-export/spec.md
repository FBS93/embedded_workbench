# Delta for obsidian-export
## Change: obsidian-auto-sync

> Extends REQ-EXPORT-01..09 from `openspec/changes/obsidian-plugin/specs/obsidian-export/spec.md`.
> Only ADDED requirements are listed here — existing requirements are unchanged.

---

## ADDED Requirements

### Requirement: REQ-GRAPH-01 — `--graph-config` Flag

The `engram obsidian-export` command MUST accept a `--graph-config` flag with accepted values: `preserve` (default), `force`, `skip`.
Any other value MUST cause the process to print a descriptive error to stderr and exit with code 1.

#### Scenario: Valid flag value accepted

- GIVEN the user runs `engram obsidian-export --vault ./v --graph-config=force`
- WHEN the flag is parsed
- THEN export proceeds with force mode

#### Scenario: Invalid flag value rejected

- GIVEN the user runs `engram obsidian-export --vault ./v --graph-config=invalid`
- WHEN the flag is parsed
- THEN an error "invalid --graph-config value: invalid (accepted: preserve, force, skip)" is printed to stderr
- AND the process exits with code 1

---

### Requirement: REQ-GRAPH-02 — `preserve` Mode

When `--graph-config=preserve` (default):
- If `{vault}/.obsidian/graph.json` does NOT exist, the system MUST write the embedded default template.
- If `{vault}/.obsidian/graph.json` DOES exist, the system MUST leave it unchanged and log "graph.json exists, preserving user config".

#### Scenario: First run — graph.json absent

- GIVEN `{vault}/.obsidian/graph.json` does not exist
- WHEN export runs with `--graph-config=preserve` (or default)
- THEN `{vault}/.obsidian/graph.json` is created with the embedded default
- AND no warning is logged

#### Scenario: Subsequent run — graph.json present

- GIVEN `{vault}/.obsidian/graph.json` already exists with custom content
- WHEN export runs with `--graph-config=preserve`
- THEN the file is unchanged
- AND "graph.json exists, preserving user config" is logged

---

### Requirement: REQ-GRAPH-03 — `force` Mode

When `--graph-config=force`, the system MUST always overwrite `{vault}/.obsidian/graph.json` with the embedded default and log "graph.json overwritten with engram default".

#### Scenario: Force overwrites existing file

- GIVEN `{vault}/.obsidian/graph.json` exists with custom content
- WHEN export runs with `--graph-config=force`
- THEN the file is overwritten with the embedded default
- AND "graph.json overwritten with engram default" is logged

---

### Requirement: REQ-GRAPH-04 — `skip` Mode

When `--graph-config=skip`, the system MUST NOT read, write, or create `{vault}/.obsidian/graph.json`. No log is emitted.

#### Scenario: Skip never touches graph.json

- GIVEN `{vault}/.obsidian/graph.json` does not exist
- WHEN export runs with `--graph-config=skip`
- THEN `{vault}/.obsidian/graph.json` is NOT created

---

### Requirement: REQ-GRAPH-05 — Embedded Default Template

The embedded default `graph.json` MUST be stored as a file in `internal/obsidian/` and included via `//go:embed`.
It MUST contain EXACTLY the following values (no deviation):

| Key | Value |
|-----|-------|
| `collapse-filter` | `false` |
| `collapse-color-groups` | `true` |
| `collapse-display` | `true` |
| `collapse-forces` | `false` |
| `showArrow` | `false` |
| `textFadeMultiplier` | `0` |
| `nodeSizeMultiplier` | `1` |
| `lineSizeMultiplier` | `1` |
| `centerStrength` | `0.515147569444444` |
| `repelStrength` | `12.7118055555556` |
| `linkStrength` | `0.729210069444444` |
| `linkDistance` | `207` |
| `scale` | `0.1` |
| `showOrphans` | `true` |

Color groups (MUST be EXACTLY these 6, in this order):

| # | query | rgb |
|---|-------|-----|
| 1 | `path:engram/_sessions` | 14736466 |
| 2 | `path:engram/_topics` | 13893887 |
| 3 | `tag:#architecture` | 7935 |
| 4 | `tag:#bugfix` | 16711680 |
| 5 | `tag:#decision` | 65322 |
| 6 | `tag:#pattern` | 16741120 |

All color group entries MUST have `"a": 1`. No discovery group. No additional groups.

#### Scenario: Embedded default validates expected structure

- GIVEN the embedded `graph.json` is parsed
- WHEN a unit test validates its contents
- THEN all keys and values from the table above match exactly
- AND `colorGroups` has exactly 6 entries with the queries and rgb values listed

---

### Requirement: REQ-GRAPH-06 — `.obsidian/` Directory Creation

If `{vault}/.obsidian/` does not exist when graph config needs to be written (modes `preserve` or `force`), the system MUST create it before writing `graph.json`.
This MUST NOT error if the directory already exists.

#### Scenario: .obsidian dir missing — created automatically

- GIVEN `{vault}/.obsidian/` does not exist
- AND `--graph-config=preserve` or `--graph-config=force`
- WHEN export runs
- THEN `{vault}/.obsidian/` is created
- AND `graph.json` is written inside it

#### Scenario: .obsidian dir already exists — no error

- GIVEN `{vault}/.obsidian/` already exists
- WHEN export runs
- THEN no error is produced related to the directory

---

### Requirement: REQ-WATCH-01 — `--watch` Flag (Daemon Mode)

The `engram obsidian-export` command MUST accept a `--watch` flag.
Without `--watch`: command runs once and exits (unchanged behavior).
With `--watch`: command loops indefinitely until SIGINT or SIGTERM is received.

#### Scenario: Without --watch — single run

- GIVEN the user runs `engram obsidian-export --vault ./v`
- WHEN export completes
- THEN the process exits with code 0

#### Scenario: With --watch — loops until signal

- GIVEN the user runs `engram obsidian-export --vault ./v --watch`
- WHEN the first export completes
- THEN the process does NOT exit; it waits for the next interval
- AND continues looping until SIGINT or SIGTERM

---

### Requirement: REQ-WATCH-02 — `--interval` Flag

The `--interval` flag MUST accept Go duration strings (e.g. `30s`, `5m`, `1h`).
It is only active when `--watch` is set. Default value: `10m`.
Values below `1m` MUST cause the process to exit with code 1 and a descriptive error.
Unparseable values MUST cause the process to exit with code 1.

#### Scenario: Valid interval accepted

- GIVEN the user runs `engram obsidian-export --vault ./v --watch --interval 5m`
- WHEN the flag is parsed
- THEN the daemon runs with a 5-minute interval

#### Scenario: Interval below minimum rejected

- GIVEN the user runs with `--watch --interval 30s`
- WHEN the flag is parsed
- THEN an error "--interval must be >= 1m; got 30s" is printed to stderr
- AND the process exits with code 1

#### Scenario: Unparseable interval rejected

- GIVEN the user runs with `--watch --interval banana`
- WHEN the flag is parsed
- THEN a parse error is printed to stderr
- AND the process exits with code 1

---

### Requirement: REQ-WATCH-03 — First-Run Semantics

When `--watch` is enabled, the FIRST export cycle MUST run IMMEDIATELY on startup without waiting for the interval.
Subsequent cycles run every `interval` after the previous cycle completes.

#### Scenario: First export is immediate

- GIVEN the user runs `engram obsidian-export --vault ./v --watch --interval 10m`
- WHEN the process starts
- THEN an export cycle begins immediately (within milliseconds, not after 10m)
- AND the next cycle starts approximately 10m after the first completes

---

### Requirement: REQ-WATCH-04 — Per-Cycle Logging

Each completed watch cycle MUST log a single line in this exact format:

```
[{RFC3339-timestamp}] sync: created={n} updated={n} deleted={n} skipped={n} hubs={n}
```

On cycle error: the error MUST be logged and the loop MUST continue. The process MUST NOT exit on transient failures.

#### Scenario: Successful cycle log line

- GIVEN a watch cycle completes with 5 created, 0 updated, 0 deleted, 1725 skipped, 277 hubs
- WHEN the cycle ends
- THEN a log line matching `[2026-04-06T21:00:00Z] sync: created=5 updated=0 deleted=0 skipped=1725 hubs=277` is printed (timestamp varies)

#### Scenario: Cycle error — loop continues

- GIVEN a watch cycle fails with a transient error (e.g. locked DB)
- WHEN the error occurs
- THEN the error is logged
- AND the next cycle runs at the next interval tick

---

### Requirement: REQ-WATCH-05 — Graceful Shutdown

On SIGINT (Ctrl+C) or SIGTERM, the system MUST:
1. Finish the current cycle if one is in progress.
2. Log "shutting down watch mode".
3. Exit with code 0.
4. NOT leave a corrupted state file.

#### Scenario: SIGINT during idle interval

- GIVEN the daemon is waiting between cycles
- WHEN SIGINT is received
- THEN "shutting down watch mode" is logged
- AND the process exits with code 0

#### Scenario: SIGTERM during active export

- GIVEN a cycle is actively writing observation files
- WHEN SIGTERM is received
- THEN the current cycle completes
- AND "shutting down watch mode" is logged
- AND the process exits with code 0
- AND `.engram-sync-state.json` is valid (not corrupted)

---

### Requirement: REQ-WATCH-06 — Flag Interaction

`--watch` MUST be compatible with `--project`, `--graph-config`, `--since`, `--force`, and `--limit`. All flags apply to every cycle with these exceptions:
- `--since` is respected ONLY on the first cycle; subsequent cycles use the state file cutoff.
- `--graph-config` is applied ONLY on the first cycle; subsequent cycles skip the graph config step.

#### Scenario: --project scopes every cycle

- GIVEN `--watch --project engram`
- WHEN the daemon runs multiple cycles
- THEN each cycle exports only observations from project "engram"

#### Scenario: --since only on first cycle

- GIVEN `--watch --since 2026-01-01`
- WHEN the first cycle runs
- THEN observations since 2026-01-01 are exported
- WHEN the second cycle runs
- THEN the state file cutoff is used (not --since)

#### Scenario: --graph-config=force only on first cycle

- GIVEN `--watch --graph-config=force`
- WHEN the first cycle runs
- THEN `graph.json` is overwritten
- WHEN the second cycle runs
- THEN `graph.json` is NOT touched (graph-config step is skipped for subsequent cycles)

---

### Requirement: REQ-WATCH-07 — Mutual Exclusion

`--interval` without `--watch` MUST cause an error exit.
`--watch` with an invalid `--interval` MUST cause an error exit.

#### Scenario: --interval without --watch is rejected

- GIVEN the user runs `engram obsidian-export --vault ./v --interval 5m`
- WHEN the flags are parsed
- THEN an error "--interval requires --watch" is printed to stderr
- AND the process exits with code 1

#### Scenario: --watch with invalid interval is rejected

- GIVEN the user runs `engram obsidian-export --vault ./v --watch --interval 10s`
- WHEN the flags are parsed
- THEN an error "--interval must be >= 1m; got 10s" is printed to stderr
- AND the process exits with code 1
