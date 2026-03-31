# go-ssh: Terminal SSH Bookmark Manager

## Overview

A Go TUI application for managing and connecting to SSH server bookmarks. Bookmarks are organized by category and stored in a JSON file. The app provides an interactive terminal interface built with bubbletea/huh for browsing, managing, and connecting to bookmarks.

## Data Model

Bookmarks are stored as a JSON map of category name to array of bookmark objects:

```json
{
  "d-services.tvmidtvest.dk": [
    { "cmd": "ssh d-services@10.21.10.160", "name": "d-services" },
    { "cmd": "ssh maor@10.21.10.160", "name": "maor" }
  ],
  "production-servers": [
    { "cmd": "ssh deploy@10.21.20.50", "name": "deploy" }
  ]
}
```

### Types

```go
type Bookmark struct {
    Cmd  string `json:"cmd"`
    Name string `json:"name"`
}

type Bookmarks map[string][]Bookmark
```

## Config Path Resolution

Resolved in this order:

1. `--config` CLI flag
2. `GO_SSH_CONFIG` environment variable
3. `~/.config/go-ssh/bookmarks.json` (default)

If the file or directory doesn't exist, it is created automatically with an empty `{}` JSON object.

## Package Structure

```
go-ssh/
  main.go                  -- Entry point, CLI flag parsing, launches TUI
  internal/
    config/
      config.go            -- Load/save JSON, resolve config path
    bookmark/
      bookmark.go          -- Bookmark/category types and CRUD operations
    model/
      model.go             -- Bubbletea model, views, navigation state
```

## Dependencies

- `github.com/charmbracelet/bubbletea` -- TUI framework
- `github.com/charmbracelet/bubbles` -- List component
- `github.com/charmbracelet/lipgloss` -- Styling
- `github.com/charmbracelet/huh` -- Forms for add/edit dialogs

## TUI Navigation

Stack-based navigation with two main views:

### Category List View

Displays all categories. Hotkeys:

- `Enter` -- drill into selected category
- `a` -- add new category (huh text input form)
- `d` -- delete selected category (with confirmation)
- `q` -- quit

### Bookmark List View

Displays bookmarks in the selected category. Hotkeys:

- `Enter` -- connect to selected bookmark (exec SSH)
- `a` -- add new bookmark (huh form: name + cmd)
- `e` -- edit selected bookmark (huh form pre-filled)
- `d` -- delete selected bookmark (with confirmation)
- `Esc` -- back to category list
- `q` -- quit

### Forms (huh)

- **Add category:** single text input for category name
- **Add bookmark:** text inputs for name and cmd
- **Edit bookmark:** pre-filled text inputs for name and cmd

### Visual Layout

```
┌─────────────────────────────┐
│  Category / Bookmark title   │
│  ─────────────────────────── │
│  > item 1                    │
│    item 2                    │
│    item 3                    │
│                              │
│  [a]dd  [e]dit  [d]elete    │
│  [esc] back  [q]uit         │
└─────────────────────────────┘
```

Styled with lipgloss for colors, borders, and a help bar at the bottom.

## SSH Execution

- On bookmark select, the bubbletea program shuts down cleanly (restores terminal state)
- The `cmd` field is split into args using `strings.Fields()`
- `syscall.Exec` replaces the current process with the SSH command
- Supports commands like `ssh user@host` and `ssh -p 2222 user@host`

## Error Handling

| Scenario | Behavior |
|---|---|
| Missing JSON file | Auto-create empty `{}` file |
| Missing config directory | Create `~/.config/go-ssh/` recursively |
| Malformed JSON | Show error message and exit |
| Invalid SSH command | Show error in TUI before exec |
| Empty bookmarks file | Show empty category list with hint to press `a` |
| Category with no bookmarks | Show empty bookmark list with hint to press `a` |
| Duplicate bookmark names | Allowed across categories |
| Delete empty category | Remove it directly |
