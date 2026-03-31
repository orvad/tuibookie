# go-ssh

A terminal-based SSH bookmark manager built with Go. Organize your SSH connections into categories, browse them with an interactive TUI, and connect with a single keypress.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Huh](https://github.com/charmbracelet/huh), and [Lip Gloss](https://github.com/charmbracelet/lipgloss) from the Charm ecosystem.

## Features

- **Interactive TUI** — Navigate bookmarks and categories with arrow keys
- **Category management** — Add, rename, and delete categories
- **Bookmark management** — Add, edit, and delete SSH bookmarks within categories
- **Instant connect** — Select a bookmark and connect via SSH immediately
- **Import/Export** — Back up your bookmarks to JSON and import from backup files
- **Alphabetical sorting** — Categories and bookmarks are always sorted (case-insensitive)
- **Configurable storage** — Choose where your bookmarks file lives

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.26 or later

### Build from source

```bash
git clone <repo-url>
cd go-ssh
go build -o go-ssh .
```

This produces a `go-ssh` binary in the current directory. Move it to a directory in your `$PATH` to use it from anywhere:

```bash
# Example: move to /usr/local/bin
sudo mv go-ssh /usr/local/bin/

# Or to a user-local bin directory
mv go-ssh ~/.local/bin/
```

### Cross-compilation

Build for a different platform:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o go-ssh .

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o go-ssh .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o go-ssh .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o go-ssh .
```

## Usage

```bash
go-ssh
```

### Navigation

The app uses a stack-based navigation model. Use arrow keys or vim-style keys to move around:

#### Category list (root view)

| Key | Action |
|---|---|
| `Up` / `k` | Move cursor up |
| `Down` / `j` | Move cursor down |
| `Enter` / `Right` / `l` | Open selected category |
| `a` | Add a new category |
| `e` | Rename selected category |
| `d` | Delete selected category |
| `s` | Open settings (import/export) |
| `q` / `Esc` | Quit |

#### Bookmark list (inside a category)

| Key | Action |
|---|---|
| `Up` / `k` | Move cursor up |
| `Down` / `j` | Move cursor down |
| `Enter` | Connect to selected bookmark via SSH |
| `a` | Add a new bookmark |
| `e` | Edit selected bookmark |
| `d` | Delete selected bookmark |
| `Left` / `Esc` / `h` | Go back to categories |
| `q` | Quit |

#### Settings view

| Key | Action |
|---|---|
| `Up` / `k` | Move cursor up |
| `Down` / `j` | Move cursor down |
| `Enter` / `Right` / `l` | Execute selected action |
| `Left` / `Esc` / `h` | Go back to categories |
| `q` | Quit |

Settings provides:

- **Export bookmarks** — Saves a backup to the current working directory as `bookmarks-backup-YYYY-MM-DD-HHMMSS.json`
- **Import bookmarks** — Lists `.json` files in the current directory to choose from, or lets you enter a file path manually. Imported bookmarks are merged into existing categories.

#### Forms (add/edit)

| Key | Action |
|---|---|
| `Enter` | Submit the form |
| `Esc` | Cancel and go back |
| `Tab` | Next field (multi-field forms) |

## Configuration

### Bookmarks file location

By default, bookmarks are stored at:

```
~/.config/go-ssh/bookmarks.json
```

Override this with:

```bash
# CLI flag (highest priority)
go-ssh --config /path/to/bookmarks.json

# Environment variable
export GO_SSH_CONFIG=/path/to/bookmarks.json
go-ssh
```

Priority order: `--config` flag > `GO_SSH_CONFIG` env var > default path.

The config directory and file are created automatically on first run.

### Bookmarks file format

The bookmarks file is plain JSON. Each key is a category name, and the value is an array of bookmarks with `name` and `cmd` fields:

```json
{
  "production": [
    {
      "cmd": "ssh deploy@10.0.1.50",
      "name": "deploy"
    },
    {
      "cmd": "ssh root@10.0.1.50 -p 2222",
      "name": "root (custom port)"
    }
  ],
  "staging": [
    {
      "cmd": "ssh dev@staging.example.com",
      "name": "dev"
    }
  ]
}
```

The `cmd` field can be any valid SSH command, including flags like `-p` for custom ports, `-i` for identity files, etc.

You can edit this file manually — the app will pick up changes on next launch.

## Project structure

```
go-ssh/
  main.go                       Entry point, flag parsing, TUI launch
  internal/
    bookmark/
      bookmark.go               Bookmark types, CRUD, import/export
      bookmark_test.go           Tests for bookmark operations
    config/
      config.go                  Config path resolution
      config_test.go             Tests for config resolution
    tui/
      model.go                   Bubble Tea model, routing, state
      styles.go                  Lip Gloss style definitions (Monokai theme)
      category.go                Category list view and key handling
      bookmarks.go               Bookmark list view and key handling
      forms.go                   Huh form handling (add/edit/import)
      settings.go                Settings view (import/export)
      exec.go                    SSH execution message types
```

## License

MIT
