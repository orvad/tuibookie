# Configuration

## Config file

On first launch, TuiBookie creates a config file at:

```
~/.config/tuibookie/config.json
```

This stores your app settings, starting with the bookmarks file path. You can change the bookmarks path directly from the Settings view in the TUI -- no flags needed.

## Bookmarks file location

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

## Settings view

The Settings view in TuiBookie provides access to:

| Setting | Description |
|---|---|
| **Bookmarks file** | View and change the path to your bookmarks JSON file |
| **Export bookmarks** | Save a backup as `bookmarks-backup-YYYY-MM-DD-HHMMSS.json` |
| **Import bookmarks** | Import from a `.json` file, merged into existing categories |
| **Push to Gist** | Upload bookmarks to a secret GitHub Gist |
| **Pull from Gist** | Download bookmarks from your gist |
| **GitHub token** | Set or remove the Personal Access Token for Gist sync |
| **Shared repo** | Set the git clone URL for shared bookmarks |
| **Shared file path** | Set the path to bookmarks within the shared repo |
| **Sync shared bookmarks** | Pull latest shared bookmarks from the remote |
| **Disconnect shared repo** | Remove shared repo configuration |
