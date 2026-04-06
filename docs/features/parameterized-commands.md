# Parameterized Commands

Commands can include named parameters using `{{name}}` or `{{name:default}}` syntax. When you run a parameterized command, TuiBookie prompts you for each parameter value before executing.

## Syntax

- **`{{name}}`** -- Prompts for a value with no default
- **`{{name:default}}`** -- Prompts for a value, pre-filled with the default

## Example

A bookmark with the command:

```
ssh {{user:admin}}@{{server}}
```

will prompt for `user` (pre-filled with `admin`) and `server` before running.

A live preview of the resolved command updates as you type.

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
        <span class="tui-title-ver">v1.7.1</span>
      </div>
      <div class="tui-sep">━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━</div>
      <div class="tui-heading">Staging servers</div>
      <div class="tui-selected">
        &gt; app
        <span class="tui-cmd">ssh <span class="tui-param">admin</span>@<span class="tui-param">server</span></span>
      </div>
      <div>
        &nbsp; copy stuff
        <span class="tui-cmd">cp -r <span class="tui-param">somefolder</span> <span class="tui-param">someotherfolder</span></span>
      </div>
      <div>
        &nbsp; database
        <span class="tui-cmd">ssh dev@staging-db.example.com</span>
      </div>
      <div class="tui-help">
        <span>[a]</span>dd <span>[e]</span>dit <span>[d]</span>elete
        <span>[enter]</span> run <span>[←/esc]</span> back
        <span>[q]</span>uit
      </div>
    </div>
  </div>
</div>

## Visual Indicators

In the bookmark list, parameters are highlighted in <span style="color: #c87a1a">orange</span> so you can tell at a glance which commands will prompt for input. Commands without parameters run immediately as before.
