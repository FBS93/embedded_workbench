// Package project provides utilities for detecting and normalizing project names.
//
// It replicates the detection logic from the Claude Code shell helpers and
// OpenCode TypeScript plugin in pure Go, so CLI and MCP server can share
// a single canonical implementation.
package project

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DetectProject detects the project name for a given directory.
// Priority: git remote origin repo name → git root basename → dir basename.
// The returned name is always non-empty and already normalized (lowercase, trimmed).
func DetectProject(dir string) string {
	// Guard empty dir — nothing useful to detect.
	if dir == "" {
		return "unknown"
	}
	// Guard against arg injection: a dir starting with "-" would be
	// interpreted as a git flag when passed to `git -C <dir>`.
	if strings.HasPrefix(dir, "-") {
		dir = "./" + dir
	}

	if name := detectFromGitRemote(dir); name != "" {
		return normalize(name)
	}
	if name := detectFromGitRoot(dir); name != "" {
		return normalize(name)
	}
	base := filepath.Base(dir)
	if base == "" || base == "." {
		return "unknown"
	}
	return normalize(base)
}

// normalize applies canonical project name rules: lowercase + trim whitespace.
// It mirrors the normalization applied by the store layer so that DetectProject
// always returns a value that is consistent with stored project names.
func normalize(name string) string {
	n := strings.TrimSpace(strings.ToLower(name))
	if n == "" {
		return "unknown"
	}
	return n
}

// detectFromGitRemote attempts to determine the project name from the git
// remote "origin" URL. Returns empty string if git is unavailable, the
// directory is not a repo, or there is no origin remote.
func detectFromGitRemote(dir string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "-C", dir, "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(string(out))
	return extractRepoName(url)
}

// detectFromGitRoot returns the basename of the git repository root.
// Falls back to empty string when git is unavailable or the directory is not
// inside a git repository.
func detectFromGitRoot(dir string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return ""
	}
	return filepath.Base(root)
}

// extractRepoName parses a git remote URL and returns just the repository name.
//
// Supported URL formats:
//   - SSH:   git@github.com:user/repo.git
//   - HTTPS: https://github.com/user/repo.git
//   - Either with or without the trailing .git suffix
func extractRepoName(url string) string {
	// Strip trailing .git suffix
	url = strings.TrimSuffix(url, ".git")

	// Split on both "/" and ":" to handle SSH and HTTPS uniformly
	parts := strings.FieldsFunc(url, func(r rune) bool {
		return r == '/' || r == ':'
	})
	if len(parts) == 0 {
		return ""
	}
	name := parts[len(parts)-1]
	return strings.TrimSpace(name)
}
