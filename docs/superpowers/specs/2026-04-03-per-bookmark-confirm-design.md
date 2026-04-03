# Per-Bookmark Confirmation Before Execution

**Date:** 2026-04-03
**Status:** Approved

## Problem

All bookmarked commands execute immediately on Enter. Some commands are dangerous (e.g., `rm -rf`, `kubectl delete`, `docker system prune`) and users want a safety net before running them.

## Solution

Add a per-bookmark `Confirm` boolean field. When enabled, pressing Enter shows a confirmation dialog with the resolved command before executing. When disabled, behavior is unchanged (execute immediately).

## Design Decisions

- **Per-bookmark, not global** — some commands are dangerous, others are harmless. Users opt in per command.
- **Default: false** — preserves current behavior. Existing bookmarks and new bookmarks execute immediately unless the user enables confirmation.
- **Visual indicator** — a bold pink `!` between the bookmark name and command in the list view, so users can see at a glance which commands require confirmation.
- **Resolved command shown** — the confirmation dialog shows the command after parameter substitution (e.g., `ssh prod-01`, not `ssh {{host:prod-01}}`).
- **No migration needed** — `json:"confirm,omitempty"` means existing `bookmarks.json` files parse without changes.

## Data Model

Add `Confirm` field to `Bookmark` struct in `internal/bookmark/bookmark.go`:

```go
type Bookmark struct {
    Cmd     string `json:"cmd"`
    Name    string `json:"name"`
    Confirm bool   `json:"confirm,omitempty"`
}
```

## Visual Indicator

New style in `internal/tui/styles.go`:

```go
confirmIndicatorStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#F92672")).
    Bold(true)
```

In the bookmark list view (`internal/tui/bookmarks.go`), when `bm.Confirm` is true, render a `!` between name and command:

```
> deploy-prod  ! docker compose up -d    ← confirm enabled
  check-logs     docker logs -f app       ← normal
```

## Execution Flow

Two code paths currently execute commands:

1. **Direct execution** — `internal/tui/bookmarks.go:56-61` (commands without parameters)
2. **Post-param-form execution** — `internal/tui/forms.go:189-194` (after parameter resolution)

Both paths get the same intercept. When the selected bookmark has `Confirm: true`:

1. Store the resolved command string in `m.pendingCmd`
2. Set `m.confirmMsg` to `"Execute: <resolved command>?"`
3. Set `m.confirmAction` to `formConfirmExec` (new action)
4. Switch to `confirmView`

In `internal/tui/confirm.go`, add a `formConfirmExec` case to `onConfirm()`:
- **Yes:** execute the command via `exec.Command` + `tea.ExecProcess`
- **No:** clear `pendingCmd` and return to `bookmarkView`

## Add/Edit Bookmark Forms

Both the add and edit bookmark forms (`internal/tui/bookmarks.go`) gain a third field:

```go
huh.NewConfirm().
    Title("Confirm before execute?").
    Key("confirm").
    Value(&confirmVal)
```

- **Add form:** `confirmVal` defaults to `false`
- **Edit form:** `confirmVal` pre-populated with `bm.Confirm`

The form completion handlers pass the value through to `AddBookmark` / `UpdateBookmark`.

## Files Changed

| File | Change |
|------|--------|
| `internal/bookmark/bookmark.go` | Add `Confirm bool` field to `Bookmark` struct |
| `internal/tui/styles.go` | Add `confirmIndicatorStyle` (bold pink) |
| `internal/tui/bookmarks.go` | Insert `!` indicator in list rendering; intercept Enter for confirm bookmarks; add confirm toggle to add/edit forms |
| `internal/tui/forms.go` | Intercept post-param execution for confirm bookmarks |
| `internal/tui/model.go` | Add `formConfirmExec` to `formAction` enum |
| `internal/tui/confirm.go` | Add `formConfirmExec` case in `onConfirm()`; "No" returns to bookmark view |

## What Stays the Same

- Config system (`AppConfig`) — untouched
- Shell integration — none needed
- Existing confirm dialog UI — fully reused
- Default behavior — all existing bookmarks execute immediately
