package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ─── extractRepoName unit tests ──────────────────────────────────────────────

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "SSH with .git suffix",
			url:  "git@github.com:user/repo.git",
			want: "repo",
		},
		{
			name: "SSH without .git suffix",
			url:  "git@github.com:user/repo",
			want: "repo",
		},
		{
			name: "HTTPS with .git suffix",
			url:  "https://github.com/user/repo.git",
			want: "repo",
		},
		{
			name: "HTTPS without .git suffix",
			url:  "https://github.com/user/repo",
			want: "repo",
		},
		{
			name: "SSH org with dots in repo name",
			url:  "git@github.com:Gentleman-Programming/engram.git",
			want: "engram",
		},
		{
			name: "HTTPS org with dots",
			url:  "https://github.com/Gentleman-Programming/engram.git",
			want: "engram",
		},
		{
			name: "Repo name without .git from SSH",
			url:  "git@gitlab.com:group/subgroup/my-project",
			want: "my-project",
		},
		{
			name: "Empty URL returns empty",
			url:  "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractRepoName(tc.url)
			if got != tc.want {
				t.Errorf("extractRepoName(%q) = %q; want %q", tc.url, got, tc.want)
			}
		})
	}
}

// ─── DetectProject integration tests ─────────────────────────────────────────

// initGit initialises a new git repository in dir. Helper for tests.
func initGit(t *testing.T, dir string) {
	t.Helper()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test User")
}

func TestDetectProject_GitRemote(t *testing.T) {
	dir := t.TempDir()
	initGit(t, dir)

	// Add a fake remote
	cmd := exec.Command("git", "-C", dir, "remote", "add", "origin",
		"git@github.com:testuser/my-cool-repo.git")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git remote add: %v\n%s", err, out)
	}

	got := DetectProject(dir)
	if got != "my-cool-repo" {
		t.Errorf("DetectProject with remote = %q; want %q", got, "my-cool-repo")
	}
}

func TestDetectProject_GitRemote_HTTPS(t *testing.T) {
	dir := t.TempDir()
	initGit(t, dir)

	cmd := exec.Command("git", "-C", dir, "remote", "add", "origin",
		"https://github.com/Gentleman-Programming/engram.git")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git remote add: %v\n%s", err, out)
	}

	got := DetectProject(dir)
	if got != "engram" {
		t.Errorf("DetectProject HTTPS remote = %q; want %q", got, "engram")
	}
}

func TestDetectProject_GitRootNoRemote(t *testing.T) {
	dir := t.TempDir()
	initGit(t, dir)
	// No remote configured — should fall back to basename of git root

	got := DetectProject(dir)
	want := filepath.Base(dir)
	// Normalize to lowercase to match DetectProject output
	wantLower := strings.ToLower(want)
	if got != wantLower {
		t.Errorf("DetectProject no-remote = %q; want %q", got, wantLower)
	}
}

func TestDetectProject_NonGitDir(t *testing.T) {
	dir := t.TempDir()
	// Not a git repo — should fall back to basename of dir

	got := DetectProject(dir)
	want := strings.ToLower(filepath.Base(dir))
	if got != want {
		t.Errorf("DetectProject non-git = %q; want %q", got, want)
	}
}

func TestDetectProject_EmptyDir_NoPanic(t *testing.T) {
	// Even an empty string for dir should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DetectProject panicked: %v", r)
		}
	}()
	got := DetectProject("")
	// Just verify it returns something non-empty (the exact value depends on OS)
	if got == "" {
		t.Error("DetectProject(\"\") returned empty string")
	}
}

func TestDetectProject_NormalizedLowercase(t *testing.T) {
	// DetectProject must always return lowercase names.
	// Create a temp dir whose basename has upper-case letters.
	parent := t.TempDir()
	upper := filepath.Join(parent, "MyProject")
	if err := os.MkdirAll(upper, 0o755); err != nil {
		t.Fatal(err)
	}

	got := DetectProject(upper)
	if got != "myproject" {
		t.Errorf("DetectProject uppercase dir = %q; want %q", got, "myproject")
	}
}

func TestDetectProject_GitRemoteCasing(t *testing.T) {
	// Remote repo name like "MyRepo.git" should be lowercased.
	dir := t.TempDir()
	initGit(t, dir)

	cmd := exec.Command("git", "-C", dir, "remote", "add", "origin",
		"git@github.com:user/MyRepo.git")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git remote add: %v\n%s", err, out)
	}

	got := DetectProject(dir)
	if got != "myrepo" {
		t.Errorf("DetectProject uppercase remote name = %q; want %q", got, "myrepo")
	}
}
