# Installation

## Quick install (recommended)

Automatically downloads the latest release for your OS and architecture:

```bash
curl -sL https://raw.githubusercontent.com/orvad/tuibookie/main/install.sh | sh
```

## Homebrew (macOS and Linux)

```bash
brew tap orvad/tuibookie
brew install tuibookie
```

## Download manually

Download the latest binary from the [Releases page](https://github.com/orvad/tuibookie/releases):

| Platform | Binary |
|---|---|
| macOS (Apple Silicon) | `tuibookie-darwin-arm64` |
| macOS (Intel) | `tuibookie-darwin-amd64` |
| Linux (x86_64) | `tuibookie-linux-amd64` |
| Linux (ARM64) | `tuibookie-linux-arm64` |

Then make it executable and move it to your PATH:

```bash
chmod +x tuibookie-*
sudo mv tuibookie-* /usr/local/bin/tuibookie
```

## Build from source

Requires [Go](https://go.dev/dl/) 1.26 or later:

```bash
git clone https://github.com/orvad/tuibookie.git
cd tuibookie
go build -o tuibookie .
sudo mv tuibookie /usr/local/bin/
```

### Cross-compilation

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o tuibookie .

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o tuibookie .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o tuibookie .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o tuibookie .
```
