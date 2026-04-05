# Shared Bookmarks via Git

TuiBookie can sync shared bookmarks from a git repository, letting teams maintain a common set of commands alongside their personal bookmarks.

## Setup

1. Go to **Settings > Shared repo** and enter a git clone URL (SSH or HTTPS)
2. Optionally set a custom file path within the repo (defaults to `bookmarks.json`)
3. Select **Sync shared bookmarks** -- TuiBookie clones the repo and loads the bookmarks

Shared bookmarks appear in a separate section below your local bookmarks, marked with a `SHARED` header. When both local and shared bookmarks exist, drill-in views show a breadcrumb like `SHARED > KUBERNETES` so you always know which group you're in.

## Sync Behavior

| Trigger | Behavior |
|---|---|
| **On startup** | Shared bookmarks load instantly from the cached local clone. A background `git pull` fetches updates without blocking the UI. |
| **Manual sync** | Press `S` (capital) in the category view, or use **Sync shared bookmarks** in Settings. |
| **On edit** | Adding, editing, or deleting a shared bookmark triggers an immediate commit and push to the remote repo. |

## Read-Only Repos

TuiBookie detects read-only access on first sync. When a repo is read-only:

- The section header shows `SHARED (READ-ONLY)`
- Add/edit/delete keys are disabled with a status message

This lets you share a curated set of bookmarks with users who shouldn't modify them.

## Requirements

- `git` must be installed and available in your PATH
- Authentication is handled by your existing git setup (SSH keys, credential helpers, tokens)
