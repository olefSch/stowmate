# stowmate

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Stowmate is a CLI tool written in Go. It acts as an intelligent orchestrator bridging GNU Stow (for symlink management) and native system package managers (for dependency installation). It detects the host OS, ensures required tools are installed, handles symlink conflicts, and executes post-install hooks.

## Installation

### Go install

```bash
go install github.com/stowmate/stowmate/cmd/stowmate@latest
```

### GitHub Releases

Download a pre-built binary for Linux or macOS from the [Releases](https://github.com/stowmate/stowmate/releases) page.

## Quick start

```bash
# Process every package in your dotfiles directory
stowmate run

# Process a single package
stowmate package nvim

# Remove a stowed package
stowmate remove nvim

# Preview changes without modifying anything
stowmate run --dry-run
```

See the [documentation](docs/) for configuration, package layout, and advanced usage.

## Development

Install [Just](https://github.com/casey/just) (`cargo install just` or `brew install just`), then run:

```bash
just build
just test
just lint
```

## License

[MIT](LICENSE)
