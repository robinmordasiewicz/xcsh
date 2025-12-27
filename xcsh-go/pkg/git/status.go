// Package git provides utilities for detecting git repository status.
package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoInfo contains information about the current git repository state.
type RepoInfo struct {
	InRepo    bool   // Whether we're inside a git repository
	RepoName  string // Repository name (directory name of repo root)
	Branch    string // Current branch name
	IsDirty   bool   // Has uncommitted changes (staged or unstaged)
	Ahead     int    // Commits ahead of remote
	Behind    int    // Commits behind remote
	HasRemote bool   // Whether the branch tracks a remote
}

// GetRepoInfo returns information about the current git repository.
// If not in a git repository, InRepo will be false and other fields empty/zero.
func GetRepoInfo() RepoInfo {
	info := RepoInfo{}

	// Check if we're in a git repository
	repoRoot, err := getRepoRoot()
	if err != nil {
		return info
	}

	info.InRepo = true
	info.RepoName = filepath.Base(repoRoot)

	// Get current branch
	info.Branch = getBranch()

	// Check for uncommitted changes
	info.IsDirty = isDirty()

	// Get ahead/behind counts
	info.Ahead, info.Behind, info.HasRemote = getAheadBehind()

	return info
}

// getRepoRoot returns the root directory of the git repository.
func getRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getBranch returns the current branch name.
func getBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// isDirty returns true if there are uncommitted changes (staged or unstaged).
func isDirty() bool {
	// git status --porcelain returns nothing if clean
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// getAheadBehind returns the number of commits ahead and behind the remote.
// Returns (0, 0, false) if there's no remote tracking branch.
func getAheadBehind() (ahead int, behind int, hasRemote bool) {
	// Get the upstream tracking branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	if _, err := cmd.Output(); err != nil {
		// No upstream configured
		return 0, 0, false
	}

	hasRemote = true

	// Get ahead count
	cmd = exec.Command("git", "rev-list", "--count", "@{upstream}..HEAD")
	if output, err := cmd.Output(); err == nil {
		if n, err := parseCount(strings.TrimSpace(string(output))); err == nil {
			ahead = n
		}
	}

	// Get behind count
	cmd = exec.Command("git", "rev-list", "--count", "HEAD..@{upstream}")
	if output, err := cmd.Output(); err == nil {
		if n, err := parseCount(strings.TrimSpace(string(output))); err == nil {
			behind = n
		}
	}

	return ahead, behind, hasRemote
}

// parseCount parses a string to an integer, returning 0 on error.
func parseCount(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	var n int
	// Simple manual parsing for positive integers
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// StatusIcon returns an appropriate icon for the git status.
// Returns: "✓" (clean), "*" (dirty), "↑" (ahead), "↓" (behind), "↕" (diverged), or "" (no repo).
func (r RepoInfo) StatusIcon() string {
	if !r.InRepo {
		return ""
	}

	// Priority: dirty > diverged > ahead/behind > clean
	if r.IsDirty {
		return "*"
	}

	if r.Ahead > 0 && r.Behind > 0 {
		return "↕"
	}

	if r.Ahead > 0 {
		return "↑"
	}

	if r.Behind > 0 {
		return "↓"
	}

	return "✓"
}

// FormatStatus returns a formatted status string like "main ✓" or "feature *".
func (r RepoInfo) FormatStatus() string {
	if !r.InRepo {
		return ""
	}

	icon := r.StatusIcon()
	if icon == "" {
		return r.Branch
	}
	return r.Branch + " " + icon
}
