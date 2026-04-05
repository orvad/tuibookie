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
- Authentication is handled by your existing git setup -- TuiBookie calls `git` directly, so whatever works for `git clone` and `git push` in your terminal will work here

## Authentication Setup

TuiBookie doesn't manage credentials itself. You need git authentication configured on your machine. Here are the most common options:

### SSH keys (recommended)

Use an SSH clone URL (e.g., `git@github.com:org/shared-bookmarks.git`):

1. Generate a key if you don't have one: `ssh-keygen -t ed25519`
2. Add the public key to your GitHub account under **Settings > SSH and GPG keys**
3. Test with `ssh -T git@github.com`

### HTTPS with a Personal Access Token

Use an HTTPS clone URL (e.g., `https://github.com/org/shared-bookmarks.git`):

1. Go to [GitHub Settings > Developer settings > Personal access tokens > Fine-grained tokens](https://github.com/settings/personal-access-tokens/new)
2. Create a token with **Contents: Read and write** permission scoped to the shared bookmarks repo
3. Configure git to store the credential so you aren't prompted each time:
   ```bash
   git config --global credential.helper store
   ```
4. The first time TuiBookie syncs, git will prompt for your username and token -- after that it's cached

### macOS Keychain

On macOS, git can use the system keychain automatically:

```bash
git config --global credential.helper osxkeychain
```

Your credentials are stored securely and you won't be prompted again after the first sync.
