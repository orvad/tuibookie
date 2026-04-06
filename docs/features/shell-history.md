# Shell History Integration

Commands executed through TuiBookie can appear in your shell history, so pressing the up arrow recalls them as if you had typed them directly.

## How It Works

After executing a command, TuiBookie writes it to `~/.config/tuibookie/lastcmd`. A small shell wrapper function reads this file and adds the command to your shell's live history.

## Setup

Add the appropriate function to your shell configuration file and restart your shell (or `source` the file).

### Zsh

Add to `~/.zshrc`:

```zsh
tuibookie() {
    command tuibookie "$@"
    local lastcmd="$HOME/.config/tuibookie/lastcmd"
    if [[ -f "$lastcmd" ]]; then
        print -s "$(cat "$lastcmd")"
        rm -f "$lastcmd"
    fi
}
```

### Bash

Add to `~/.bashrc`:

```bash
tuibookie() {
    command tuibookie "$@"
    local lastcmd="$HOME/.config/tuibookie/lastcmd"
    if [[ -f "$lastcmd" ]]; then
        history -s "$(cat "$lastcmd")"
        rm -f "$lastcmd"
    fi
}
```

### Fish

Add to `~/.config/fish/functions/tuibookie.fish`:

```fish
function tuibookie
    command tuibookie $argv
    set lastcmd "$HOME/.config/tuibookie/lastcmd"
    if test -f "$lastcmd"
        builtin history merge
        builtin history add (cat "$lastcmd")
        rm -f "$lastcmd"
    end
end
```

## Notes

- The wrapper only affects history when a command was actually executed. If you quit TuiBookie without running anything, nothing is added.
- Commands with resolved parameters (e.g., `ssh admin@prod-server`) are stored, not the original template.
- If you use a custom config directory via `TUIBOOKIE_CONFIG` or the `-config` flag, the `lastcmd` file is still written to `~/.config/tuibookie/lastcmd`.
