# Proposal: obsidian-auto-sync — graph config bootstrap + watch mode

## Intent

After the `obsidian-plugin` smoke test (1731 obs + 277 hubs exported), two friction points remain: (1) users must set up cron/launchd to keep the vault in sync, and (2) Obsidian's default graph view is unreadable for our memory model without manual tuning. This change makes `engram obsidian-export` both self-syncing and visually opinionated out of the box.

## Scope

### In Scope
- `--graph-config=<preserve|force|skip>` flag (default `preserve`) that writes `.obsidian/graph.json` to the vault
- Embedded default `graph.json` template (color groups, forces, display) covering `_sessions`, `_topics`, and type tags (`#architecture`, `#bugfix`, `#decision`, `#discovery`, `#pattern`)
- `--watch` + `--interval <duration>` flags turning the command into a daemon
- Immediate-first-run semantics: export once, then tick
- Minimum interval validation (≥1m), default `10m`
- Graceful shutdown on SIGINT/SIGTERM with in-flight cycle drain
- Per-cycle log line: timestamp + created/updated/deleted/skipped counts
- `graph_test.go` unit coverage
- README flag documentation

### Out of Scope
- Integration with `engram serve` (daemon colocation is future work)
- `fsnotify`/inotify reactive file-watching mode
- Cross-vault or multi-vault sync
- Bidirectional sync (vault → engram)

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `obsidian-export`: adds graph-config bootstrap behavior and watch-mode daemon semantics to the existing CLI contract

## Approach

1. Add `graph.json` embedded template via `//go:embed` in `internal/obsidian/graph.go`, plus a `WriteGraphConfig(vault, mode)` function with three modes.
2. Extend `Exporter` to accept a graph-config mode and call the writer once per run before the observation pass.
3. In `cmd/engram/main.go`, parse `--graph-config`, `--watch`, `--interval`; wrap the single-run export in a `time.Ticker` loop when `--watch` is set; register `signal.NotifyContext(os.Interrupt, syscall.SIGTERM)` for graceful shutdown.
4. First run executes immediately on startup, then every interval; interval validated `>= 1 * time.Minute`.
5. Log each cycle with the existing `ExportResult` counts.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `cmd/engram/main.go` | Modified | Parse 3 new flags, daemon loop, signal handling |
| `internal/obsidian/exporter.go` | Modified | Accept `GraphConfigMode`, invoke writer once per run |
| `internal/obsidian/graph.go` | New | Embedded `graph.json` + `WriteGraphConfig` |
| `internal/obsidian/graph.json` | New | Default graph view template |
| `internal/obsidian/graph_test.go` | New | Unit tests for the three modes |
| `README.md` | Modified | Document new flags + watch-mode example |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Overwriting a user's customized `graph.json` | Med | Default mode `preserve`; `force` is explicit opt-in |
| Watch loop leaks or hangs on shutdown | Low | `signal.NotifyContext` + `ticker.Stop()` + drain in-flight run |
| Too-frequent interval hammers SQLite | Low | Enforce `interval >= 1m`; incremental sync skips unchanged observations |
| Embedded template drifts from real Obsidian schema | Low | Unit test validates JSON parses + has required top-level keys |

## Rollback Plan

Purely additive. Rollback = delete `internal/obsidian/graph.go`, `graph.json`, `graph_test.go`, and revert the flag parsing + daemon block in `cmd/engram/main.go`. No data migration, no store changes, existing single-run exports remain byte-identical.

## Dependencies

- None (`time.Ticker`, `os/signal`, `embed` are stdlib)

## Success Criteria

- [ ] `engram obsidian-export --vault ./vault` writes `.obsidian/graph.json` on a fresh vault and leaves it alone on re-run
- [ ] `--graph-config=force` overwrites existing graph.json
- [ ] `--graph-config=skip` never touches `.obsidian/`
- [ ] `engram obsidian-export --vault ./vault --watch --interval 2m` runs first export immediately, then again after 2 minutes
- [ ] SIGINT during watch mode exits cleanly within one cycle
- [ ] `--interval 30s` is rejected with a clear error
- [ ] `go test ./internal/obsidian/...` passes; existing exporter tests unchanged
