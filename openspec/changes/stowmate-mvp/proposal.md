## Why

Managing dotfiles across multiple machines and operating systems requires manual intervention: detecting which tools are installed, resolving file conflicts, installing missing dependencies, and running GNU Stow with the correct flags. Stowmate automates this orchestration layer while keeping GNU Stow as the symlink backend, eliminating repetitive setup friction for developers who maintain cross-platform dotfile repositories.

## What Changes

- **New Go CLI binary** (`stowmate`) with three commands: `run`, `package <name>`, `remove <name>`
- **OS and package manager auto-detection** for darwin/linux with brew, apt, dnf, pacman, zypper support
- **Automatic dependency bootstrapping** — installs GNU Stow and target package binaries via detected package manager
- **TOML-based configuration** (`.stowmate.toml`) for package name overrides, target directories, pre-clean paths, and hook scripts
- **Physical file conflict detection and resolution** with interactive prompts before stowing
- **Pre/post hook execution** for cache cleanup and post-install scripts (e.g., TPM cloning)
- **Graceful error handling** — per-package failure catching with continue/skip prompts
- **Interactive terminal UI** using Charmbracelet ecosystem (huh for prompts, log for output, spinners for long operations)
- **Dry-run mode** that prints planned operations without modifying the filesystem
- **Domain-driven Go project structure** with `CommandRunner` interface for testability

## Capabilities

### New Capabilities
- `environment-detection`: OS detection (darwin/linux), package manager detection (brew/apt/dnf/pacman/zypper), and GNU Stow bootstrapping
- `configuration`: TOML parsing via viper with package name resolution hierarchy, target overrides, pre-clean paths, and hook script definitions
- `package-orchestration`: Core execution loop — discovery, user prompt, dependency install, pre-clean, conflict resolution, stow invocation, and hook execution
- `cli-interface`: Cobra-based CLI with `run`, `package`, `remove` commands and global/execution flags (--dotfiles, --verbose, --dry-run, --yes, --force, --target)
- `conflict-resolution`: Scanning target directories for non-symlink physical files and prompting for deletion before stowing
- `runner`: CommandRunner interface wrapping os/exec with mock implementation for testing

### Modified Capabilities
<!-- None — this is a greenfield project -->

## Impact

- **New codebase**: Entire Go project created from scratch under `cmd/` and `internal/`
- **Dependencies**: spf13/cobra, spf13/viper, BurntSushi/toml, charmbracelet/huh, charmbracelet/log, charmbracelet/x/ansi (or lipgloss)
- **Testing**: Table-driven tests for config/sysdetect, afero-based filesystem mocking for conflict resolution, mock CommandRunner for shell command assertions
- **CI/CD**: GitHub Actions workflow for `go test`, `go vet`, `golangci-lint`; GoReleaser config for cross-compilation (darwin/linux × amd64/arm64)
- **No breaking changes** — greenfield project, no existing consumers
