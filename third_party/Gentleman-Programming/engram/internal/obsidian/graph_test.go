package obsidian

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParseGraphConfigMode covers REQ-GRAPH-01: valid values are accepted,
// any other string returns an error. Parsing is case-sensitive.
func TestParseGraphConfigMode(t *testing.T) {
	tests := []struct {
		input   string
		want    GraphConfigMode
		wantErr bool
	}{
		{input: "preserve", want: GraphConfigPreserve, wantErr: false},
		{input: "force", want: GraphConfigForce, wantErr: false},
		{input: "skip", want: GraphConfigSkip, wantErr: false},
		{input: "invalid", want: "", wantErr: true},
		{input: "", want: "", wantErr: true},
		{input: "PRESERVE", want: "", wantErr: true}, // case-sensitive
		{input: "Force", want: "", wantErr: true},    // case-sensitive
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseGraphConfigMode(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseGraphConfigMode(%q): expected error, got nil (mode=%q)", tc.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseGraphConfigMode(%q): unexpected error: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("ParseGraphConfigMode(%q): got %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestEmbeddedGraphTemplate covers REQ-GRAPH-05: the embedded graph.json
// must parse as valid JSON and contain EXACTLY the user-locked values.
func TestEmbeddedGraphTemplate(t *testing.T) {
	var cfg map[string]any
	if err := json.Unmarshal(defaultGraphTemplate, &cfg); err != nil {
		t.Fatalf("defaultGraphTemplate is not valid JSON: %v", err)
	}

	// --- numeric fields ---
	assertFloat(t, cfg, "centerStrength", 0.515147569444444)
	assertFloat(t, cfg, "repelStrength", 12.7118055555556)
	assertFloat(t, cfg, "linkStrength", 0.729210069444444)
	assertFloat(t, cfg, "linkDistance", 207)
	assertFloat(t, cfg, "scale", 0.1)
	assertFloat(t, cfg, "textFadeMultiplier", 0)
	assertFloat(t, cfg, "nodeSizeMultiplier", 1)
	assertFloat(t, cfg, "lineSizeMultiplier", 1)

	// --- boolean fields ---
	assertBool(t, cfg, "showArrow", false)
	assertBool(t, cfg, "showOrphans", true)

	// --- collapse section booleans ---
	assertBool(t, cfg, "collapse-color-groups", true)
	assertBool(t, cfg, "collapse-filter", false)
	assertBool(t, cfg, "collapse-display", true)
	assertBool(t, cfg, "collapse-forces", false)

	// --- color groups: exactly 6, in exact order ---
	raw, ok := cfg["colorGroups"]
	if !ok {
		t.Fatal("colorGroups key missing from embedded graph.json")
	}
	groups, ok := raw.([]any)
	if !ok {
		t.Fatalf("colorGroups is not an array, got %T", raw)
	}
	if len(groups) != 6 {
		t.Fatalf("colorGroups: got %d entries, want exactly 6", len(groups))
	}

	want := []struct {
		query string
		rgb   float64
	}{
		{"path:engram/_sessions", 14736466},
		{"path:engram/_topics", 13893887},
		{"tag:#architecture", 7935},
		{"tag:#bugfix", 16711680},
		{"tag:#decision", 65322},
		{"tag:#pattern", 16741120},
	}

	for i, g := range groups {
		entry, ok := g.(map[string]any)
		if !ok {
			t.Errorf("colorGroups[%d]: not an object", i)
			continue
		}
		gotQuery, _ := entry["query"].(string)
		if gotQuery != want[i].query {
			t.Errorf("colorGroups[%d].query: got %q, want %q", i, gotQuery, want[i].query)
		}
		color, ok := entry["color"].(map[string]any)
		if !ok {
			t.Errorf("colorGroups[%d].color: not an object", i)
			continue
		}
		gotA, _ := color["a"].(float64)
		if gotA != 1 {
			t.Errorf("colorGroups[%d].color.a: got %v, want 1", i, gotA)
		}
		gotRGB, _ := color["rgb"].(float64)
		if gotRGB != want[i].rgb {
			t.Errorf("colorGroups[%d].color.rgb: got %v, want %v", i, gotRGB, want[i].rgb)
		}
	}
}

// TestWriteGraphConfig covers REQ-GRAPH-02..06:
// preserve mode, force mode, skip mode, and .obsidian/ directory creation.
func TestWriteGraphConfig(t *testing.T) {
	sentinel := []byte(`{"custom":true}`)

	t.Run("preserve: creates graph.json when absent", func(t *testing.T) {
		vault := t.TempDir()
		if err := WriteGraphConfig(vault, GraphConfigPreserve); err != nil {
			t.Fatalf("WriteGraphConfig preserve on empty vault: %v", err)
		}
		path := filepath.Join(vault, ".obsidian", "graph.json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("graph.json not created: %v", err)
		}
		// file must be valid JSON and match the embedded template
		if string(data) != string(defaultGraphTemplate) {
			t.Errorf("graph.json content differs from embedded template")
		}
	})

	t.Run("preserve: skips when graph.json already exists", func(t *testing.T) {
		vault := t.TempDir()
		obsDir := filepath.Join(vault, ".obsidian")
		if err := os.MkdirAll(obsDir, 0755); err != nil {
			t.Fatal(err)
		}
		graphPath := filepath.Join(obsDir, "graph.json")
		if err := os.WriteFile(graphPath, sentinel, 0644); err != nil {
			t.Fatal(err)
		}

		if err := WriteGraphConfig(vault, GraphConfigPreserve); err != nil {
			t.Fatalf("WriteGraphConfig preserve with existing file: %v", err)
		}

		data, err := os.ReadFile(graphPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != string(sentinel) {
			t.Errorf("preserve mode must not overwrite existing graph.json, got %s", data)
		}
	})

	t.Run("force: overwrites existing graph.json", func(t *testing.T) {
		vault := t.TempDir()
		obsDir := filepath.Join(vault, ".obsidian")
		if err := os.MkdirAll(obsDir, 0755); err != nil {
			t.Fatal(err)
		}
		graphPath := filepath.Join(obsDir, "graph.json")
		if err := os.WriteFile(graphPath, sentinel, 0644); err != nil {
			t.Fatal(err)
		}

		if err := WriteGraphConfig(vault, GraphConfigForce); err != nil {
			t.Fatalf("WriteGraphConfig force: %v", err)
		}

		data, err := os.ReadFile(graphPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != string(defaultGraphTemplate) {
			t.Errorf("force mode must overwrite with embedded template")
		}
	})

	t.Run("force: creates graph.json when absent", func(t *testing.T) {
		vault := t.TempDir()
		if err := WriteGraphConfig(vault, GraphConfigForce); err != nil {
			t.Fatalf("WriteGraphConfig force on empty vault: %v", err)
		}
		path := filepath.Join(vault, ".obsidian", "graph.json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("graph.json not created by force mode: %v", err)
		}
		if string(data) != string(defaultGraphTemplate) {
			t.Errorf("force mode must write embedded template when file absent")
		}
	})

	t.Run("skip: never creates graph.json", func(t *testing.T) {
		vault := t.TempDir()
		if err := WriteGraphConfig(vault, GraphConfigSkip); err != nil {
			t.Fatalf("WriteGraphConfig skip: %v", err)
		}
		path := filepath.Join(vault, ".obsidian", "graph.json")
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("skip mode must not create graph.json, but it exists (err=%v)", err)
		}
	})

	t.Run("obsidian dir created when missing (preserve)", func(t *testing.T) {
		vault := t.TempDir()
		// .obsidian does NOT exist
		if err := WriteGraphConfig(vault, GraphConfigPreserve); err != nil {
			t.Fatalf("WriteGraphConfig should create .obsidian dir: %v", err)
		}
		obsDir := filepath.Join(vault, ".obsidian")
		if info, err := os.Stat(obsDir); os.IsNotExist(err) {
			t.Errorf(".obsidian dir was not created")
		} else if err == nil && !info.IsDir() {
			t.Errorf(".obsidian is not a directory")
		}
	})
}

// --- helpers ---

func assertFloat(t *testing.T, cfg map[string]any, key string, want float64) {
	t.Helper()
	raw, ok := cfg[key]
	if !ok {
		t.Errorf("key %q missing from graph.json", key)
		return
	}
	got, ok := raw.(float64)
	if !ok {
		t.Errorf("key %q: expected float64, got %T (%v)", key, raw, raw)
		return
	}
	// Use a tight epsilon for precision-sensitive values
	const epsilon = 1e-9
	diff := got - want
	if diff < -epsilon || diff > epsilon {
		t.Errorf("key %q: got %.15f, want %.15f", key, got, want)
	}
}

func assertBool(t *testing.T, cfg map[string]any, key string, want bool) {
	t.Helper()
	raw, ok := cfg[key]
	if !ok {
		t.Errorf("key %q missing from graph.json", key)
		return
	}
	got, ok := raw.(bool)
	if !ok {
		t.Errorf("key %q: expected bool, got %T (%v)", key, raw, raw)
		return
	}
	if got != want {
		t.Errorf("key %q: got %v, want %v", key, got, want)
	}
}
