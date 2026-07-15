# Configuration

Stowmate reads per-package settings from `.stowmate.toml` in the root of the dotfiles directory.

## Location

```
~/dotfiles
├── .stowmate.toml
├── nvim/
├── git/
└── ...
```

If the file is missing, stowmate uses an empty configuration.

## Structure

Each package is configured under the `[packages.<name>]` table.

```toml
[packages.nvim]
target = "$HOME/.config"
sys_package = "neovim"
post_install = ["nvim --headless +PackerSync +q"]
```

## Field reference

| Field | Type | Description |
|-------|------|-------------|
| `target` | string | Directory where the package is stowed. Supports `$HOME` and `~` expansion. |
| `sys_package` | string | Default system package name for all platforms. |
| `sys_package_macos` | string | macOS-specific system package name. |
| `sys_package_linux` | string | Linux-specific system package name. |
| `sys_package_linux_apt` | string | Linux-specific name for `apt`. |
| `sys_package_linux_dnf` | string | Linux-specific name for `dnf`. |
| `sys_package_linux_pacman` | string | Linux-specific name for `pacman`. |
| `sys_package_linux_zypper` | string | Linux-specific name for `zypper`. |
| `pre_clean` | list of strings | Files or directories to remove before stowing. |
| `post_install` | list of strings | Shell commands to run after stowing. |
| `post_remove` | list of strings | Shell commands to run after removing the package. |

## System package resolution

Stowmate picks the system package name using this cascade, from most specific to least specific:

1. `sys_package_<os>_<manager>` — for example, `sys_package_linux_apt` on Ubuntu.
2. `sys_package_<os>` — for example, `sys_package_linux` on any Linux.
3. `sys_package` — the default value for any platform.
4. The package folder name — the final fallback.

For example, with this config:

```toml
[packages.nvim]
sys_package = "neovim"
sys_package_linux_apt = "neovim"
sys_package_linux_dnf = "neovim"
```

On Ubuntu (apt) stowmate installs `neovim`. On Fedora (dnf) it installs `neovim`. On macOS (brew) it falls back to `sys_package` and installs `neovim`.

## Path expansion

The `target` and `pre_clean` fields support simple home-directory expansion:

- `$HOME` and `$HOME/...` are expanded to the user's home directory.
- `~` and `~/...` are expanded to the user's home directory.
- Absolute paths are used as-is.

If expansion fails, the raw value is used.

## Full example

```toml
[packages.nvim]
target = "$HOME/.config"
sys_package = "neovim"
sys_package_macos = "neovim"
sys_package_linux = "neovim"
sys_package_linux_apt = "neovim"
sys_package_linux_dnf = "neovim"
sys_package_linux_pacman = "neovim"
sys_package_linux_zypper = "neovim"
pre_clean = ["$HOME/.cache/nvim"]
post_install = ["nvim --headless +PackerSync +q"]
post_remove = ["echo 'nvim removed'"]

[packages.git]
target = "$HOME"
sys_package = "git"

[packages.tmux]
target = "$HOME"
sys_package = "tmux"
sys_package_linux_apt = "tmux"
sys_package_linux_dnf = "tmux"
post_install = ["tmux source-file ~/.tmux.conf"]

[packages.zsh]
target = "$HOME"
sys_package = "zsh"
post_install = ["chsh -s $(which zsh)"]
```

## Tips

- Use the most specific key only when the package name differs between platforms or package managers.
- Keep `sys_package` as a sensible default so new platforms are covered automatically.
- Use `pre_clean` to remove stale cache or config files that would otherwise block symlinks.
- Use `post_install` for commands that need the dotfiles to be in place first.
