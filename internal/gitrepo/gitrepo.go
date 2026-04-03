package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IsCloned returns true if localPath contains a .git directory.
func IsCloned(localPath string) bool {
	info, err := os.Stat(localPath + "/.git")
	return err == nil && info.IsDir()
}

// Clone performs a shallow clone of repoURL into localPath.
func Clone(repoURL, localPath string) error {
	return runGit("", "clone", "--depth", "1", repoURL, localPath)
}

// Pull fetches and fast-forward merges the current branch.
func Pull(repoPath string) error {
	return runGit(repoPath, "pull", "--ff-only")
}

// CommitAndPush stages a file, commits with the given message, and pushes.
func CommitAndPush(repoPath, filePath, message string) error {
	if err := runGit(repoPath, "add", filePath); err != nil {
		return err
	}
	if err := runGit(repoPath, "commit", "-m", message); err != nil {
		return err
	}
	return runGit(repoPath, "push")
}

// CanPush performs a dry-run push to check write access.
// Returns true if the push would succeed, false if permission is denied.
func CanPush(repoPath string) (bool, error) {
	err := runGit(repoPath, "push", "--dry-run")
	if err == nil {
		return true, nil
	}
	// Permission denied or similar → read-only
	return false, nil
}

// ResetToRemote fetches from origin and hard-resets to the remote default branch.
func ResetToRemote(repoPath string) error {
	if err := runGit(repoPath, "fetch", "origin"); err != nil {
		return err
	}
	// Determine the default branch name
	branch, err := gitOutput(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}
	return runGit(repoPath, "reset", "--hard", "origin/"+branch)
}

// GitInstalled returns true if the git binary is available in PATH.
func GitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return nil
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}
