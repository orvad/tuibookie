# TuiBookie

A fast, interactive terminal bookmark manager for CLI commands.

Organize your frequently used commands into categories, browse them with an intuitive Terminal User Interface, and execute with a single keypress.

<div class="terminal">
  <div class="terminal-bar">
    <div class="terminal-dot red"></div>
    <div class="terminal-dot yellow"></div>
    <div class="terminal-dot green"></div>
  </div>
  <div class="terminal-body">
    <div class="tui">
      <div class="tui-badge">
        <span class="tui-title-prefix">Tui</span><span class="tui-title-name">Bookie</span>
        <span class="tui-title-ver">v1.7.1</span>
      </div>
      <div class="tui-sep">━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━</div>
      <div class="tui-selected">&gt; Deployments (3)</div>
      <div>&nbsp; Docker (4)</div>
      <div>&nbsp; Kubernetes (3)</div>
      <div>&nbsp; FFmpeg (9)</div>
      <div>&nbsp; Production servers (3)</div>
      <div>&nbsp; Staging servers (2)</div>
      <div class="tui-help">
        <span>[a]</span>dd <span>[e]</span>dit <span>[d]</span>elete
        <span>[s]</span>ettings <span>[enter/→]</span> open
        <span>[q/esc]</span> quit
      </div>
    </div>
  </div>
</div>

## Why TuiBookie?

Tired of scrolling through shell history to find that one SSH command? TuiBookie was born from a simple frustration: **too many servers, too many commands, no good way to organize them in the terminal.**

Built originally as an **SSH bookmark manager** (called *go-ssh*), TuiBookie grew into a general-purpose CLI command organizer. Save any command -- SSH connections, Docker workflows, kubectl operations, deployment scripts -- and run them instantly from a clean, keyboard-driven interface.

Ever wanted to share complex CLI commands with your team? With TuiBookie you can -- sync shared bookmarks through any git repo and keep everyone on the same page.

**No browser. No GUI. Just your terminal, organized.**

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

TuiBookie is written in [Go](https://go.dev/) and compiles to a **single, self-contained binary** -- no runtime, no dependencies, no package manager required. Download it, put it in your PATH, and it just works. This makes it easy to install on remote servers, air-gapped machines, or anywhere you don't want to manage a toolchain.

The terminal interface is built with the [Charm](https://charm.sh/) ecosystem:

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
