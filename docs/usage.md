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
        <span class="tui-title-ver">v1.6.0</span>
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

**Browse bookmarks** -- Drill into a category to see commands. Select one and press Enter to run it.

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
        <span class="tui-title-ver">v1.6.0</span>
      </div>
      <div class="tui-sep">━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━</div>
      <div class="tui-heading">Kubernetes</div>
      <div class="tui-selected">
        &gt; API logs
        <span class="tui-cmd">kubectl logs -f deployment/api -n production</span>
      </div>
      <div>
        &nbsp; Node resources
        <span class="tui-cmd">kubectl top nodes</span>
      </div>
      <div>
        &nbsp; Prod pods
        <span class="tui-cmd">kubectl get pods -n production</span>
      </div>
      <div class="tui-help">
        <span>[a]</span>dd <span>[e]</span>dit <span>[d]</span>elete
        <span>[enter]</span> run <span>[←/esc]</span> back
        <span>[q]</span>uit
      </div>
    </div>
  </div>
</div>

**Settings** -- Configure your bookmarks file path, export backups, or import from a JSON file.

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
        <span class="tui-title-ver">v1.6.0</span>
      </div>
      <div class="tui-sep">━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━</div>
      <div class="tui-heading">Settings</div>
      <div class="tui-section-label">CONFIG</div>
      <div>&nbsp; Bookmarks file: ~/.config/tuibookie/bookmarks.json</div>
      <div>&nbsp;</div>
      <div class="tui-section-label">DATA</div>
      <div>&nbsp; Export bookmarks</div>
      <div>&nbsp; Import bookmarks</div>
      <div>&nbsp;</div>
      <div class="tui-section-label">SYNC</div>
      <div>&nbsp; Push to Gist</div>
      <div>&nbsp; Pull from Gist</div>
      <div>&nbsp; GitHub token: ****hEss</div>
      <div>&nbsp;</div>
      <div class="tui-section-label">SHARED</div>
      <div class="tui-selected">&gt; Shared repo: git@github.com:team/devops.git</div>
      <div>&nbsp; Shared file path: bookmarks.json</div>
      <div>&nbsp; Sync shared bookmarks</div>
      <div>&nbsp; Disconnect shared repo</div>
      <div class="tui-help">
        <span>[enter/→]</span> select
        <span>[←/esc]</span> back <span>[q]</span> quit
      </div>
    </div>
  </div>
</div>

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
