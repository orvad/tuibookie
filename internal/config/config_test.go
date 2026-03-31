package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathFlag(t *testing.T) {
	path := ResolvePath("/tmp/custom.json", "", t.TempDir())
	if path != "/tmp/custom.json" {
		t.Fatalf("expected /tmp/custom.json, got %s", path)
	}
}

func TestResolvePathEnv(t *testing.T) {
	path := ResolvePath("", "/tmp/env.json", t.TempDir())
	if path != "/tmp/env.json" {
		t.Fatalf("expected /tmp/env.json, got %s", path)
	}
}

func TestResolvePathFlagOverridesEnv(t *testing.T) {
	path := ResolvePath("/tmp/flag.json", "/tmp/env.json", t.TempDir())
	if path != "/tmp/flag.json" {
		t.Fatalf("expected flag to override env, got %s", path)
	}
}

func TestResolvePathDefault(t *testing.T) {
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "tuibookie", "bookmarks.json")
	path := ResolvePath("", "", t.TempDir())
	if path != expected {
		t.Fatalf("expected %s, got %s", expected, path)
	}
}

func TestResolvePathFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := AppConfig{BookmarksPath: "/tmp/from-config.json"}
	if err := SaveAppConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	path := ResolvePath("", "", dir)
	if path != "/tmp/from-config.json" {
		t.Fatalf("expected /tmp/from-config.json, got %s", path)
	}
}

func TestResolvePathFlagOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := AppConfig{BookmarksPath: "/tmp/from-config.json"}
	if err := SaveAppConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	path := ResolvePath("/tmp/flag.json", "", dir)
	if path != "/tmp/flag.json" {
		t.Fatalf("expected flag to override config, got %s", path)
	}
}

func TestResolvePathEnvOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := AppConfig{BookmarksPath: "/tmp/from-config.json"}
	if err := SaveAppConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	path := ResolvePath("", "/tmp/env.json", dir)
	if path != "/tmp/env.json" {
		t.Fatalf("expected env to override config, got %s", path)
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

func TestSaveAppConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := AppConfig{BookmarksPath: "/tmp/custom.json"}
	if err := SaveAppConfig(dir, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := LoadAppConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}
	if loaded.BookmarksPath != "/tmp/custom.json" {
		t.Fatalf("expected /tmp/custom.json, got %s", loaded.BookmarksPath)
	}
}

func TestSaveAppConfigCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	cfg := AppConfig{BookmarksPath: "/tmp/test.json"}
	if err := SaveAppConfig(dir, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.json")); err != nil {
		t.Fatalf("config.json not created: %v", err)
	}
}
