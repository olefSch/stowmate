# Introduction

Stowmate is a CLI tool for managing dotfiles with [GNU Stow](https://www.gnu.org/software/stow/). It turns a directory of package folders into symlinks in your home directory, while also handling the boring parts that Stow leaves to you: installing system packages, cleaning up stale files, resolving symlink conflicts, and running post-install hooks.

## Why stowmate exists

Dotfiles repositories often look like this:

```
~/dotfiles
├── nvim/
├── git/
├── tmux/
└── zsh/
```

Each folder is a *Stow package*. With plain Stow, turning those folders into symlinks is already simple, but you still have to:

- Make sure Stow is installed.
- Install the actual programs the dotfiles configure.
- Remove stale files that would block the symlink.
- Run setup commands after linking (e.g. `nvim --headless +PackerSync +q`).

Stowmate automates that surrounding work. It also gives you a per-package config file (`.stowmate.toml`) so each package can declare its own target directory, system dependency, and hooks.

---

## How it works

Running `stowmate run` executes a pipeline for each package.

### Phase 1: Discovery & Detection

1. **Discover** — Find every subdirectory in the dotfiles directory.
2. **Detect environment** — Identify the OS and package manager (brew, apt, dnf, pacman, zypper).
3. **Ensure Stow** — Install GNU Stow if it is missing.

### Phase 2: Preparation

4. **Load config** — Read `.stowmate.toml` in the dotfiles directory.
5. **Resolve system package** — Pick the right dependency name for the current OS and package manager.
6. **Install dependency** — Install the system package via the detected package manager.
7. **Pre-clean** — Remove files listed in the package's `pre_clean` array.

### Phase 3: Installation

8. **Detect conflicts** — Find files that would block Stow from creating symlinks.
9. **Resolve conflicts** — Prompt to delete conflicts (or delete automatically with `--force`).
10. **Stow** — Run `stow` to create the symlinks.
11. **Post-install hooks** — Run the commands listed in `post_install`.

!!! tip "Skip prompts with `--yes`"
    Pass `--yes` to answer all setup and continuation prompts automatically. Use `--force` to delete file conflicts without asking.

---

## Installation

### Go install

Requires Go 1.25 or later.

```bash
go install github.com/stowmate/stowmate/cmd/stowmate@latest
```

### GitHub Releases

Download a pre-built binary for Linux or macOS (amd64 or arm64) from the [Releases](https://github.com/stowmate/stowmate/releases) page.

## Prerequisites

- **For the binary:** None. GNU Stow is installed automatically if it is missing.
- **For Go install:** Go 1.25 or later.
- **For package installation:** A supported package manager must be available on the system.

---

## Supported platforms

| OS      | Package managers |
|---------|------------------|
| macOS   | `brew`           |
| Linux   | `apt`, `dnf`, `pacman`, `zypper` |

---

## Dotfiles directory layout

A dotfiles directory is a collection of Stow packages. Each package is a subdirectory. For example:

```
~/dotfiles
├── nvim/
│   └── .config/
│       └── nvim/
│           └── init.lua
├── git/
│   └── .gitconfig
├── tmux/
│   └── .tmux.conf
└── zsh/
    └── .zshrc
```

Each package folder is processed independently. When you run `stowmate run`, it stows every package into `$HOME` by default.

!!! tip "Configuration file"
    Place `.stowmate.toml` in the root of the dotfiles directory to customize per-package targets, dependencies, and hooks. See the [Configuration](configuration.md) reference.

---

## Next steps

- Read the [Commands](commands/run.md) reference for each subcommand.
- Read the [Configuration](configuration.md) reference for `.stowmate.toml`.
- See the [Examples](examples.md) page for a complete worked setup.
