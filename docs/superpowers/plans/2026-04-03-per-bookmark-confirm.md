# Per-Bookmark Confirmation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a per-bookmark `Confirm` toggle that shows a "are you sure?" dialog before executing dangerous commands.

**Architecture:** Add a `Confirm bool` field to the `Bookmark` struct. Intercept both execution paths (direct and post-param) to show the existing confirm dialog when enabled. Reuse the existing `confirmView` infrastructure with a new `formConfirmExec` action.

**Tech Stack:** Go, Bubble Tea, Huh forms, Lip Gloss styling

**Spec:** `docs/superpowers/specs/2026-04-03-per-bookmark-confirm-design.md`

---

### Task 1: Add `Confirm` field to `Bookmark` struct

**Files:**
- Modify: `internal/bookmark/bookmark.go:13-16`

- [ ] **Step 1: Add the field**

In `internal/bookmark/bookmark.go`, change the `Bookmark` struct from:

```go
type Bookmark struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}
```

to:

```go
type Bookmark struct {
	Cmd     string `json:"cmd"`
	Name    string `json:"name"`
	Confirm bool   `json:"confirm,omitempty"`
}
```

- [ ] **Step 2: Run existing tests to verify no regressions**

Run: `go test ./internal/bookmark/ -v`
Expected: All existing tests pass. The `omitempty` tag means existing test JSON without `confirm` still parses correctly.

- [ ] **Step 3: Commit**

```bash
git add internal/bookmark/bookmark.go
git commit -m "feat: add Confirm field to Bookmark struct"
```

---

### Task 2: Test backward compatibility and round-trip

**Files:**
- Modify: `internal/bookmark/bookmark_test.go`

- [ ] **Step 1: Write test for loading JSON without confirm field**

Add to `internal/bookmark/bookmark_test.go`:

```go
func TestLoadBookmarkWithoutConfirmField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")

	data := `{"servers":[{"cmd":"ssh user@host","name":"myserver"}]}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	bm, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bm["servers"][0].Confirm != false {
		t.Fatal("expected Confirm to default to false")
	}
}
```

- [ ] **Step 2: Write test for round-tripping confirm field**

Add to `internal/bookmark/bookmark_test.go`:

```go
func TestSaveAndLoadWithConfirm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")

	bm := Bookmarks{
		"danger": {
			{Cmd: "rm -rf /tmp/test", Name: "cleanup", Confirm: true},
			{Cmd: "ls -la", Name: "list", Confirm: false},
		},
	}

	if err := Save(path, bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded["danger"][0].Confirm != true {
		t.Fatal("expected cleanup to have Confirm=true")
	}
	if loaded["danger"][1].Confirm != false {
		t.Fatal("expected list to have Confirm=false")
	}
}
```

- [ ] **Step 3: Write test that omitempty excludes false values from JSON**

Add to `internal/bookmark/bookmark_test.go`:

```go
func TestSaveOmitsConfirmWhenFalse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")

	bm := Bookmarks{
		"cat": {
			{Cmd: "echo hi", Name: "safe"},
		},
	}

	if err := Save(path, bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "confirm") {
		t.Fatalf("expected no confirm key in JSON when false, got: %s", data)
	}
}
```

Note: add `"strings"` to the import block at the top of the test file.

- [ ] **Step 4: Run all tests**

Run: `go test ./internal/bookmark/ -v`
Expected: All tests pass including the three new ones.

- [ ] **Step 5: Commit**

```bash
git add internal/bookmark/bookmark_test.go
git commit -m "test: add backward compat and round-trip tests for Confirm field"
```

---

### Task 3: Add `formConfirmExec` action and confirm indicator style

**Files:**
- Modify: `internal/tui/model.go:26-38`
- Modify: `internal/tui/styles.go`

- [ ] **Step 1: Add `formConfirmExec` to the `formAction` enum**

In `internal/tui/model.go`, add `formConfirmExec` at the end of the `formAction` const block:

```go
const (
	formAddCategory formAction = iota
	formEditCategory
	formAddBookmark
	formEditBookmark
	formImport
	formImportManual
	formChangeBookmarksPath
	formConfirmBookmarksPath
	formSetGistToken
	formConfirmPull
	formRunParam
	formConfirmExec
)
```

- [ ] **Step 2: Add `confirmIndicatorStyle` to styles**

In `internal/tui/styles.go`, add after the `paramStyle` definition:

```go
confirmIndicatorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#F92672")).
	Bold(true)
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/model.go internal/tui/styles.go
git commit -m "feat: add formConfirmExec action and confirm indicator style"
```

---

### Task 4: Add confirmation handler in confirm.go

**Files:**
- Modify: `internal/tui/confirm.go`

- [ ] **Step 1: Add exec import and formConfirmExec case**

In `internal/tui/confirm.go`, add `"os/exec"` and `"strings"` to the import block (keep the existing imports). The full import block becomes:

```go
import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
	"example/tuibookie/internal/config"
	"example/tuibookie/internal/gist"
)
```

Then in the `onConfirm()` method, add a `formConfirmExec` case at the end of the switch, before the closing brace:

```go
	case formConfirmExec:
		cmd := m.pendingCmd
		m.pendingCmd = ""
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			c := exec.Command(parts[0], parts[1:]...)
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return execDoneMsg{err: err}
			})
		}
		m.currentView = bookmarkView
```

- [ ] **Step 2: Update the "No" path in `resolveConfirm`**

In the same file, update `resolveConfirm()` to return to `bookmarkView` when the action is `formConfirmExec`. Change:

```go
func (m Model) resolveConfirm() (tea.Model, tea.Cmd) {
	confirmed := m.confirmCursor == 0
	if !confirmed {
		m.pendingConfigPath = ""
		m.currentView = settingsView
		return m, nil
	}
	return m.onConfirm()
}
```

to:

```go
func (m Model) resolveConfirm() (tea.Model, tea.Cmd) {
	confirmed := m.confirmCursor == 0
	if !confirmed {
		if m.confirmAction == formConfirmExec {
			m.pendingCmd = ""
			m.currentView = bookmarkView
		} else {
			m.pendingConfigPath = ""
			m.currentView = settingsView
		}
		return m, nil
	}
	return m.onConfirm()
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/confirm.go
git commit -m "feat: add formConfirmExec handler in confirm view"
```

---

### Task 5: Intercept execution in bookmarks.go (direct commands)

**Files:**
- Modify: `internal/tui/bookmarks.go`

- [ ] **Step 1: Add confirm intercept for direct execution**

In `internal/tui/bookmarks.go`, in the `updateBookmark` method, replace the direct execution block (the `case "enter"` handler, lines 32-61) with:

```go
		case "enter":
			if len(items) > 0 {
				bm := items[m.bmCursor]
				cmd := bm.Cmd
				params := bookmark.ParseParams(cmd)
				if len(params) > 0 {
					m.pendingCmd = cmd
					m.pendingParams = params
					m.formAction = formRunParam
					m.paramValues = make(map[string]*string)
					groups := make([]huh.Field, len(params))
					for i, p := range params {
						val := p.Default
						m.paramValues[p.Name] = &val
						groups[i] = huh.NewInput().
							Title(p.Name).
							Key(p.Name).
							Value(&val)
					}
					m.form = huh.NewForm(
						huh.NewGroup(groups...),
					).WithTheme(formTheme)
					m.currentView = formView
					return m, m.form.Init()
				}
				if bm.Confirm {
					m.pendingCmd = cmd
					m.confirmMsg = "Execute: " + cmd + "?"
					m.confirmAction = formConfirmExec
					m.confirmCursor = 0
					m.currentView = confirmView
					return m, nil
				}
				parts := strings.Fields(cmd)
				if len(parts) > 0 {
					c := exec.Command(parts[0], parts[1:]...)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return execDoneMsg{err: err}
					})
				}
			}
```

This inserts the confirm check between the param-form check and the direct execution.

- [ ] **Step 2: Add confirm indicator to bookmark list rendering**

In the same file, update the `viewBookmark()` method. Replace the rendering loop:

```go
	for i, bm := range items {
		if i == m.bmCursor {
			b.WriteString(selectedStyle.Render("> "+bm.Name) + "  " + renderCmd(bm.Cmd))
		} else {
			b.WriteString(normalStyle.Render("  "+bm.Name) + "  " + renderCmd(bm.Cmd))
		}
		b.WriteString("\n")
	}
```

with:

```go
	for i, bm := range items {
		indicator := ""
		if bm.Confirm {
			indicator = " " + confirmIndicatorStyle.Render("!")
		}
		if i == m.bmCursor {
			b.WriteString(selectedStyle.Render("> "+bm.Name) + indicator + "  " + renderCmd(bm.Cmd))
		} else {
			b.WriteString(normalStyle.Render("  "+bm.Name) + indicator + "  " + renderCmd(bm.Cmd))
		}
		b.WriteString("\n")
	}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/bookmarks.go
git commit -m "feat: intercept direct execution with confirm dialog and add ! indicator"
```

---

### Task 6: Intercept execution in forms.go (post-param commands)

**Files:**
- Modify: `internal/tui/forms.go`

- [ ] **Step 1: Add confirm intercept after param resolution**

In `internal/tui/forms.go`, in the `formRunParam` case (inside the `m.form.State == huh.StateCompleted` block), replace:

```go
		case formRunParam:
			values := make(map[string]string)
			for _, p := range m.pendingParams {
				values[p.Name] = m.form.GetString(p.Name)
			}
			resolved := bookmark.ResolveParams(m.pendingCmd, values)
			m.pendingCmd = ""
			m.pendingParams = nil
			m.paramValues = nil
			m.currentView = bookmarkView
			parts := strings.Fields(resolved)
			if len(parts) > 0 {
				c := exec.Command(parts[0], parts[1:]...)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return execDoneMsg{err: err}
				})
			}
```

with:

```go
		case formRunParam:
			values := make(map[string]string)
			for _, p := range m.pendingParams {
				values[p.Name] = m.form.GetString(p.Name)
			}
			resolved := bookmark.ResolveParams(m.pendingCmd, values)
			m.pendingParams = nil
			m.paramValues = nil
			items := m.bookmarks[m.selectedCat]
			if len(items) > m.bmCursor && items[m.bmCursor].Confirm {
				m.pendingCmd = resolved
				m.confirmMsg = "Execute: " + resolved + "?"
				m.confirmAction = formConfirmExec
				m.confirmCursor = 0
				m.currentView = confirmView
				m.form = nil
				return m, nil
			}
			m.pendingCmd = ""
			m.currentView = bookmarkView
			parts := strings.Fields(resolved)
			if len(parts) > 0 {
				c := exec.Command(parts[0], parts[1:]...)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return execDoneMsg{err: err}
				})
			}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add internal/tui/forms.go
git commit -m "feat: intercept post-param execution with confirm dialog"
```

---

### Task 7: Add confirm toggle to add/edit bookmark forms

**Files:**
- Modify: `internal/tui/bookmarks.go`
- Modify: `internal/tui/model.go`
- Modify: `internal/tui/forms.go`

- [ ] **Step 1: Add `editConfirm` field to Model**

In `internal/tui/model.go`, add a field to the `Model` struct after `editIndex`:

```go
	editIndex         int
	editConfirm       bool
```

- [ ] **Step 2: Update the add-bookmark form**

In `internal/tui/bookmarks.go`, replace the `case "a"` block:

```go
		case "a":
			m.formAction = formAddBookmark
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Bookmark name").
						Key("name"),
					huh.NewInput().
						Title("Command").
						Key("cmd"),
				),
			).WithTheme(formTheme)
			m.currentView = formView
			return m, m.form.Init()
```

with:

```go
		case "a":
			m.editConfirm = false
			m.formAction = formAddBookmark
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Bookmark name").
						Key("name"),
					huh.NewInput().
						Title("Command").
						Key("cmd"),
					huh.NewConfirm().
						Title("Confirm before execute?").
						Key("confirm").
						Value(&m.editConfirm),
				),
			).WithTheme(formTheme)
			m.currentView = formView
			return m, m.form.Init()
```

- [ ] **Step 3: Update the edit-bookmark form**

In the same file, replace the `case "e"` block:

```go
		case "e":
			if len(items) > 0 {
				bm := items[m.bmCursor]
				editName := bm.Name
				editCmd := bm.Cmd
				m.editIndex = m.bmCursor
				m.formAction = formEditBookmark
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Bookmark name").
							Key("name").
							Value(&editName),
						huh.NewInput().
							Title("Command").
							Key("cmd").
							Value(&editCmd),
					),
				).WithTheme(formTheme)
				m.currentView = formView
				return m, m.form.Init()
			}
```

with:

```go
		case "e":
			if len(items) > 0 {
				bm := items[m.bmCursor]
				editName := bm.Name
				editCmd := bm.Cmd
				m.editConfirm = bm.Confirm
				m.editIndex = m.bmCursor
				m.formAction = formEditBookmark
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Bookmark name").
							Key("name").
							Value(&editName),
						huh.NewInput().
							Title("Command").
							Key("cmd").
							Value(&editCmd),
						huh.NewConfirm().
							Title("Confirm before execute?").
							Key("confirm").
							Value(&m.editConfirm),
					),
				).WithTheme(formTheme)
				m.currentView = formView
				return m, m.form.Init()
			}
```

- [ ] **Step 4: Update form completion handlers to use editConfirm**

In `internal/tui/forms.go`, update the `formAddBookmark` case. Replace:

```go
		case formAddBookmark:
			if name != "" && cmd != "" {
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, bookmark.Bookmark{
					Name: name,
					Cmd:  cmd,
				})
				m.save()
			}
			m.currentView = bookmarkView
```

with:

```go
		case formAddBookmark:
			if name != "" && cmd != "" {
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, bookmark.Bookmark{
					Name:    name,
					Cmd:     cmd,
					Confirm: m.editConfirm,
				})
				m.save()
			}
			m.currentView = bookmarkView
```

Then update the `formEditBookmark` case. Replace:

```go
		case formEditBookmark:
			if name != "" && cmd != "" {
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, bookmark.Bookmark{
					Name: name,
					Cmd:  cmd,
				})
				m.save()
			}
			m.currentView = bookmarkView
```

with:

```go
		case formEditBookmark:
			if name != "" && cmd != "" {
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, bookmark.Bookmark{
					Name:    name,
					Cmd:     cmd,
					Confirm: m.editConfirm,
				})
				m.save()
			}
			m.currentView = bookmarkView
```

- [ ] **Step 5: Verify it compiles**

Run: `go build ./...`
Expected: Build succeeds.

- [ ] **Step 6: Run all tests**

Run: `go test ./... -v`
Expected: All tests pass.

- [ ] **Step 7: Commit**

```bash
git add internal/tui/bookmarks.go internal/tui/model.go internal/tui/forms.go
git commit -m "feat: add confirm toggle to add/edit bookmark forms"
```

---

### Task 8: Manual smoke test

- [ ] **Step 1: Run the app**

Run: `go run . -config /tmp/test-confirm-bookmarks.json`

- [ ] **Step 2: Test add with confirm**

1. Press `a` to add a category, name it "test"
2. Press Enter to open it
3. Press `a` to add a bookmark: name "dangerous", command "echo DANGER", set "Confirm before execute?" to Yes
4. Verify the `!` indicator appears next to "dangerous" in the list
5. Press Enter — verify the confirmation dialog appears showing "Execute: echo DANGER?"
6. Press `n` — verify it returns to the bookmark list without executing
7. Press Enter again, press `y` — verify the command executes

- [ ] **Step 3: Test add without confirm**

1. Press `a` to add another bookmark: name "safe", command "echo SAFE", leave confirm as No
2. Verify no `!` indicator
3. Press Enter — verify it executes immediately without confirmation

- [ ] **Step 4: Test edit toggle**

1. Select "safe", press `e`
2. Change "Confirm before execute?" to Yes
3. Verify `!` indicator now appears
4. Press Enter — verify confirmation dialog shows

- [ ] **Step 5: Clean up**

```bash
rm /tmp/test-confirm-bookmarks.json
```

- [ ] **Step 6: Commit (if any fixes were needed)**

Only if changes were made during smoke testing:
```bash
git add -A
git commit -m "fix: address issues found during smoke testing"
```
