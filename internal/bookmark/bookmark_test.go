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
