# stowmate run

Discover and process all packages in the dotfiles directory.

## Synopsis

```bash
stowmate run [flags]
```

## Description

`stowmate run` scans the dotfiles directory for subdirectories (packages), then processes each one through the full setup pipeline: install dependency, pre-clean, resolve conflicts, stow, and run post-install hooks.

When a package fails, stowmate asks whether to continue with the remaining packages. In `--yes` mode it continues automatically.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dotfiles` | `-d` | `$HOME/dotfiles` | Path to the dotfiles directory |
| `--target` | `-t` | `$HOME` | Base target directory for stow operations |
| `--dry-run` | | `false` | Print planned operations without executing them |
| `--yes` | `-y` | `false` | Answer yes to all prompts automatically |
| `--force` | `-f` | `false` | Delete physical file conflicts without prompting |
| `--verbose` | `-v` | `false` | Enable debug logging and print shell commands before execution |

## Examples

```bash
# Process all packages in $HOME/dotfiles
stowmate run

# Use a custom dotfiles directory
stowmate run -d ~/projects/dotfiles

# Preview changes without modifying anything
stowmate run --dry-run

# Skip all prompts and delete conflicts automatically
stowmate run --yes --force

# Stow into a custom target directory
stowmate run -t /tmp/demo
```

## Step-by-step behavior

For each package, `stowmate run` loads configuration, installs the declared system dependency, removes stale files, resolves symlink conflicts, creates symlinks with GNU Stow, and runs post-install hooks.

See [How it works](../introduction.md#how-it-works) for the full pipeline.

!!! note "Per-package confirmation"
    By default stowmate asks whether to set up each package. Use `--yes` to skip these prompts and continue automatically.

!!! warning "Dry-run mode"
    `--dry-run` only prints the plan; it does not install dependencies, clean files, or create symlinks.
