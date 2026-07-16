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
Setting up git...
Setting up nvim...
Setting up tmux...
Setting up zsh...
```

!!! note "`--yes` skips prompts"
    `--yes` answers the per-package confirmation prompts automatically; it does not show them.

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

---

## Example 2: Cross-platform package names

Some programs are packaged under different names on different platforms. You can use the most specific keys for those cases.

For example, the `fd` command-line tool is `fd` in Homebrew but `fd-find` in apt:

```toml
[packages.fd]
target = "$HOME"
sys_package = "fd-find"
sys_package_macos = "fd"
```

On macOS stowmate installs `fd`; on Ubuntu it installs `fd-find`. If you later add a platform where the name differs again, add the specific key for it.

!!! tip "Keep it simple"
    Only add platform-specific keys when the name actually differs. If the package name is the same everywhere, a single `sys_package` is enough.
