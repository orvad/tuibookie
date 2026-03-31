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
	expected := filepath.Join(home, ".config", "tuibookie", "bookmarks.json")
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

func TestLoadAppConfigMissing(t *testing.T) {
	dir := t.TempDir()
	cfg, err := LoadAppConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BookmarksPath != "" {
		t.Fatalf("expected empty BookmarksPath, got %s", cfg.BookmarksPath)
	}
}

func TestLoadAppConfigValid(t *testing.T) {
	dir := t.TempDir()
	data := []byte(`{"bookmarks_path": "/tmp/my-bookmarks.json"}`)
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadAppConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BookmarksPath != "/tmp/my-bookmarks.json" {
		t.Fatalf("expected /tmp/my-bookmarks.json, got %s", cfg.BookmarksPath)
	}
}

func TestLoadAppConfigEmpty(t *testing.T) {
	dir := t.TempDir()
	data := []byte(`{}`)
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadAppConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BookmarksPath != "" {
		t.Fatalf("expected empty BookmarksPath, got %s", cfg.BookmarksPath)
	}
}
