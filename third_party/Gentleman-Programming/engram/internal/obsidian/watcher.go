package obsidian

import (
	"context"
	"log"
	"time"
)

// Exportable is the interface the Watcher uses to run export cycles.
// *Exporter satisfies this interface. Fake implementations are used in tests.
type Exportable interface {
	// Export performs one export cycle and returns its result.
	Export() (*ExportResult, error)
	// SetGraphConfig overrides the GraphConfig mode for subsequent calls.
	SetGraphConfig(mode GraphConfigMode)
	// GraphConfig returns the current GraphConfig mode.
	GraphConfig() GraphConfigMode
}

// Watcher wraps an Exportable and runs it on a fixed interval in a loop.
type Watcher struct {
	exporter Exportable
	interval time.Duration
	logf     func(format string, args ...any) // injectable; defaults to log.Printf
}

// WatcherConfig holds constructor parameters for Watcher.
type WatcherConfig struct {
	// Exporter is the export driver. Required.
	Exporter Exportable
	// Interval between cycles. Required.
	Interval time.Duration
	// Logf is the log function. Defaults to log.Printf if nil.
	Logf func(format string, args ...any)
}

// NewWatcher constructs a Watcher from the given config.
func NewWatcher(cfg WatcherConfig) *Watcher {
	logf := cfg.Logf
	if logf == nil {
		logf = log.Printf
	}
	return &Watcher{
		exporter: cfg.Exporter,
		interval: cfg.Interval,
		logf:     logf,
	}
}

// Run executes the watch loop.
//
// Behavior:
//   - The first cycle runs immediately (REQ-WATCH-03).
//   - Subsequent cycles run after each interval tick.
//   - On cycle error: the error is logged and the loop continues (REQ-WATCH-04).
//   - The first cycle uses the exporter's current GraphConfig;
//     subsequent cycles force GraphConfigSkip (REQ-WATCH-06).
//   - Returns ctx.Err() when the context is canceled or times out (REQ-WATCH-05).
func (w *Watcher) Run(ctx context.Context) error {
	// Save the original GraphConfig for the first cycle.
	originalGraphConfig := w.exporter.GraphConfig()

	// ── First cycle: immediate ────────────────────────────────────────────────
	w.runCycle()

	// ── Subsequent cycles: override graph config to skip ─────────────────────
	// REQ-WATCH-06: --graph-config applies only to the first cycle.
	w.exporter.SetGraphConfig(GraphConfigSkip)

	// ── Ticker loop ───────────────────────────────────────────────────────────
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	_ = originalGraphConfig // used above; keep reference for clarity

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.runCycle()
		}
	}
}

// runCycle executes one export cycle and logs the result.
func (w *Watcher) runCycle() {
	result, err := w.exporter.Export()
	if err != nil {
		w.logf("[%s] sync error: %v", time.Now().UTC().Format(time.RFC3339), err)
		return
	}
	w.logf("[%s] sync: created=%d updated=%d deleted=%d skipped=%d hubs=%d",
		time.Now().UTC().Format(time.RFC3339),
		result.Created,
		result.Updated,
		result.Deleted,
		result.Skipped,
		result.HubsCreated,
	)
}
