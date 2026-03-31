# go-ssh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a TUI SSH bookmark manager that lets users browse, manage, and connect to SSH servers organized by category.

**Architecture:** Bubbletea v2 app with stack-based navigation (category list → bookmark list). Huh v2 forms for CRUD operations. Bookmarks stored in a JSON file with configurable path. On bookmark selection, `syscall.Exec` replaces the process with SSH.

**Tech Stack:** Go 1.26, bubbletea v2, bubbles v2 (list), huh v2 (forms), lipgloss v2 (styling)

**Import paths (v2 charm.land vanity domain):**
```go
tea "charm.land/bubbletea/v2"
"charm.land/bubbles/v2/list"
"charm.land/bubbles/v2/key"
"charm.land/huh/v2"
"charm.land/lipgloss/v2"
```

---

### Task 1: Project Setup & Dependencies

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Initialize dependencies**

```bash
cd /Users/maor/Projects/go-ssh
go get charm.land/bubbletea/v2@latest
go get charm.land/bubbles/v2@latest
go get charm.land/huh/v2@latest
go get charm.land/lipgloss/v2@latest
```

- [ ] **Step 2: Create directory structure**

```bash
mkdir -p internal/config internal/bookmark internal/tui
```

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add charm dependencies (bubbletea, bubbles, huh, lipgloss)"
```

---

### Task 2: Bookmark Types & JSON Operations

**Files:**
- Create: `internal/bookmark/bookmark.go`
- Create: `internal/bookmark/bookmark_test.go`

- [ ] **Step 1: Write failing tests for bookmark types and operations**

Create `internal/bookmark/bookmark_test.go`:

```go
package bookmark

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")

	bm, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bm) != 0 {
		t.Fatalf("expected empty bookmarks, got %d categories", len(bm))
	}
}

func TestLoadExisting(t *testing.T) {
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
	if len(bm) != 1 {
		t.Fatalf("expected 1 category, got %d", len(bm))
	}
	if len(bm["servers"]) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bm["servers"]))
	}
	if bm["servers"][0].Cmd != "ssh user@host" {
		t.Fatalf("unexpected cmd: %s", bm["servers"][0].Cmd)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")

	bm := Bookmarks{
		"cat1": {
			{Cmd: "ssh a@b", Name: "a"},
			{Cmd: "ssh c@d", Name: "c"},
		},
	}

	if err := Save(path, bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded["cat1"]) != 2 {
		t.Fatalf("expected 2 bookmarks, got %d", len(loaded["cat1"]))
	}
	if loaded["cat1"][0].Name != "a" {
		t.Fatalf("unexpected name: %s", loaded["cat1"][0].Name)
	}
}

func TestCategories(t *testing.T) {
	bm := Bookmarks{
		"bravo":   {{Cmd: "ssh b@b", Name: "b"}},
		"alpha":   {{Cmd: "ssh a@a", Name: "a"}},
		"charlie": {{Cmd: "ssh c@c", Name: "c"}},
	}

	cats := Categories(bm)
	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}
	if cats[0] != "alpha" || cats[1] != "bravo" || cats[2] != "charlie" {
		t.Fatalf("expected sorted categories, got %v", cats)
	}
}

func TestAddCategory(t *testing.T) {
	bm := Bookmarks{}
	AddCategory(bm, "new-cat")
	if _, ok := bm["new-cat"]; !ok {
		t.Fatal("expected category to exist")
	}
	if len(bm["new-cat"]) != 0 {
		t.Fatal("expected empty bookmark list")
	}
}

func TestDeleteCategory(t *testing.T) {
	bm := Bookmarks{
		"cat1": {{Cmd: "ssh a@b", Name: "a"}},
	}
	DeleteCategory(bm, "cat1")
	if _, ok := bm["cat1"]; ok {
		t.Fatal("expected category to be deleted")
	}
}

func TestAddBookmark(t *testing.T) {
	bm := Bookmarks{"cat1": {}}
	AddBookmark(bm, "cat1", Bookmark{Cmd: "ssh x@y", Name: "x"})
	if len(bm["cat1"]) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bm["cat1"]))
	}
}

func TestDeleteBookmark(t *testing.T) {
	bm := Bookmarks{
		"cat1": {
			{Cmd: "ssh a@b", Name: "a"},
			{Cmd: "ssh c@d", Name: "c"},
		},
	}
	DeleteBookmark(bm, "cat1", 0)
	if len(bm["cat1"]) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bm["cat1"]))
	}
	if bm["cat1"][0].Name != "c" {
		t.Fatalf("expected 'c', got '%s'", bm["cat1"][0].Name)
	}
}

func TestUpdateBookmark(t *testing.T) {
	bm := Bookmarks{
		"cat1": {{Cmd: "ssh a@b", Name: "a"}},
	}
	UpdateBookmark(bm, "cat1", 0, Bookmark{Cmd: "ssh x@y", Name: "x"})
	if bm["cat1"][0].Name != "x" {
		t.Fatalf("expected 'x', got '%s'", bm["cat1"][0].Name)
	}
	if bm["cat1"][0].Cmd != "ssh x@y" {
		t.Fatalf("expected 'ssh x@y', got '%s'", bm["cat1"][0].Cmd)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/maor/Projects/go-ssh
go test ./internal/bookmark/...
```

Expected: compilation errors — types and functions not defined.

- [ ] **Step 3: Implement bookmark types and operations**

Create `internal/bookmark/bookmark.go`:

```go
package bookmark

import (
	"encoding/json"
	"os"
	"sort"
)

type Bookmark struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}

type Bookmarks map[string][]Bookmark

func Load(path string) (Bookmarks, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Bookmarks{}, nil
		}
		return nil, err
	}

	var bm Bookmarks
	if err := json.Unmarshal(data, &bm); err != nil {
		return nil, err
	}
	return bm, nil
}

func Save(path string, bm Bookmarks) error {
	data, err := json.MarshalIndent(bm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Categories(bm Bookmarks) []string {
	cats := make([]string, 0, len(bm))
	for k := range bm {
		cats = append(cats, k)
	}
	sort.Strings(cats)
	return cats
}

func AddCategory(bm Bookmarks, name string) {
	bm[name] = []Bookmark{}
}

func DeleteCategory(bm Bookmarks, name string) {
	delete(bm, name)
}

func AddBookmark(bm Bookmarks, category string, b Bookmark) {
	bm[category] = append(bm[category], b)
}

func DeleteBookmark(bm Bookmarks, category string, index int) {
	items := bm[category]
	bm[category] = append(items[:index], items[index+1:]...)
}

func UpdateBookmark(bm Bookmarks, category string, index int, b Bookmark) {
	bm[category][index] = b
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/maor/Projects/go-ssh
go test ./internal/bookmark/... -v
```

Expected: all 9 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/bookmark/
git commit -m "feat: add bookmark types and CRUD operations with tests"
```

---

### Task 3: Config Path Resolution

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests for config resolution**

Create `internal/config/config_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathFlag(t *testing.T) {
	path := ResolvePath("/tmp/custom.json", "")
	if path != "/tmp/custom.json" {
		t.Fatalf("expected /tmp/custom.json, got %s", path)
	}
}

func TestResolvePathEnv(t *testing.T) {
	path := ResolvePath("", "/tmp/env.json")
	if path != "/tmp/env.json" {
		t.Fatalf("expected /tmp/env.json, got %s", path)
	}
}

func TestResolvePathFlagOverridesEnv(t *testing.T) {
	path := ResolvePath("/tmp/flag.json", "/tmp/env.json")
	if path != "/tmp/flag.json" {
		t.Fatalf("expected flag to override env, got %s", path)
	}
}

func TestResolvePathDefault(t *testing.T) {
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "go-ssh", "bookmarks.json")
	path := ResolvePath("", "")
	if path != expected {
		t.Fatalf("expected %s, got %s", expected, path)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "bookmarks.json")
	err := EnsureConfigDir(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, "sub"))
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/config/... -v
```

Expected: compilation errors.

- [ ] **Step 3: Implement config resolution**

Create `internal/config/config.go`:

```go
package config

import (
	"os"
	"path/filepath"
)

func ResolvePath(flagPath, envPath string) string {
	if flagPath != "" {
		return flagPath
	}
	if envPath != "" {
		return envPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "go-ssh", "bookmarks.json")
}

func EnsureConfigDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/config/... -v
```

Expected: all 5 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat: add config path resolution with flag/env/default"
```

---

### Task 4: TUI Model & Category List View

**Files:**
- Create: `internal/tui/model.go`
- Create: `internal/tui/styles.go`
- Create: `internal/tui/category.go`

- [ ] **Step 1: Create styles**

Create `internal/tui/styles.go`:

```go
package tui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)
```

- [ ] **Step 2: Create the main model with view state enum**

Create `internal/tui/model.go`:

```go
package tui

import (
	"github.com/charmbracelet/huh"

	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

type view int

const (
	categoryView view = iota
	bookmarkView
	formView
)

type formAction int

const (
	formAddCategory formAction = iota
	formAddBookmark
	formEditBookmark
)

type Model struct {
	bookmarks    bookmark.Bookmarks
	configPath   string
	currentView  view
	categories   []string
	catCursor    int
	selectedCat  string
	bmCursor     int
	form         *huh.Form
	formAction   formAction
	editIndex    int
	err          error
	width        int
	height       int

	// form field bindings
	formName string
	formCmd  string
}

func NewModel(bm bookmark.Bookmarks, configPath string) Model {
	cats := bookmark.Categories(bm)
	return Model{
		bookmarks:   bm,
		configPath:  configPath,
		currentView: categoryView,
		categories:  cats,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	switch m.currentView {
	case categoryView:
		return m.updateCategory(msg)
	case bookmarkView:
		return m.updateBookmark(msg)
	case formView:
		return m.updateForm(msg)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case categoryView:
		return m.viewCategory()
	case bookmarkView:
		return m.viewBookmark()
	case formView:
		return m.viewForm()
	}
	return ""
}

func (m *Model) refreshCategories() {
	m.categories = bookmark.Categories(m.bookmarks)
}

func (m *Model) save() {
	m.err = bookmark.Save(m.configPath, m.bookmarks)
}
```

- [ ] **Step 3: Create category view with navigation and actions**

Create `internal/tui/category.go`:

```go
package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

func (m Model) updateCategory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.catCursor > 0 {
				m.catCursor--
			}
		case "down", "j":
			if m.catCursor < len(m.categories)-1 {
				m.catCursor++
			}
		case "enter":
			if len(m.categories) > 0 {
				m.selectedCat = m.categories[m.catCursor]
				m.bmCursor = 0
				m.currentView = bookmarkView
			}
		case "a":
			m.formName = ""
			m.formAction = formAddCategory
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Category name").
						Value(&m.formName),
				),
			)
			m.currentView = formView
			return m, m.form.Init()
		case "d":
			if len(m.categories) > 0 {
				cat := m.categories[m.catCursor]
				bookmark.DeleteCategory(m.bookmarks, cat)
				m.refreshCategories()
				m.save()
				if m.catCursor >= len(m.categories) && m.catCursor > 0 {
					m.catCursor--
				}
			}
		}
	}
	return m, nil
}

func (m Model) viewCategory() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SSH Bookmarks"))
	b.WriteString("\n\n")

	if len(m.categories) == 0 {
		b.WriteString(dimStyle.Render("No categories yet. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, cat := range m.categories {
			count := len(m.bookmarks[cat])
			label := fmt.Sprintf("%s (%d)", cat, count)

			if i == m.catCursor {
				b.WriteString(selectedStyle.Render("> " + label))
			} else {
				b.WriteString(normalStyle.Render("  " + label))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[a]dd  [d]elete  [enter] open  [q]uit"))

	return b.String()
}
```

- [ ] **Step 4: Verify it compiles**

```bash
cd /Users/maor/Projects/go-ssh
go build ./internal/tui/...
```

Expected: compiles with no errors (updateBookmark, updateForm, viewBookmark, viewForm not yet implemented — will be added in next tasks, so we need stubs).

Add stubs to `internal/tui/model.go` temporarily — actually, we'll implement them in the next tasks. For now, add minimal stubs at the bottom of `internal/tui/model.go`:

Add to the bottom of `internal/tui/model.go`:

```go
// Stubs — implemented in subsequent tasks
func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd)     { return m, nil }
func (m Model) viewBookmark() string                             { return "" }
func (m Model) viewForm() string                                 { return "" }
```

Then verify:

```bash
go build ./internal/tui/...
```

Expected: compiles successfully.

- [ ] **Step 5: Commit**

```bash
git add internal/tui/
git commit -m "feat: add TUI model and category list view"
```

---

### Task 5: Bookmark List View

**Files:**
- Create: `internal/tui/bookmarks.go`
- Modify: `internal/tui/model.go` (remove `updateBookmark` and `viewBookmark` stubs)

- [ ] **Step 1: Remove stubs from model.go**

Remove these two lines from the bottom of `internal/tui/model.go`:

```go
func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m Model) viewBookmark() string                             { return "" }
```

- [ ] **Step 2: Create bookmark view**

Create `internal/tui/bookmarks.go`:

```go
package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

type execSSHMsg struct {
	cmd string
}

func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) {
	items := m.bookmarks[m.selectedCat]

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.currentView = categoryView
			m.bmCursor = 0
		case "up", "k":
			if m.bmCursor > 0 {
				m.bmCursor--
			}
		case "down", "j":
			if m.bmCursor < len(items)-1 {
				m.bmCursor++
			}
		case "enter":
			if len(items) > 0 {
				return m, func() tea.Msg {
					return execSSHMsg{cmd: items[m.bmCursor].Cmd}
				}
			}
		case "a":
			m.formName = ""
			m.formCmd = ""
			m.formAction = formAddBookmark
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Bookmark name").
						Value(&m.formName),
					huh.NewInput().
						Title("SSH command").
						Placeholder("ssh user@host").
						Value(&m.formCmd),
				),
			)
			m.currentView = formView
			return m, m.form.Init()
		case "e":
			if len(items) > 0 {
				bm := items[m.bmCursor]
				m.formName = bm.Name
				m.formCmd = bm.Cmd
				m.editIndex = m.bmCursor
				m.formAction = formEditBookmark
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Bookmark name").
							Value(&m.formName),
						huh.NewInput().
							Title("SSH command").
							Value(&m.formCmd),
					),
				)
				m.currentView = formView
				return m, m.form.Init()
			}
		case "d":
			if len(items) > 0 {
				bookmark.DeleteBookmark(m.bookmarks, m.selectedCat, m.bmCursor)
				m.save()
				if m.bmCursor >= len(m.bookmarks[m.selectedCat]) && m.bmCursor > 0 {
					m.bmCursor--
				}
			}
		}
	}
	return m, nil
}

func (m Model) viewBookmark() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.selectedCat))
	b.WriteString("\n\n")

	items := m.bookmarks[m.selectedCat]

	if len(items) == 0 {
		b.WriteString(dimStyle.Render("No bookmarks yet. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, bm := range items {
			label := fmt.Sprintf("%s  %s", bm.Name, dimStyle.Render(bm.Cmd))

			if i == m.bmCursor {
				b.WriteString(selectedStyle.Render("> "+bm.Name) + "  " + dimStyle.Render(bm.Cmd))
			} else {
				_ = label
				b.WriteString(normalStyle.Render("  "+bm.Name) + "  " + dimStyle.Render(bm.Cmd))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[a]dd  [e]dit  [d]elete  [enter] connect  [esc] back  [q]uit"))

	return b.String()
}
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./internal/tui/...
```

Expected: compiles successfully.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/bookmarks.go internal/tui/model.go
git commit -m "feat: add bookmark list view with navigation and actions"
```

---

### Task 6: Form Handling (Add/Edit)

**Files:**
- Create: `internal/tui/forms.go`
- Modify: `internal/tui/model.go` (remove `updateForm` and `viewForm` stubs)

- [ ] **Step 1: Remove stubs from model.go**

Remove these two lines from the bottom of `internal/tui/model.go`:

```go
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd)     { return m, nil }
func (m Model) viewForm() string                                 { return "" }
```

- [ ] **Step 2: Create form handling**

Create `internal/tui/forms.go`:

```go
package tui

import (
	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow cancelling the form
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "esc" {
			switch m.formAction {
			case formAddCategory:
				m.currentView = categoryView
			case formAddBookmark, formEditBookmark:
				m.currentView = bookmarkView
			}
			m.form = nil
			return m, nil
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		switch m.formAction {
		case formAddCategory:
			if m.formName != "" {
				bookmark.AddCategory(m.bookmarks, m.formName)
				m.refreshCategories()
				m.save()
			}
			m.currentView = categoryView

		case formAddBookmark:
			if m.formName != "" && m.formCmd != "" {
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, bookmark.Bookmark{
					Name: m.formName,
					Cmd:  m.formCmd,
				})
				m.save()
			}
			m.currentView = bookmarkView

		case formEditBookmark:
			if m.formName != "" && m.formCmd != "" {
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, bookmark.Bookmark{
					Name: m.formName,
					Cmd:  m.formCmd,
				})
				m.save()
			}
			m.currentView = bookmarkView
		}
		m.form = nil
	}

	return m, cmd
}

func (m Model) viewForm() string {
	if m.form == nil {
		return ""
	}
	return m.form.View()
}
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./internal/tui/...
```

Expected: compiles successfully.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/forms.go internal/tui/model.go
git commit -m "feat: add huh form handling for add/edit operations"
```

---

### Task 7: SSH Execution

**Files:**
- Create: `internal/tui/exec.go`

- [ ] **Step 1: Create SSH exec handler**

Create `internal/tui/exec.go`:

```go
package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	tea "charm.land/bubbletea/v2"
)

func execSSH(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	binary, err := exec.LookPath(parts[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", parts[0])
	}

	return syscall.Exec(binary, parts, os.Environ())
}

// ExecCmd returns a tea.Cmd that quits bubbletea.
// The actual exec happens after the program exits.
func ExecCmd(cmd string) tea.Cmd {
	return tea.Exec(tea.WrapExecCommand(exec.Command(strings.Fields(cmd)[0], strings.Fields(cmd)[1:]...)), func(err error) tea.Msg {
		return execDoneMsg{err: err}
	})
}

type execDoneMsg struct {
	err error
}
```

- [ ] **Step 2: Update bookmark view to use tea.Exec**

In `internal/tui/bookmarks.go`, change the `"enter"` case to:

Replace the current enter handler:

```go
		case "enter":
			if len(items) > 0 {
				return m, func() tea.Msg {
					return execSSHMsg{cmd: items[m.bmCursor].Cmd}
				}
			}
```

With:

```go
		case "enter":
			if len(items) > 0 {
				cmd := items[m.bmCursor].Cmd
				parts := strings.Fields(cmd)
				if len(parts) > 0 {
					c := exec.Command(parts[0], parts[1:]...)
					c.Stdin = os.Stdin
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					return m, tea.Exec(c, func(err error) tea.Msg {
						return execDoneMsg{err: err}
					})
				}
			}
```

Add these imports to `internal/tui/bookmarks.go`:

```go
import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)
```

Also add handling for `execDoneMsg` in `internal/tui/model.go` Update method, before the view switch:

```go
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case execDoneMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		return m, tea.Quit
	}
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./internal/tui/...
```

Expected: compiles successfully.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/
git commit -m "feat: add SSH execution via tea.Exec"
```

---

### Task 8: Main Entry Point

**Files:**
- Create: `main.go`

- [ ] **Step 1: Create main.go**

Create `main.go`:

```go
package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
	"example/go-ssh/internal/config"
	"example/go-ssh/internal/tui"
)

func main() {
	configFlag := flag.String("config", "", "path to bookmarks JSON file")
	flag.Parse()

	configPath := config.ResolvePath(*configFlag, os.Getenv("GO_SSH_CONFIG"))

	if err := config.EnsureConfigDir(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	bm, err := bookmark.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bookmarks: %v\n", err)
		os.Exit(1)
	}

	model := tui.NewModel(bm, configPath)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 2: Build the binary**

```bash
cd /Users/maor/Projects/go-ssh
go build -o go-ssh .
```

Expected: binary `go-ssh` created successfully.

- [ ] **Step 3: Smoke test**

```bash
./go-ssh
```

Expected: TUI launches showing "SSH Bookmarks" title with empty category list and hint to press `a`.

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add main entry point with config resolution and TUI launch"
```

---

### Task 9: Clean Up & Final Verification

**Files:**
- Modify: `internal/tui/exec.go` (remove unused `execSSH` function and `execSSHMsg` type from bookmarks.go)

- [ ] **Step 1: Remove unused code**

Remove from `internal/tui/exec.go` the `execSSH` function (unused since we use `tea.Exec`) and the `execSSHMsg` type from `internal/tui/bookmarks.go` (replaced by `execDoneMsg`).

The cleaned `internal/tui/exec.go` should just contain:

```go
package tui

import tea "charm.land/bubbletea/v2"

type execDoneMsg struct {
	err error
}
```

Remove the `execSSHMsg` type from `internal/tui/bookmarks.go` if it's still there.

- [ ] **Step 2: Run all tests**

```bash
cd /Users/maor/Projects/go-ssh
go test ./... -v
```

Expected: all tests pass.

- [ ] **Step 3: Build final binary**

```bash
go build -o go-ssh .
```

Expected: compiles cleanly.

- [ ] **Step 4: End-to-end manual test**

```bash
./go-ssh
```

Test the following flow:
1. App starts with empty category list
2. Press `a`, type a category name, submit → category appears
3. Press `enter` on category → bookmark list (empty)
4. Press `a`, type name and SSH command, submit → bookmark appears
5. Press `e` → edit form pre-filled
6. Press `d` → bookmark deleted
7. Press `esc` → back to categories
8. Press `d` on category → category deleted
9. Press `q` → app exits
10. Restart app → data persists from JSON file

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: clean up unused code, final verification"
```
