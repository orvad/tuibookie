# Confirm Before Execute

By default, commands fire immediately when you press Enter -- that's the whole point of TuiBookie. But some commands are dangerous, and you may want a safety net before running them.

## How It Works

Individual bookmarks can be marked to require confirmation:

1. When you add or edit a bookmark, set **"Confirm before execute?"** to Yes
2. Bookmarks with confirmation enabled show a bold pink **!** indicator in the list
3. Pressing Enter will display the resolved command in a confirmation dialog
4. You must confirm with `y` before it runs

## Use Cases

This is useful for commands like:

- `rm -rf` -- Destructive file operations
- `kubectl delete` -- Kubernetes resource deletion
- Any command you don't want to fire accidentally
