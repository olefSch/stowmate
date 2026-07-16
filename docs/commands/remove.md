# stowmate remove

Un-stow a named package and run post-remove hooks.

## Synopsis

```bash
stowmate remove <name> [flags]
```

## Description

`stowmate remove <name>` runs `stow -D` to remove the symlinks created by the named package, then runs the commands in its `post_remove` array. It does not uninstall the system package that was installed when the package was set up.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dotfiles` | `-d` | `$HOME/dotfiles` | Path to the dotfiles directory |
| `--target` | `-t` | `$HOME` | Base target directory for stow operations |
| `--dry-run` | | `false` | Print planned operations without executing them |
| `--verbose` | `-v` | `false` | Enable debug logging and print shell commands before execution |

!!! note "No `--yes` or `--force`"
    `remove` does not prompt during the normal flow and does not take `--yes` or `--force`. Use `--dry-run` to preview what will be removed.

## Examples

```bash
# Remove a package
stowmate remove nvim

# Preview what removing a package would do
stowmate remove nvim --dry-run

# Remove a package from a custom dotfiles directory
stowmate remove git -d ~/projects/dotfiles
```

## What happens

1. Load `.stowmate.toml` from the dotfiles directory.
2. Run `stow -D <package>` to delete the symlinks.
3. Run the package's `post_remove` hooks.

If the package is not found in the dotfiles directory, stowmate exits with an error.
