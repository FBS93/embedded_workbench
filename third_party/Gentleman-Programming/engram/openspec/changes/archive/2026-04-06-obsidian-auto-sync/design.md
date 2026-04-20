# Design: obsidian-auto-sync — graph bootstrap + watch mode

## Technical Approach

Purely additive extension to `internal/obsidian` and `cmd/engram/main.go`. A new `graph.go` holds an embedded `graph.json` template (the user's locked-in engram-default) plus `WriteGraphConfig(vault, mode)`. A new `watcher.go` wraps the existing `Exporter` in a `time.Ticker` loop driven by `signal.NotifyContext`. `cmdObsidianExport` parses three new flags (`--graph-config`, `--watch`, `--interval`) and dispatches to single-run or watch-loop paths. No store changes, no migrations, no new deps.

## Architecture Decisions

| Decision | Choice | Alternative | Rationale |
|---|---|---|---|
| Watcher location | `internal/obsidian/watcher.go` (own type) | Inline loop in `main.go` | Testable in isolation with fake clock + counter exporter; matches `newObsidianExporter` injection pattern |
| Signal handling | `signal.NotifyContext(os.Interrupt, syscall.SIGTERM)` | Manual `signal.Notify` + channel | Go 1.16+ idiom; context propagates naturally to `ticker.C` select |
| Graph config in loop | First cycle only (transient `firstCycle` flag on `ExportConfig`) | Every cycle | Prevents clobbering user customizations made mid-session, even in `force` mode |
| Embedded template | `//go:embed graph.json` | Runtime file read | Versioned with binary, zero runtime deps, matches goreleaser `CGO_ENABLED=0` build |
| Default mode | `preserve` | `force` | Safe default — respects any existing `.obsidian/graph.json` |
| Watcher API | `Run(ctx context.Context) error` | `Start()` / `Stop()` | Context idiom; no internal state machine; cancellation via parent ctx |
| Error handling in loop | Log + continue | Exit daemon on error | Transient store/FS errors shouldn't kill a long-running daemon |
| First run | Immediate, then tick | Wait for first tick | User sees export feedback without waiting 10 minutes |
| Min interval | `>= 1 * time.Minute` | No minimum | Prevents hammering SQLite; matches proposal success criteria |
| `graph.json` placement | `{vault}/.obsidian/graph.json` | `{vault}/engram/.obsidian/graph.json` | Obsidian reads the config at the vault root `.obsidian/`, not per-folder |

## Data Flow

```
  cmdObsidianExport
         │
         ├── parse flags (--vault --graph-config --watch --interval ...)
         ├── validate (interval>=1m, mode in {preserve,force,skip})
         │
         ├── !watch ──→ runSingleCycle(exporter, firstCycle=true)
         │                    │
         │                    ├── WriteGraphConfig(vault, mode)   [first cycle only]
         │                    └── Exporter.Export() ──→ observations + hubs
         │
         └── watch ───→ signal.NotifyContext(SIGINT, SIGTERM)
                            │
                            └── Watcher.Run(ctx)
                                     │
                                     ├── runCycle(firstCycle=true)    [immediate]
                                     └── for { select ctx.Done | ticker.C → runCycle(first=false) }
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/obsidian/graph.go` | Create | `GraphConfigMode` enum, `ParseGraphConfigMode`, `//go:embed graph.json`, `WriteGraphConfig(vault, mode)` |
| `internal/obsidian/graph.json` | Create | Embedded default template — exact user-locked config (6 color groups, centerStrength 0.515, repelStrength 12.71, linkStrength 0.729, linkDistance 207, scale 0.1, no arrows, no text fade) |
| `internal/obsidian/graph_test.go` | Create | Unit tests for `ParseGraphConfigMode`, all 3 modes of `WriteGraphConfig`, JSON validity + value assertions on embedded template |
| `internal/obsidian/watcher.go` | Create | `Watcher` struct wrapping `*Exporter`, `Run(ctx)` with immediate-first + ticker loop, log line per cycle |
| `internal/obsidian/watcher_test.go` | Create | Test `Run()` with 10ms interval + context cancel after N cycles, counter-based fake exporter |
| `internal/obsidian/exporter.go` | Modify | Add `GraphConfig GraphConfigMode` + unexported `firstCycle bool` to `ExportConfig`; call `WriteGraphConfig` at top of `Export()` when `firstCycle && mode != Skip` |
| `cmd/engram/main.go` | Modify | Parse `--graph-config`, `--watch`, `--interval`; add `newObsidianWatcher` injectable; dispatch single vs watch |
| `cmd/engram/main_test.go` | Modify | Cover new flag parsing, rejection of `--interval 30s`, rejection of bare `--interval` without `--watch`, invalid `--graph-config` value |
| `README.md` | Modify | Document new flags + watch-mode example |

## Interfaces / Contracts

```go
// internal/obsidian/graph.go
type GraphConfigMode string
const (
    GraphConfigPreserve GraphConfigMode = "preserve" // default
    GraphConfigForce    GraphConfigMode = "force"
    GraphConfigSkip     GraphConfigMode = "skip"
)

//go:embed graph.json
var defaultGraphTemplate []byte

func ParseGraphConfigMode(s string) (GraphConfigMode, error)
func WriteGraphConfig(vaultPath string, mode GraphConfigMode) error

// internal/obsidian/exporter.go (additions)
type ExportConfig struct {
    // ... existing ...
    GraphConfig GraphConfigMode // default preserve
    firstCycle  bool            // unexported; set by Watcher / single-run path
}

// internal/obsidian/watcher.go
type Watcher struct {
    exporter *Exporter
    interval time.Duration
    logf     func(format string, args ...any) // injectable, defaults to log.Printf
}
func NewWatcher(exp *Exporter, interval time.Duration) *Watcher
func (w *Watcher) Run(ctx context.Context) error
```

`WriteGraphConfig` semantics:
- `preserve`: if `{vault}/.obsidian/graph.json` exists → no-op; else create parent dir + write template
- `force`: always `MkdirAll(.obsidian, 0755)` + `WriteFile(graph.json, template, 0644)`
- `skip`: no-op, returns nil

Flag validation errors (all exit 1 with stderr message):
- `--interval` without `--watch` → "error: --interval requires --watch"
- `--interval` < 1m → "error: --interval must be >= 1m"
- `--graph-config` not in set → "error: invalid --graph-config value %q (expected preserve|force|skip)"

## Testing Strategy

| Layer | What | Approach |
|---|---|---|
| Unit | `ParseGraphConfigMode` | Table test: valid/invalid inputs |
| Unit | `WriteGraphConfig` 3 modes | tempdir vault; assert file presence/content per mode; idempotency of `preserve` |
| Unit | Embedded template integrity | `json.Unmarshal` into map; assert 6 color groups + exact numeric values from user preference |
| Unit | `Watcher.Run` | Fake exporter (counter), interval 10ms, `context.WithTimeout(50ms)`; assert cycles >= 3 + ctx.Err() returned |
| Unit | `Watcher.Run` first-cycle | Assert first `runCycle` invoked before any tick; log line emitted |
| Unit | `cmdObsidianExport` flags | Inject `newObsidianExporter` + `newObsidianWatcher`; assert parsed values, error paths exit 1 with correct stderr |
| Integration | Existing exporter tests | Unchanged — new `GraphConfig` field defaults to `preserve` and no-ops on test vaults without `.obsidian/` expectations |

Daemon loop is NOT exercised via `cmdObsidianExport` in unit tests — we test `Watcher.Run` directly and stub `newObsidianWatcher` in `main_test.go`.

## Migration / Rollout

No migration. Purely additive flags with safe defaults (`preserve`, no watch). Existing single-run exports produce byte-identical output. Rollback = revert the files listed above.

## Open Questions

None — all decisions resolved against the proposal and the user's locked-in graph config.
