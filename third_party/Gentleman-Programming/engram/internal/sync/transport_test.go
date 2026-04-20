package sync

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFileTransportReadManifestMissing(t *testing.T) {
	ft := NewFileTransport(t.TempDir())
	m, err := ft.ReadManifest()
	if err != nil {
		t.Fatalf("read missing manifest: %v", err)
	}
	if m.Version != 1 || len(m.Chunks) != 0 {
		t.Fatalf("expected empty v1 manifest, got %+v", m)
	}
}

func TestFileTransportManifestRoundtrip(t *testing.T) {
	dir := t.TempDir()
	ft := NewFileTransport(dir)

	want := &Manifest{
		Version: 1,
		Chunks: []ChunkEntry{
			{ID: "abc12345", CreatedBy: "alice", CreatedAt: "2025-06-01T00:00:00Z", Sessions: 1, Memories: 2, Prompts: 3},
		},
	}

	if err := ft.WriteManifest(want); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	got, err := ft.ReadManifest()
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if len(got.Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(got.Chunks))
	}
	if got.Chunks[0].ID != "abc12345" {
		t.Fatalf("chunk ID mismatch: got %q", got.Chunks[0].ID)
	}
	if got.Chunks[0].Memories != 2 {
		t.Fatalf("memories mismatch: got %d", got.Chunks[0].Memories)
	}
}

func TestFileTransportReadManifestInvalid(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	ft := NewFileTransport(dir)
	_, err := ft.ReadManifest()
	if err == nil {
		t.Fatal("expected error for invalid manifest JSON")
	}
}

func TestFileTransportReadManifestNotDir(t *testing.T) {
	// syncDir is a regular file, not a directory.
	badPath := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(badPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	ft := NewFileTransport(badPath)
	_, err := ft.ReadManifest()
	if err == nil {
		t.Fatal("expected error when syncDir is a file")
	}
}

func TestFileTransportChunkRoundtrip(t *testing.T) {
	dir := t.TempDir()
	ft := NewFileTransport(dir)

	payload := []byte(`{"sessions":[{"id":"s1"}],"observations":[],"prompts":[]}`)
	entry := ChunkEntry{ID: "aabbccdd", CreatedBy: "bob"}

	if err := ft.WriteChunk("aabbccdd", payload, entry); err != nil {
		t.Fatalf("write chunk: %v", err)
	}

	// Verify file was created.
	chunkPath := filepath.Join(dir, "chunks", "aabbccdd.jsonl.gz")
	if _, err := os.Stat(chunkPath); err != nil {
		t.Fatalf("chunk file not created: %v", err)
	}

	got, err := ft.ReadChunk("aabbccdd")
	if err != nil {
		t.Fatalf("read chunk: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("chunk data mismatch: got %q want %q", got, payload)
	}
}

func TestFileTransportReadChunkMissing(t *testing.T) {
	ft := NewFileTransport(t.TempDir())
	_, err := ft.ReadChunk("nonexist")
	if err == nil {
		t.Fatal("expected error for missing chunk")
	}
}

func TestFileTransportWriteChunkBadDir(t *testing.T) {
	badPath := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(badPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	ft := NewFileTransport(badPath)
	err := ft.WriteChunk("aabb", []byte("data"), ChunkEntry{})
	if err == nil {
		t.Fatal("expected error writing chunk to non-directory syncDir")
	}
}

func TestFileTransportWriteManifestMarshalError(t *testing.T) {
	resetSyncTestHooks(t)
	ft := NewFileTransport(t.TempDir())

	jsonMarshalManifest = func(v any, prefix, indent string) ([]byte, error) {
		_, _ = prefix, indent
		_ = v
		return nil, os.ErrInvalid
	}

	err := ft.WriteManifest(&Manifest{Version: 1})
	if err == nil {
		t.Fatal("expected error from marshal failure")
	}
}

func TestFileTransportChunkDataIntegrity(t *testing.T) {
	// Write a real ChunkData payload, read it back, verify JSON roundtrip.
	dir := t.TempDir()
	ft := NewFileTransport(dir)

	original := map[string]any{
		"sessions": []map[string]any{
			{"id": "s1", "project": "proj", "directory": "/tmp", "started_at": "2025-01-01 00:00:00"},
		},
		"observations": []map[string]any{
			{"id": 1, "session_id": "s1", "type": "decision", "title": "test", "content": "data", "scope": "project", "created_at": "2025-01-01 00:00:00", "updated_at": "2025-01-01 00:00:00"},
		},
		"prompts": []map[string]any{
			{"id": 1, "session_id": "s1", "content": "hello", "created_at": "2025-01-01 00:00:00"},
		},
	}

	payload, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if err := ft.WriteChunk("deadbeef", payload, ChunkEntry{ID: "deadbeef"}); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := ft.ReadChunk("deadbeef")
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("unmarshal roundtrip: %v", err)
	}

	sessions, ok := result["sessions"].([]any)
	if !ok || len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %v", result["sessions"])
	}
}

func TestNewFileTransport(t *testing.T) {
	dir := t.TempDir()
	ft := NewFileTransport(dir)
	if ft == nil {
		t.Fatal("expected non-nil FileTransport")
	}
	if ft.syncDir != dir {
		t.Fatalf("syncDir mismatch: got %q want %q", ft.syncDir, dir)
	}
}

func TestNewWithTransport(t *testing.T) {
	s := newTestStore(t)
	ft := NewFileTransport(t.TempDir())
	sy := NewWithTransport(s, ft)
	if sy == nil {
		t.Fatal("expected non-nil syncer")
	}
	if sy.transport != ft {
		t.Fatal("transport not preserved")
	}
}

func TestNewLocal(t *testing.T) {
	s := newTestStore(t)
	dir := t.TempDir()
	sy := NewLocal(s, dir)
	if sy == nil {
		t.Fatal("expected non-nil syncer")
	}
	if sy.syncDir != dir {
		t.Fatalf("syncDir mismatch: got %q want %q", sy.syncDir, dir)
	}
	if sy.transport == nil {
		t.Fatal("expected non-nil transport")
	}
}
