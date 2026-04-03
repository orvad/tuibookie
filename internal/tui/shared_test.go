package tui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
)

// localBM returns a simple local bookmarks map with the given category names.
func localBM(cats ...string) bookmark.Bookmarks {
	bm := bookmark.Bookmarks{}
	for _, c := range cats {
		bm[c] = []bookmark.Bookmark{{Name: "bm", Cmd: "echo hi"}}
	}
	return bm
}

// sharedBM returns a simple shared bookmarks map with the given category names.
func sharedBM(cats ...string) bookmark.Bookmarks {
	bm := bookmark.Bookmarks{}
	for _, c := range cats {
		bm[c] = []bookmark.Bookmark{{Name: "sbm", Cmd: "echo shared"}}
	}
	return bm
}

// baseModel returns a minimal Model with sensible defaults for testing.
func baseModel() Model {
	ApplyTheme(false)
	return Model{
		currentView: categoryView,
		width:       80,
	}
}

// ─── 1. Section headers only when both sections have data ────────────────────

func TestViewCategory_OnlyLocal_NoHeaders(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("alpha")
	m.categories = bookmark.Categories(m.bookmarks)

	out := m.viewCategory()
	if strings.Contains(out, "LOCAL") || strings.Contains(out, "SHARED") {
		t.Errorf("expected no section headers with only local bookmarks, got:\n%s", out)
	}
}

func TestViewCategory_OnlyShared_NoHeaders(t *testing.T) {
	m := baseModel()
	m.sharedBookmarks = sharedBM("beta")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)

	out := m.viewCategory()
	if strings.Contains(out, "LOCAL") || strings.Contains(out, "SHARED") {
		t.Errorf("expected no section headers with only shared bookmarks, got:\n%s", out)
	}
}

func TestViewCategory_BothSections_ShowsHeaders(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("alpha")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("beta")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)

	out := m.viewCategory()
	if !strings.Contains(out, "LOCAL") {
		t.Errorf("expected LOCAL header, got:\n%s", out)
	}
	if !strings.Contains(out, "SHARED") {
		t.Errorf("expected SHARED header, got:\n%s", out)
	}
}

// ─── 2. Read-only header ─────────────────────────────────────────────────────

func TestViewCategory_ReadOnly_ShowsReadOnly(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("alpha")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("beta")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
	m.sharedReadOnly = true

	out := m.viewCategory()
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "read-only") {
		t.Errorf("expected 'read-only' in output when sharedReadOnly=true, got:\n%s", out)
	}
}

// ─── 3. Unified cursor spans both sections ────────────────────────────────────

func TestCategoryAtCursor(t *testing.T) {
	m := baseModel()
	// Use two local and two shared categories with predictable sort order.
	m.bookmarks = localBM("aaa", "bbb")
	m.categories = bookmark.Categories(m.bookmarks) // ["aaa", "bbb"]
	m.sharedBookmarks = sharedBM("ccc", "ddd")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks) // ["ccc", "ddd"]

	tests := []struct {
		cursor     int
		wantName   string
		wantShared bool
	}{
		{0, "aaa", false},
		{1, "bbb", false},
		{2, "ccc", true},
		{3, "ddd", true},
	}

	for _, tt := range tests {
		m.catCursor = tt.cursor
		name, shared := m.categoryAtCursor()
		if name != tt.wantName || shared != tt.wantShared {
			t.Errorf("cursor=%d: got (%q, %v), want (%q, %v)",
				tt.cursor, name, shared, tt.wantName, tt.wantShared)
		}
	}
}

// ─── 4. totalCategoryItems ───────────────────────────────────────────────────

func TestTotalCategoryItems(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("a", "b")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("x", "y", "z")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)

	if got := m.totalCategoryItems(); got != 5 {
		t.Errorf("totalCategoryItems() = %d, want 5", got)
	}
}

// ─── 5. Read-only gating in category view ────────────────────────────────────

func readOnlyModel() Model {
	m := baseModel()
	m.bookmarks = localBM("local-cat")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("shared-cat")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
	m.sharedReadOnly = true
	m.isSharedContext = true
	m.catCursor = 1 // pointing at the shared category
	return m
}

func sendKey(m Model, key string) Model {
	msg := tea.KeyPressMsg{Code: -1, Text: key}
	result, _ := m.Update(msg)
	return result.(Model)
}

func TestReadOnly_AddKey_SetsStatus(t *testing.T) {
	m := readOnlyModel()
	m2 := sendKey(m, "a")
	if m2.statusMsg != "Shared bookmarks are read-only" {
		t.Errorf("expected read-only status for 'a', got %q", m2.statusMsg)
	}
}

func TestReadOnly_EditKey_SetsStatus(t *testing.T) {
	m := readOnlyModel()
	m2 := sendKey(m, "e")
	if m2.statusMsg != "Shared bookmarks are read-only" {
		t.Errorf("expected read-only status for 'e', got %q", m2.statusMsg)
	}
}

func TestReadOnly_DeleteKey_SetsStatus(t *testing.T) {
	m := readOnlyModel()
	m2 := sendKey(m, "d")
	if m2.statusMsg != "Shared bookmarks are read-only" {
		t.Errorf("expected read-only status for 'd', got %q", m2.statusMsg)
	}
}

// ─── 6. Breadcrumb in bookmark view ──────────────────────────────────────────

func TestViewBookmark_LocalContext_Breadcrumb(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("mycat")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("sharedcat")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
	m.selectedCat = "mycat"
	m.isSharedContext = false
	m.currentView = bookmarkView

	out := m.viewBookmark()
	if !strings.Contains(out, "LOCAL") {
		t.Errorf("expected LOCAL in breadcrumb for local context, got:\n%s", out)
	}
	if !strings.Contains(out, "›") {
		t.Errorf("expected '›' separator in breadcrumb for local context, got:\n%s", out)
	}
}

func TestViewBookmark_SharedContext_Breadcrumb(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("mycat")
	m.categories = bookmark.Categories(m.bookmarks)
	m.sharedBookmarks = sharedBM("sharedcat")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
	m.selectedCat = "sharedcat"
	m.isSharedContext = true
	m.currentView = bookmarkView

	out := m.viewBookmark()
	if !strings.Contains(out, "SHARED") {
		t.Errorf("expected SHARED in breadcrumb for shared context, got:\n%s", out)
	}
	if !strings.Contains(out, "›") {
		t.Errorf("expected '›' separator in breadcrumb for shared context, got:\n%s", out)
	}
}

func TestViewBookmark_OnlyLocal_NoBreadcrumb(t *testing.T) {
	m := baseModel()
	m.bookmarks = localBM("mycat")
	m.categories = bookmark.Categories(m.bookmarks)
	m.selectedCat = "mycat"
	m.isSharedContext = false
	m.currentView = bookmarkView

	out := m.viewBookmark()
	if strings.Contains(out, "›") {
		t.Errorf("expected no '›' separator when only local bookmarks exist, got:\n%s", out)
	}
}

// ─── 7. Sync error clears stale shared data ───────────────────────────────────

func TestSyncError_ClearsSharedData(t *testing.T) {
	m := baseModel()
	m.sharedBookmarks = sharedBM("stale")
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
	m.configDir = t.TempDir()

	msg := sharedSyncMsg{err: errors.New("fail")}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.sharedBookmarks != nil {
		t.Errorf("expected sharedBookmarks to be nil after sync error, got %v", m2.sharedBookmarks)
	}
	if m2.sharedCategories != nil {
		t.Errorf("expected sharedCategories to be nil after sync error, got %v", m2.sharedCategories)
	}
	if !m2.statusIsError {
		t.Errorf("expected statusIsError=true after sync error")
	}
}

// ─── 8. Successful sync populates shared data ────────────────────────────────

func TestSyncSuccess_PopulatesSharedData(t *testing.T) {
	m := baseModel()
	m.configDir = t.TempDir()

	bm := sharedBM("newcat")
	msg := sharedSyncMsg{bookmarks: bm, readOnly: true}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if len(m2.sharedBookmarks) == 0 {
		t.Errorf("expected sharedBookmarks to be populated after successful sync")
	}
	if !m2.sharedReadOnly {
		t.Errorf("expected sharedReadOnly=true after sync with readOnly=true")
	}
	if m2.statusMsg != "Synced" {
		t.Errorf("expected statusMsg 'Synced', got %q", m2.statusMsg)
	}
}

// ─── 9. currentBookmarkItems routing ─────────────────────────────────────────

func TestCurrentBookmarkItems_LocalContext(t *testing.T) {
	m := baseModel()
	m.bookmarks = bookmark.Bookmarks{
		"mycat": {{Name: "local-bm", Cmd: "echo local"}},
	}
	m.sharedBookmarks = bookmark.Bookmarks{
		"mycat": {{Name: "shared-bm", Cmd: "echo shared"}},
	}
	m.selectedCat = "mycat"
	m.isSharedContext = false

	items := m.currentBookmarkItems()
	if len(items) != 1 || items[0].Name != "local-bm" {
		t.Errorf("expected local bookmark, got %v", items)
	}
}

func TestCurrentBookmarkItems_SharedContext(t *testing.T) {
	m := baseModel()
	m.bookmarks = bookmark.Bookmarks{
		"mycat": {{Name: "local-bm", Cmd: "echo local"}},
	}
	m.sharedBookmarks = bookmark.Bookmarks{
		"mycat": {{Name: "shared-bm", Cmd: "echo shared"}},
	}
	m.selectedCat = "mycat"
	m.isSharedContext = true

	items := m.currentBookmarkItems()
	if len(items) != 1 || items[0].Name != "shared-bm" {
		t.Errorf("expected shared bookmark, got %v", items)
	}
}
