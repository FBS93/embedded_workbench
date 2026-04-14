package upgrade

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestEnumerateFilesInDir_SkipsSymlinks verifies that enumerateFilesInDir does not
// traverse into symlinks pointing to directories (simulates Windows junctions /
// reparse points). Directory symlinks are skipped entirely.
//
// Note: filepath.WalkDir does NOT follow symlinks — it reports them as entries
// with d.Type()&ModeSymlink set. Without the symlink check, a symlink-to-dir
// would pass the !d.IsDir() check and be added to the file list, causing
// "is a directory" errors when snapshotPath attempts to copy it as a file.
//
// On Windows, os.Symlink requires elevated privileges or developer mode enabled.
// We skip on Windows since the fix uses a platform-safe check.
func TestEnumerateFilesInDir_SkipsSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Symlink requires elevated privileges on Windows — skipping symlink traversal test")
	}

	home := t.TempDir()

	// Create a real file inside home that should be discovered.
	realFile := filepath.Join(home, "config.json")
	if err := os.WriteFile(realFile, []byte(`{"key":"value"}`), 0o644); err != nil {
		t.Fatalf("WriteFile real file: %v", err)
	}

	// Create a target directory (simulating a junctioned skills directory).
	targetDir := filepath.Join(home, "skills-target")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("MkdirAll target dir: %v", err)
	}
	// Create a file inside the target dir so WalkDir finds something if it follows the link.
	if err := os.WriteFile(filepath.Join(targetDir, "skill.md"), []byte("# Skill"), 0o644); err != nil {
		t.Fatalf("WriteFile inside target dir: %v", err)
	}

	// Create a symlink inside home that points to the target directory.
	// This simulates a Windows junction / reparse point.
	symlinkPath := filepath.Join(home, "skills-link")
	if err := os.Symlink(targetDir, symlinkPath); err != nil {
		t.Fatalf("os.Symlink: %v", err)
	}

	// enumerateFilesInDir must not return the symlink itself or files from inside the symlink target.
	files, err := enumerateFilesInDir(home, nil)
	if err != nil {
		t.Fatalf("enumerateFilesInDir() with symlink returned error: %v — must not fail on symlinks", err)
	}

	// The real file must be in the list.
	realFileFound := false
	for _, f := range files {
		if f == realFile {
			realFileFound = true
		}
		// The symlink ITSELF must NOT be included — on macOS, WalkDir returns the
		// symlink as a non-directory entry (d.IsDir()=false, d.Type()=ModeSymlink).
		// If included, snapshotPath will call os.Stat which resolves to a directory → "is a directory" error.
		if f == symlinkPath {
			t.Errorf("enumerateFilesInDir included symlink path %q — symlinks must be skipped to avoid 'is a directory' error in snapshot", f)
		}
		// Must NOT include files reached by following the symlink into the target dir.
		symlinkTarget := filepath.Join(symlinkPath, "skill.md")
		if f == symlinkTarget {
			t.Errorf("enumerateFilesInDir traversed into symlink target and included %q — junctions/symlinks must be skipped", f)
		}
	}

	if !realFileFound {
		t.Errorf("enumerateFilesInDir did not return the real file %q; got: %v", realFile, files)
	}
}

// TestEnumerateFilesInDir_SymlinkToFileIsIncluded verifies that symlinks pointing
// to regular files ARE included in the backup. This supports dotfile managers
// (stow, chezmoi, bare git) where config files like CLAUDE.md are symlinks to
// files in a dotfiles repository.
func TestEnumerateFilesInDir_SymlinkToFileIsIncluded(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Symlink requires elevated privileges on Windows")
	}

	home := t.TempDir()

	// Create a real file outside the config dir (simulates dotfiles repo).
	dotfilesDir := filepath.Join(home, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0o755); err != nil {
		t.Fatalf("MkdirAll dotfiles: %v", err)
	}
	realConfig := filepath.Join(dotfilesDir, "CLAUDE.md")
	if err := os.WriteFile(realConfig, []byte("# My Claude rules"), 0o644); err != nil {
		t.Fatalf("WriteFile real config: %v", err)
	}

	// Create the config dir with a symlink-to-file.
	configDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll config dir: %v", err)
	}
	symlinkFile := filepath.Join(configDir, "CLAUDE.md")
	if err := os.Symlink(realConfig, symlinkFile); err != nil {
		t.Fatalf("os.Symlink file: %v", err)
	}

	// Also create a regular file next to the symlink.
	regularFile := filepath.Join(configDir, "settings.json")
	if err := os.WriteFile(regularFile, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("WriteFile regular: %v", err)
	}

	files, err := enumerateFilesInDir(configDir, nil)
	if err != nil {
		t.Fatalf("enumerateFilesInDir error: %v", err)
	}

	pathSet := make(map[string]struct{}, len(files))
	for _, f := range files {
		pathSet[f] = struct{}{}
	}

	// Symlink-to-file MUST be included (dotfile manager support).
	if _, ok := pathSet[symlinkFile]; !ok {
		t.Errorf("symlink-to-file %q must be included in backup — dotfile manager support; got: %v", symlinkFile, files)
	}

	// Regular file must also be included.
	if _, ok := pathSet[regularFile]; !ok {
		t.Errorf("regular file %q must be included; got: %v", regularFile, files)
	}
}

// TestEnumerateFilesInDir_BrokenSymlinkIsSkipped verifies that broken symlinks
// (pointing to a nonexistent target) are silently skipped without error.
func TestEnumerateFilesInDir_BrokenSymlinkIsSkipped(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Symlink requires elevated privileges on Windows")
	}

	home := t.TempDir()

	// Create a broken symlink.
	brokenLink := filepath.Join(home, "broken-link")
	if err := os.Symlink("/nonexistent/target/file.md", brokenLink); err != nil {
		t.Fatalf("os.Symlink: %v", err)
	}

	// Create a real file too.
	realFile := filepath.Join(home, "real.txt")
	if err := os.WriteFile(realFile, []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	files, err := enumerateFilesInDir(home, nil)
	if err != nil {
		t.Fatalf("enumerateFilesInDir error: %v", err)
	}

	pathSet := make(map[string]struct{}, len(files))
	for _, f := range files {
		pathSet[f] = struct{}{}
	}

	if _, ok := pathSet[brokenLink]; ok {
		t.Errorf("broken symlink %q must NOT be included in results", brokenLink)
	}
	if _, ok := pathSet[realFile]; !ok {
		t.Errorf("real file %q must be included; got: %v", realFile, files)
	}
}

// TestEnumerateFilesInDir_SymlinkInSubdirDoesNotBreakBackup verifies that when a
// subdirectory contains a symlink to another directory (common with skill directories
// junctioned on Windows), enumerateFilesInDir succeeds and the real files are found.
func TestEnumerateFilesInDir_SymlinkInSubdirDoesNotBreakBackup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Symlink requires elevated privileges on Windows — skipping symlink traversal test")
	}

	home := t.TempDir()

	// .claude/ directory with a real config file and a junctioned skills dir.
	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll .claude: %v", err)
	}

	// Real config file inside .claude.
	realConf := filepath.Join(claudeDir, "CLAUDE.md")
	if err := os.WriteFile(realConf, []byte("# Claude rules"), 0o644); err != nil {
		t.Fatalf("WriteFile realConf: %v", err)
	}

	// A skills directory outside .claude that will be symlinked in.
	skillsTarget := filepath.Join(home, "skills-real")
	if err := os.MkdirAll(skillsTarget, 0o755); err != nil {
		t.Fatalf("MkdirAll skillsTarget: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillsTarget, "go-testing.md"), []byte("# Go testing"), 0o644); err != nil {
		t.Fatalf("WriteFile skill: %v", err)
	}

	// Junction-like symlink inside .claude pointing to skills-real.
	junctionLink := filepath.Join(claudeDir, "skills")
	if err := os.Symlink(skillsTarget, junctionLink); err != nil {
		t.Fatalf("os.Symlink: %v", err)
	}

	// enumerateFilesInDir on .claude must not error.
	files, err := enumerateFilesInDir(claudeDir, nil)
	if err != nil {
		t.Fatalf("enumerateFilesInDir(%q) returned error with symlink inside: %v", claudeDir, err)
	}

	// realConf must be in the list.
	realConfFound := false
	for _, f := range files {
		if f == realConf {
			realConfFound = true
		}
		// The junction link itself must NOT be in the list — it's a symlink to a directory.
		// If included, snapshotPath will receive it and os.Stat resolves to a dir → error.
		if f == junctionLink {
			t.Errorf("enumerateFilesInDir included junction/symlink path %q — symlinks must be skipped", f)
		}
		// Must NOT contain anything from inside the symlink target directory.
		linkedSkillFile := filepath.Join(junctionLink, "go-testing.md")
		if f == linkedSkillFile {
			t.Errorf("enumerateFilesInDir followed symlink and returned %q — must skip junctions/symlinks", f)
		}
	}

	if !realConfFound {
		t.Errorf("enumerateFilesInDir missing real conf file %q; got: %v", realConf, files)
	}
}
