# Contributing to TuiBookie

First off — thank you for considering contributing to TuiBookie! Whether it's a bug report, a feature idea, a documentation fix, or a code change, every contribution is welcome and appreciated.

## Getting Started

1. **Fork** the repository and clone your fork locally.
2. Make sure you have [Go](https://go.dev/) installed (see `go.mod` for the required version).
3. Build and run locally:
   ```bash
   go build -o tuibookie .
   ./tuibookie
   ```

## How to Contribute

### Reporting Bugs

Open an [issue](https://github.com/orvad/tuibookie/issues) and include:

- A clear description of what happened vs. what you expected.
- Steps to reproduce the problem.
- Your OS and terminal emulator (TUI rendering can vary).

### Suggesting Features

Have an idea? Open an issue and describe:

- What problem it solves.
- How you imagine it working.

No idea is too small — I'd love to hear it.

### Submitting Changes

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-change
   ```
2. Make your changes. Try to keep commits focused and well-described.
3. Make sure the project builds cleanly:
   ```bash
   go build ./...
   ```
4. Open a pull request against `main`. Describe what you changed and why.

### Code Style

- Follow standard Go conventions (`gofmt` / `goimports`).
- Keep things simple and readable — this is a terminal app, clarity matters.

## Good First Contributions

If you're new here, look for issues labeled **good first issue**. Documentation improvements, typo fixes, and small UI tweaks are always a great way to get started.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thanks! :heart:

— [orvad](https://github.com/orvad)
