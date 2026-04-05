# TuiBookie

A fast, interactive terminal bookmark manager for CLI commands.

Organize your frequently used commands into categories, browse them with an intuitive Terminal User Interface, and execute with a single keypress.

## Why TuiBookie?

If you find yourself scrolling through shell history or keeping notes of complex CLI commands, TuiBookie gives you a better way. Bookmark any command, organize it into categories, and run it instantly from a clean terminal UI.

## Key Features

- **Interactive TUI** -- Navigate bookmarks and categories with arrow keys
- **Categories** -- Add, rename, and delete categories
- **Parameterized Commands** -- Define reusable parameters with `{{name}}` syntax
- **Import/Export** -- Back up your bookmarks to JSON and import from backup files
- **Gist Sync** -- Push bookmarks to a secret GitHub Gist and pull them on any machine
- **Shared Bookmarks via Git** -- Sync shared bookmarks from any git repo for collaborative work
- **Configurable storage** -- Choose where your bookmarks file lives
- **Any CLI command** -- SSH, rsync, docker, kubectl, or any command you use regularly

## Built With

Built with the [Charm](https://charm.sh/) ecosystem:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) -- TUI framework
- [Huh](https://github.com/charmbracelet/huh) -- Interactive forms
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) -- Styling

## Quick Start

```bash
# Install
curl -sL https://raw.githubusercontent.com/orvad/tuibookie/main/install.sh | sh

# Run
tuibookie
```

See the [Installation](installation.md) page for all installation methods.
