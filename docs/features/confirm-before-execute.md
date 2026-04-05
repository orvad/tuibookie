# Confirm Before Execute

By default, commands fire immediately when you press Enter -- that's the whole point of TuiBookie. But some commands are dangerous, and you may want a safety net before running them.

## How It Works

Individual bookmarks can be marked to require confirmation:

1. When you add or edit a bookmark, set **"Confirm before execute?"** to Yes
2. Bookmarks with confirmation enabled show a bold pink **!** indicator in the list
3. Pressing Enter will display the resolved command in a confirmation dialog
4. You must confirm with `y` before it runs

Bookmarks with confirmation enabled are easy to spot -- notice the <span class="tui-confirm-indicator">!</span> next to "Remove stuff":

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
      <div class="tui-heading">LOCAL COMMANDS</div>
      <div>&nbsp; List stuff <span class="tui-cmd">ls -ahl</span></div>
      <div class="tui-selected">&gt; Remove stuff <span class="tui-confirm-indicator">!</span> <span class="tui-cmd">rm -rf important_stuff</span></div>
      <div class="tui-help">
        <span>[a]</span>dd <span>[e]</span>dit <span>[d]</span>elete
        <span>[enter]</span>run <span>[←/esc]</span> back
        <span>[q]</span>uit
      </div>
    </div>
  </div>
</div>

When you select a confirmed bookmark, you'll see the command and must choose Yes or No before it runs:

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
      <div style="color: #f8f8f2">Execute: rm -rf important_stuff?</div>
      <div class="tui-selected">&gt; Yes</div>
      <div>&nbsp; No</div>
      <div class="tui-help">
        <span>[enter]</span> select <span>[y]</span> yes <span>[n]</span> no <span>[←/esc]</span> back <span>[q]</span> quit
      </div>
    </div>
  </div>
</div>

## Use Cases

This is useful for commands like:

- `rm -rf` -- Destructive file operations
- `kubectl delete` -- Kubernetes resource deletion
- Any command you don't want to fire accidentally
