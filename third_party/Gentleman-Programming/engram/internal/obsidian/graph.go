package obsidian

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// GraphConfigMode controls how WriteGraphConfig handles an existing graph.json.
type GraphConfigMode string

const (
	// GraphConfigPreserve writes the default template only if graph.json is absent (default).
	GraphConfigPreserve GraphConfigMode = "preserve"

	// GraphConfigForce always overwrites graph.json with the embedded default.
	GraphConfigForce GraphConfigMode = "force"

	// GraphConfigSkip never reads, writes, or creates graph.json.
	GraphConfigSkip GraphConfigMode = "skip"
)

//go:embed graph.json
var defaultGraphTemplate []byte

// ParseGraphConfigMode parses s into a GraphConfigMode.
// Returns an error for any value not in the accepted set {preserve, force, skip}.
// Parsing is case-sensitive.
func ParseGraphConfigMode(s string) (GraphConfigMode, error) {
	switch s {
	case string(GraphConfigPreserve):
		return GraphConfigPreserve, nil
	case string(GraphConfigForce):
		return GraphConfigForce, nil
	case string(GraphConfigSkip):
		return GraphConfigSkip, nil
	default:
		return "", fmt.Errorf("invalid --graph-config value: %s (accepted: preserve, force, skip)", s)
	}
}

// WriteGraphConfig writes the embedded graph.json default into {vaultPath}/.obsidian/graph.json
// according to the given mode.
//
//   - preserve: creates the file only when it does not already exist.
//   - force:    always creates or overwrites the file with the embedded default.
//   - skip:     no-op; returns nil immediately.
//
// The .obsidian/ directory is created with 0755 permissions if it does not exist
// (except in skip mode where nothing is written).
func WriteGraphConfig(vaultPath string, mode GraphConfigMode) error {
	if mode == GraphConfigSkip {
		return nil
	}

	obsidianDir := filepath.Join(vaultPath, ".obsidian")
	graphPath := filepath.Join(obsidianDir, "graph.json")

	if mode == GraphConfigPreserve {
		if _, err := os.Stat(graphPath); err == nil {
			// File exists — preserve it
			return nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("checking graph.json: %w", err)
		}
		// File does not exist — fall through to create
	}

	// force or preserve-with-absent-file: create dir + write
	if err := os.MkdirAll(obsidianDir, 0755); err != nil {
		return fmt.Errorf("creating .obsidian dir: %w", err)
	}
	if err := os.WriteFile(graphPath, defaultGraphTemplate, 0644); err != nil {
		return fmt.Errorf("writing graph.json: %w", err)
	}
	return nil
}
