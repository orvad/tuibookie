# Gist Sync

Push your bookmarks to a secret GitHub Gist and pull them on any machine. Versioned backup with full revision history, powered by a Personal Access Token.

## Setup

1. Create a [Personal Access Token](https://github.com/settings/tokens) with `gist` scope
2. In TuiBookie, go to **Settings > GitHub token** and enter the token
3. The token is stored in `config.json` and displayed masked in the UI

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
