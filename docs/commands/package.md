# stowmate package

Process a single named package.

## Synopsis

```bash
stowmate package <name> [flags]
```

## Description

`stowmate package <name>` processes exactly one package from the dotfiles directory. It runs the same pipeline as `stowmate run` but only for the package you name. The package must exist as a subdirectory of the dotfiles directory.

This is useful when you want to add, update, or reinstall a single package without touching the rest of your dotfiles.

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
# Process a single package
stowmate package nvim

# Preview what the nvim package would do
stowmate package nvim --dry-run

# Reinstall the zsh package without prompts
stowmate package zsh --yes --force

# Process a package from a custom dotfiles directory
stowmate package git -d ~/projects/dotfiles
```

## Differences from `stowmate run`

| `stowmate run` | `stowmate package <name>` |
|----------------|---------------------------|
| Processes every package | Processes exactly one package |
| Requires no arguments | Requires exactly one argument: the package name |
| Continues on package failure (with prompt) | Stops if the single package fails |

If the named package does not exist, stowmate exits with an error.
