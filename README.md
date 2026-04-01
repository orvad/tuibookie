<p align="center">
  <img src="logo/tuibookie-logo.png" alt="TuiBookie" width="200">
</p>

<h1 align="center">TuiBookie</h1>

<p align="center">
  A fast, interactive terminal bookmark manager for CLI commands.<br>
  Organize your frequently used commands into categories, browse them with an intuitive Terminal User Interface, and execute with a single keypress.
</p>

## Features

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Huh](https://github.com/charmbracelet/huh), and [Lip Gloss](https://github.com/charmbracelet/lipgloss) from the Charm ecosystem.


- **Interactive TUI** — Navigate bookmarks and categories with arrow keys
- **Categories** — Add, rename, and delete categories
- **Import/Export** — Back up your bookmarks to JSON and import from backup files
- **Gist Sync** — Push your bookmarks to a secret GitHub Gist and pull them on any machine. Versioned backup with full revision history, powered by a Personal Access Token.
- **Configurable storage** — Choose where your bookmarks file lives
- **Any CLI command** — SSH, rsync, docker, kubectl, or any command you use regularly

## Screenshots

**Browse categories** — See all your command groups at a glance with bookmark counts.

<img src="examples/screenshot-01.png" alt="Category list" width="600">

**Browse bookmarks** — Drill into a category to see commands. Select one and press Enter to run it.

<img src="examples/screenshot-02.png" alt="Bookmark list" width="600">

**Settings** — Configure your bookmarks file path, export backups, or import from a JSON file.

<img src="examples/screenshot-03.png" alt="Settings" width="600">

## Installation

### Quick install (recommended)

Automatically downloads the latest release for your OS and architecture:

```bash
curl -sL https://raw.githubusercontent.com/orvad/tuibookie/main/install.sh | sh
```

### Homebrew (macOS and Linux)

```bash
brew tap orvad/tuibookie
brew install tuibookie
```

### Download manually

Download the latest binary from the [Releases page](https://github.com/orvad/tuibookie/releases):

| Platform | Binary |
|---|---|
| macOS (Apple Silicon) | `tuibookie-darwin-arm64` |
| macOS (Intel) | `tuibookie-darwin-amd64` |
| Linux (x86_64) | `tuibookie-linux-amd64` |
| Linux (ARM64) | `tuibookie-linux-arm64` |

Then make it executable and move it to your PATH:

```bash
chmod +x tuibookie-*
sudo mv tuibookie-* /usr/local/bin/tuibookie
```

### Build from source

Requires [Go](https://go.dev/dl/) 1.26 or later:

```bash
git clone https://github.com/orvad/tuibookie.git
cd tuibookie
go build -o tuibookie .
sudo mv tuibookie /usr/local/bin/
```

#### Cross-compilation

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o tuibookie .

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o tuibookie .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o tuibookie .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o tuibookie .
```

## Usage

```bash
tuibookie
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
| `Enter` | Run the selected command |
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

- **Bookmarks file** — View and change the path to your bookmarks JSON file. When switching, you'll see a confirmation with the number of categories and bookmarks in the target file. If the file doesn't exist, you can create a new empty one.
- **Export bookmarks** — Saves a backup to the current working directory as `bookmarks-backup-YYYY-MM-DD-HHMMSS.json`
- **Import bookmarks** — Lists `.json` files in the current directory to choose from, or lets you enter a file path manually. Imported bookmarks are merged into existing categories.
- **Push to Gist** — Uploads your bookmarks to a secret GitHub Gist. On first push, a new gist is created; subsequent pushes update it.
- **Pull from Gist** — Downloads bookmarks from your gist and replaces the local file. Shows a confirmation with category and bookmark counts before overwriting.
- **GitHub token** — Set or remove the Personal Access Token used for Gist sync. The token is stored in `config.json` and displayed masked in the UI.

#### Forms (add/edit)

| Key | Action |
|---|---|
| `Enter` | Submit the form |
| `Esc` | Cancel and go back |
| `Tab` | Next field (multi-field forms) |

## Configuration

### Config file

On first launch, TuiBookie creates a config file at:

```
~/.config/tuibookie/config.json
```

This stores your app settings, starting with the bookmarks file path. You can change the bookmarks path directly from the Settings view in the TUI — no flags needed.

### Bookmarks file location

By default, bookmarks are stored at:

```
~/.config/tuibookie/bookmarks.json
```

You can change this in the Settings view, or override with flags for scripting:

```bash
# CLI flag (highest priority)
tuibookie --config /path/to/bookmarks.json

# Environment variable
export TUIBOOKIE_CONFIG=/path/to/bookmarks.json
tuibookie
```

Priority order: `--config` flag > `TUIBOOKIE_CONFIG` env var > `config.json` setting > default path.

The config directory and files are created automatically on first run.

### Bookmarks file format

The bookmarks file is plain JSON. Each key is a category name, and the value is an array of bookmarks with `name` and `cmd` fields:

```json
{
  "servers": [
    {
      "cmd": "ssh deploy@10.0.1.50",
      "name": "deploy"
    },
    {
      "cmd": "ssh root@10.0.1.50 -p 2222",
      "name": "root (custom port)"
    }
  ],
  "docker": [
    {
      "cmd": "docker compose up -d",
      "name": "start stack"
    },
    {
      "cmd": "docker compose logs -f",
      "name": "follow logs"
    }
  ],
  "misc": [
    {
      "cmd": "rsync -avz ./dist/ user@server:/var/www/",
      "name": "deploy frontend"
    },
    {
      "cmd": "kubectl get pods -n production",
      "name": "check prod pods"
    }
  ]
}
```

The `cmd` field can be any valid shell command. You can edit this file manually — the app will pick up changes on next launch.

## License

MIT
