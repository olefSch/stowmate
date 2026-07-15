# Stowmate

A friendly dotfiles manager powered by GNU Stow.

Stowmate sits between your dotfiles repository and your system: it discovers packages, installs their system dependencies, resolves symlink conflicts, and runs hooks—so you can go from a fresh machine to a fully configured environment with one command.

!!! tip "What stowmate is not"
    Stowmate is not a replacement for GNU Stow. It *orchestrates* Stow, plus the package-manager and shell-hook steps that usually surround it.

## What it does

- **:package: Package-aware** — Discovers every package in your dotfiles directory and processes each one individually.
- **:gear: Dependency handling** — Installs system packages automatically via the detected package manager.
- **:link: Conflict-free symlinks** — Detects and resolves stale files before creating Stow symlinks.
- **:rocket: One-command setup** — Goes from a fresh machine to a configured environment with `stowmate run`.

---

## Quick install

=== "Go install"

    Requires Go 1.25 or later.

    ```bash
    go install github.com/olefSch/stowmate/cmd/stowmate@latest
    ```

=== "Binary download"

    Download a pre-built binary for Linux or macOS (amd64 or arm64) from the [GitHub Releases](https://github.com/olefSch/stowmate/releases) page.

---

## Quick start

```bash
# Process every package in $HOME/dotfiles
stowmate run

# Process a single package
stowmate package nvim

# Preview changes without touching the system
stowmate run --dry-run

# Remove a stowed package
stowmate remove nvim
```

!!! note "Default dotfiles path"
    stowmate looks for packages in `$HOME/dotfiles` unless you pass `--dotfiles/-d`.

See the [Introduction](introduction.md) for a detailed walkthrough of how stowmate works.
