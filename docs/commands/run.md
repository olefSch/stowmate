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

When you run `stowmate run`, stowmate does the following for each package:

1. Load `.stowmate.toml` from the dotfiles directory.
2. Detect the OS and package manager.
3. Ensure GNU Stow is installed.
4. Ask to set up the package (skipped with `--yes`).
5. Resolve the system package name for the current OS and package manager.
6. Install the system package.
7. Run `pre_clean` commands to remove stale files.
8. Detect and resolve symlink conflicts.
9. Run `stow` to create the symlinks.
10. Run `post_install` hooks.

If a package fails, stowmate prompts you to continue. Use `--yes` to continue without asking.

!!! warning "Dry-run mode"
    `--dry-run` only prints the plan; it does not install dependencies, clean files, or create symlinks.
