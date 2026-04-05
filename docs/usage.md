# Usage

## Launch TuiBookie

```bash
tuibookie
```

## Navigation

TuiBookie presents a two-level interface:

1. **Category list** -- See all your command groups at a glance with bookmark counts
2. **Bookmark list** -- Drill into a category to see commands. Select one and press Enter to run it

Use arrow keys to navigate and Enter to select.

## Screenshots

**Browse categories** -- See all your command groups at a glance with bookmark counts.

![Category list](https://raw.githubusercontent.com/orvad/tuibookie/main/examples/screenshot-01.png)

**Browse bookmarks** -- Drill into a category to see commands. Select one and press Enter to run it.

![Bookmark list](https://raw.githubusercontent.com/orvad/tuibookie/main/examples/screenshot-02.png)

**Settings** -- Configure your bookmarks file path, export backups, or import from a JSON file.

![Settings](https://raw.githubusercontent.com/orvad/tuibookie/main/examples/screenshot-settings.png)

## CLI Flags

```bash
# Use a specific bookmarks file
tuibookie --config /path/to/bookmarks.json
```

You can also set the bookmarks path via environment variable:

```bash
export TUIBOOKIE_CONFIG=/path/to/bookmarks.json
tuibookie
```

Priority order: `--config` flag > `TUIBOOKIE_CONFIG` env var > `config.json` setting > default path.
