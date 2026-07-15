# stowmate

[![CI](https://github.com/olefSch/stowmate/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/olefSch/stowmate/actions/workflows/ci.yml)
[![Docs](https://github.com/olefSch/stowmate/actions/workflows/docs.yml/badge.svg?branch=main)](https://github.com/olefSch/stowmate/actions/workflows/docs.yml)
[![Release](https://github.com/olefSch/stowmate/actions/workflows/release.yml/badge.svg)](https://github.com/olefSch/stowmate/actions/workflows/release.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/olefSch/stowmate)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A friendly dotfiles manager powered by [GNU Stow](https://www.gnu.org/software/stow/).

Stowmate discovers packages in your dotfiles directory, installs their system dependencies, resolves symlink conflicts, and runs hooks — so you can go from a fresh machine to a fully configured environment with one command.

## Features

- **Package-aware** — Each subdirectory in your dotfiles repo is a Stow package, processed individually
- **Dependency handling** — Installs system packages via the detected package manager (brew, apt, dnf, pacman, zypper)
- **Conflict resolution** — Detects and removes stale files before creating symlinks
- **Hooks** — Runs `pre_clean`, `post_install`, and `post_remove` shell commands per package
- **Cross-platform** — Supports macOS and Linux with per-OS and per-manager package name overrides
- **Dry-run mode** — Preview every operation before it touches your system

## Installation

### Go install

Requires Go 1.25 or later.

```bash
go install github.com/olefSch/stowmate/cmd/stowmate@latest
```

### Binary download

Download a pre-built binary for Linux or macOS (amd64 / arm64) from the [Releases](https://github.com/olefSch/stowmate/releases) page.

## Quick start

```bash
# Process every package in $HOME/dotfiles
stowmate run

# Process a single package
stowmate package nvim

# Preview changes without modifying anything
stowmate run --dry-run

# Remove a stowed package
stowmate remove nvim
```

## Configuration

Place a `.stowmate.toml` file in the root of your dotfiles directory to configure per-package behavior:

```toml
[packages.nvim]
target = "$HOME/.config"
sys_package = "neovim"
pre_clean = ["$HOME/.cache/nvim"]
post_install = ["nvim --headless +PackerSync +q"]

[packages.fd]
sys_package = "fd-find"
sys_package_macos = "fd"
```

See the [Configuration reference](https://olefsch.github.io/stowmate/configuration/) for all fields and the system package resolution cascade.

## Documentation

Full documentation is available at **[olefsch.github.io/stowmate](https://olefsch.github.io/stowmate/)**:

- [Introduction](https://olefsch.github.io/stowmate/introduction/) — How stowmate works, supported platforms, directory layout
- [Commands](https://olefsch.github.io/stowmate/commands/run/) — Detailed reference for `run`, `package`, and `remove`
- [Configuration](https://olefsch.github.io/stowmate/configuration/) — `.stowmate.toml` field reference
- [Examples](https://olefsch.github.io/stowmate/examples/) — Complete worked dotfiles setup

## Development

Install [Just](https://github.com/casey/just), then:

```bash
just build          # Build the binary
just test           # Run tests with race detection
just lint           # Run golangci-lint
just check          # Run vet, lint, test, and vuln check
just docs-serve     # Preview documentation locally
```

## License

[MIT](LICENSE)
