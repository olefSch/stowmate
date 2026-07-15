# Examples

## Example 1: A full dotfiles setup

This example shows a realistic dotfiles repository for a developer who uses Neovim, Git, Tmux, and Zsh across macOS and Linux.

### Directory structure

```
~/dotfiles
├── .stowmate.toml
├── git/
│   └── .gitconfig
├── nvim/
│   └── .config/
│       └── nvim/
│           └── init.lua
├── tmux/
│   └── .tmux.conf
└── zsh/
    └── .zshrc
```

### Configuration

```toml
[packages.git]
target = "$HOME"
sys_package = "git"

[packages.nvim]
target = "$HOME/.config"
sys_package = "neovim"
sys_package_linux_apt = "neovim"
sys_package_linux_dnf = "neovim"
sys_package_linux_pacman = "neovim"
sys_package_linux_zypper = "neovim"
pre_clean = ["$HOME/.cache/nvim"]
post_install = ["nvim --headless +PackerSync +q"]
post_remove = ["echo 'nvim removed'"]

[packages.tmux]
target = "$HOME"
sys_package = "tmux"
post_install = ["tmux source-file ~/.tmux.conf"]

[packages.zsh]
target = "$HOME"
sys_package = "zsh"
post_install = ["chsh -s $(which zsh)"]
```

### Running stowmate

```bash
$ cd ~/dotfiles
$ stowmate run --dry-run
[dry-run] Would setup git (target: /home/alex)
[dry-run] Would setup nvim (target: /home/alex/.config, sys_package: neovim)
[dry-run] Would setup tmux (target: /home/alex)
[dry-run] Would setup zsh (target: /home/alex)
```

Once the preview looks correct, run it for real:

```bash
$ stowmate run --yes
Do you want to setup git? (Y/n) y
Do you want to setup nvim? (Y/n) y
Do you want to setup tmux? (Y/n) y
Do you want to setup zsh? (Y/n) y
```

With `--yes`, the prompts are skipped and stowmate proceeds through all packages.

After it finishes, your home directory contains symlinks such as:

```
~/.config/nvim -> /home/alex/dotfiles/nvim/.config/nvim
~/.gitconfig   -> /home/alex/dotfiles/git/.gitconfig
~/.tmux.conf   -> /home/alex/dotfiles/tmux/.tmux.conf
~/.zshrc       -> /home/alex/dotfiles/zsh/.zshrc
```

### Updating a single package

If you edit `nvim/.config/nvim/init.lua`, you can reinstall just the nvim package:

```bash
stowmate package nvim --yes
```

### Removing a package

```bash
stowmate remove tmux
```

This removes the `~/.tmux.conf` symlink. It does not uninstall the `tmux` package from the system.

## Example 2: Cross-platform package names

Some programs are packaged under different names on different platforms. You can use the most specific keys for those cases.

```toml
[packages.ripgrep]
target = "$HOME"
sys_package = "ripgrep"
sys_package_macos = "ripgrep"
sys_package_linux_apt = "ripgrep"
sys_package_linux_dnf = "ripgrep"
sys_package_linux_pacman = "ripgrep"
sys_package_linux_zypper = "ripgrep"
```

Stowmate resolves the name for the current OS and package manager automatically. If you later add a platform where the package name differs, add the specific key for it.

!!! tip "Keep it simple"
    Only add platform-specific keys when the name actually differs. If the package name is the same everywhere, a single `sys_package` is enough.
