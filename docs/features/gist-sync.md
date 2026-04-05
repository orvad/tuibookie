# Gist Sync

Push your bookmarks to a secret GitHub Gist and pull them on any machine. Versioned backup with full revision history, powered by a Personal Access Token.

## Setup

### Creating a Personal Access Token

TuiBookie needs a GitHub Personal Access Token (PAT) with permission to create and update gists.

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click **Generate new token** (classic)
3. Give it a descriptive name like "TuiBookie Gist Sync"
4. Under **Select scopes**, check only **gist** -- no other permissions are needed
5. Click **Generate token** and copy the token immediately (you won't see it again)

### Adding the token to TuiBookie

1. Open TuiBookie and go to **Settings > GitHub token**
2. Paste the token and confirm
3. The token is stored in `~/.config/tuibookie/config.json` and displayed masked in the UI

!!! warning
    Keep your token secure. Anyone with the token can read and modify your gists. If you suspect it's been compromised, revoke it immediately from [GitHub token settings](https://github.com/settings/tokens) and generate a new one.

## Push to Gist

Go to **Settings > Push to Gist** to upload your bookmarks:

- On first push, a new secret gist is created
- Subsequent pushes update the existing gist
- Full revision history is maintained by GitHub

## Pull from Gist

Go to **Settings > Pull from Gist** to download bookmarks:

- Shows a confirmation with category and bookmark counts before overwriting
- Replaces the local bookmarks file with the gist contents

## Use Case

Gist sync is ideal for keeping your personal bookmarks consistent across multiple machines -- your work laptop, home desktop, and cloud VMs all stay in sync.
