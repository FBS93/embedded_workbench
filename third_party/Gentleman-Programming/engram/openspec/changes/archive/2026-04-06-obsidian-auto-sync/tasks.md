# Tasks: obsidian-auto-sync

## Phase 1: Graph Config Bootstrap

- [x] 1.1 **[RED]** Create `internal/obsidian/graph_test.go` — `TestParseGraphConfigMode`: table test for `preserve`/`force`/`skip` (ok) and any other string (error). (REQ-GRAPH-01)
- [x] 1.2 **[GREEN]** Create `internal/obsidian/graph.go` — `GraphConfigMode` type + constants + `ParseGraphConfigMode(s string) (GraphConfigMode, error)`.
- [x] 1.3 **[RED]** Add `TestEmbeddedGraphTemplate` to `graph_test.go` — unmarshal embedded bytes; assert 6 color groups with exact queries + rgb values; assert all numeric keys match REQ-GRAPH-05 table. (REQ-GRAPH-05)
- [x] 1.4 **[GREEN]** Create `internal/obsidian/graph.json` with exact user-locked values (6 color groups, centerStrength 0.515147569444444, etc.) + add `//go:embed graph.json` var to `graph.go`.
- [x] 1.5 **[RED]** Add `TestWriteGraphConfig` to `graph_test.go` — tempdir vault; assert preserve creates on missing, skips on present; assert force always overwrites; assert skip never creates; assert `.obsidian/` dir is created when absent. (REQ-GRAPH-02..06)
- [x] 1.6 **[GREEN]** Implement `WriteGraphConfig(vaultPath string, mode GraphConfigMode) error` in `graph.go`.

## Phase 2: Exporter Integration

- [x] 2.1 **[RED]** Add `TestExporterCallsGraphConfig` to `internal/obsidian/exporter_test.go` — construct `ExportConfig{GraphConfig: GraphConfigForce}`, stub fs, verify `.obsidian/graph.json` is created; verify `GraphConfigSkip` does NOT create it; verify zero-value defaults to skip (backward compat). (REQ-WATCH-06)
- [x] 2.2 **[GREEN]** Modify `internal/obsidian/exporter.go` — add `GraphConfig GraphConfigMode` to `ExportConfig`; add `SetGraphConfig`/`GraphConfig()` methods to `*Exporter` (satisfies `Exportable`); call `WriteGraphConfig(vault, graphMode)` at top of `Export()` with empty string defaulting to `GraphConfigSkip`.

## Phase 3: Watcher

- [x] 3.1 **[RED]** Create `internal/obsidian/watcher_test.go` — `TestWatcherRunsImmediatelyThenTicks`: fake counter exporter, interval 10ms, `context.WithTimeout(55ms)`; assert cycles ≥ 3 and first cycle fires before first tick. (REQ-WATCH-03)
- [x] 3.2 **[GREEN]** Create `internal/obsidian/watcher.go` — `Exportable` interface + `Watcher` struct + `NewWatcher(WatcherConfig) *Watcher` + `Run(ctx context.Context) error` with immediate first run then ticker loop.
- [x] 3.3 **[RED]** Add `TestWatcherLogsCycleResults` to `watcher_test.go` — inject `logf` capture; assert log line matches `[{RFC3339}] sync: created=N updated=N deleted=N skipped=N hubs=N`. (REQ-WATCH-04)
- [x] 3.4 **[GREEN]** Add per-cycle log line in `Run()` using injectable `logf func(format string, args ...any)`.
- [x] 3.5 **[RED]** Add `TestWatcherContinuesOnError` — fake exporter returns error on first call; assert loop continues, second cycle increments counter. (REQ-WATCH-04)
- [x] 3.6 **[GREEN]** Add error branch in `Run()`: log error + continue (no return).
- [x] 3.7 **[RED]** Add `TestWatcherGracefulShutdown` — cancel context during ticker sleep; assert `Run()` returns `ctx.Err()` without panic or goroutine leak. (REQ-WATCH-05)
- [x] 3.8 **[GREEN]** Ensure `Run()` select block returns cleanly on `ctx.Done()` at both sleep and mid-cycle points.
- [x] 3.9 **[RED]** Add `TestWatcherFirstCycleGraphConfigOnly` — assert cycle 1 has original GraphConfig value (GraphConfigForce), cycles 2+ have GraphConfigSkip. (REQ-WATCH-06)
- [x] 3.10 **[GREEN]** Implement first-cycle flag logic inside `Run()` — save original GraphConfig, run first cycle, then call `SetGraphConfig(GraphConfigSkip)` for all subsequent cycles.

## Phase 4: CLI Wiring

- [x] 4.1 **[RED]** Modify `cmd/engram/main_test.go` — add cases: `--graph-config invalid` exits 1 with correct stderr; `--interval 5m` without `--watch` exits 1; `--watch --interval 30s` exits 1; `--watch --interval 2m --vault ...` invokes injected `newObsidianWatcher`. (REQ-GRAPH-01, REQ-WATCH-02, REQ-WATCH-07)
- [x] 4.2 **[GREEN]** Modify `cmd/engram/main.go` — parse `--graph-config` (default `preserve`), `--watch` (bool), `--interval` (default `10m`); validate (`ParseGraphConfigMode`, interval ≥ 1m, `--interval` requires `--watch`); add `var newObsidianWatcher = obsidian.NewWatcher` injectable; dispatch single-run vs `signal.NotifyContext(SIGINT, SIGTERM)` → `watcher.Run(ctx)` + log "shutting down watch mode".

## Phase 5: Documentation

- [x] 5.1 Modify `README.md` — add `--graph-config`, `--watch`, `--interval` rows to the flags table; add "Auto-sync" section with `engram obsidian-export --vault ~/Obsidian/engram --watch` example.
