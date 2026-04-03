package gitrepo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initBareRepo creates a bare git repo with one commit containing a bookmarks.json file.
// Returns the path to the bare repo (usable as a clone URL).
func initBareRepo(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()

	work := filepath.Join(dir, "work")
	bare := filepath.Join(dir, "bare.git")

	run(t, "", "git", "init", work)
	run(t, work, "git", "config", "user.email", "test@test.com")
	run(t, work, "git", "config", "user.name", "Test")

	if err := os.WriteFile(filepath.Join(work, "bookmarks.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, work, "git", "add", "bookmarks.json")
	run(t, work, "git", "commit", "-m", "initial")
	run(t, "", "git", "clone", "--bare", work, bare)

	return bare
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}

func TestIsCloned(t *testing.T) {
	dir := t.TempDir()

	if IsCloned(dir) {
		t.Fatal("expected empty dir to not be cloned")
	}

	if IsCloned(filepath.Join(dir, "nonexistent")) {
		t.Fatal("expected nonexistent dir to not be cloned")
	}

	if err := os.MkdirAll(filepath.Join(dir, "repo", ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if !IsCloned(filepath.Join(dir, "repo")) {
		t.Fatal("expected dir with .git to be cloned")
	}
}

func TestClone(t *testing.T) {
	bare := initBareRepo(t, `{"test":[]}`)
	cloneDir := filepath.Join(t.TempDir(), "clone")

	if err := Clone(bare, cloneDir); err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	if !IsCloned(cloneDir) {
		t.Fatal("expected clone dir to be cloned")
	}

	data, err := os.ReadFile(filepath.Join(cloneDir, "bookmarks.json"))
	if err != nil {
		t.Fatalf("expected bookmarks.json in clone: %v", err)
	}
	if string(data) != `{"test":[]}` {
		t.Fatalf("unexpected content: %s", data)
	}
}

func TestCloneInvalidURL(t *testing.T) {
	cloneDir := filepath.Join(t.TempDir(), "clone")
	err := Clone("/nonexistent/repo.git", cloneDir)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestPull(t *testing.T) {
	bare := initBareRepo(t, `{"v1":[]}`)
	cloneDir := filepath.Join(t.TempDir(), "clone")

	if err := Clone(bare, cloneDir); err != nil {
		t.Fatal(err)
	}

	second := filepath.Join(t.TempDir(), "second")
	run(t, "", "git", "clone", bare, second)
	run(t, second, "git", "config", "user.email", "test@test.com")
	run(t, second, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(second, "bookmarks.json"), []byte(`{"v2":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, second, "git", "add", "bookmarks.json")
	run(t, second, "git", "commit", "-m", "update")
	run(t, second, "git", "push")

	if err := Pull(cloneDir); err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cloneDir, "bookmarks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"v2":[]}` {
		t.Fatalf("expected v2 after pull, got: %s", data)
	}
}

func TestCommitAndPush(t *testing.T) {
	bare := initBareRepo(t, `{"v1":[]}`)
	cloneDir := filepath.Join(t.TempDir(), "clone")

	if err := Clone(bare, cloneDir); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(cloneDir, "bookmarks.json"), []byte(`{"v2":[]}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CommitAndPush(cloneDir, "bookmarks.json", "update bookmarks"); err != nil {
		t.Fatalf("CommitAndPush failed: %v", err)
	}

	verify := filepath.Join(t.TempDir(), "verify")
	run(t, "", "git", "clone", bare, verify)
	data, err := os.ReadFile(filepath.Join(verify, "bookmarks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"v2":[]}` {
		t.Fatalf("expected v2 in remote, got: %s", data)
	}
}

func TestCanPush(t *testing.T) {
	bare := initBareRepo(t, `{"test":[]}`)
	cloneDir := filepath.Join(t.TempDir(), "clone")

	if err := Clone(bare, cloneDir); err != nil {
		t.Fatal(err)
	}

	canPush, err := CanPush(cloneDir)
	if err != nil {
		t.Fatalf("CanPush failed: %v", err)
	}
	if !canPush {
		t.Fatal("expected CanPush to return true for writable repo")
	}
}

func TestResetToRemote(t *testing.T) {
	bare := initBareRepo(t, `{"v1":[]}`)
	cloneDir := filepath.Join(t.TempDir(), "clone")

	if err := Clone(bare, cloneDir); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(cloneDir, "bookmarks.json"), []byte(`{"local":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, cloneDir, "git", "add", "bookmarks.json")
	run(t, cloneDir, "git", "commit", "-m", "local change")

	second := filepath.Join(t.TempDir(), "second")
	run(t, "", "git", "clone", bare, second)
	run(t, second, "git", "config", "user.email", "test@test.com")
	run(t, second, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(second, "bookmarks.json"), []byte(`{"remote":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, second, "git", "add", "bookmarks.json")
	run(t, second, "git", "commit", "-m", "remote change")
	run(t, second, "git", "push")

	if err := ResetToRemote(cloneDir); err != nil {
		t.Fatalf("ResetToRemote failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cloneDir, "bookmarks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"remote":[]}` {
		t.Fatalf("expected remote content after reset, got: %s", data)
	}
}

func TestGitInstalled(t *testing.T) {
	if !GitInstalled() {
		t.Fatal("expected git to be installed")
	}
}
