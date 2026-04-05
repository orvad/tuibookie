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

![Parameterized commands](https://raw.githubusercontent.com/orvad/tuibookie/main/examples/screenshot-parameterized-commands.png)

## Visual Indicators

In the bookmark list, parameters are highlighted so you can tell at a glance which commands will prompt for input. Commands without parameters run immediately as before.
