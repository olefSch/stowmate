## 1. Project Foundation

- [x] 1.1 Initialize Go module (`go mod init github.com/stowmate/stowmate`) and create directory structure: `cmd/stowmate/`, `internal/cli/`, `internal/sysdetect/`, `internal/orchestrator/`, `internal/config/`, `internal/runner/`
- [x] 1.2 Add dependencies: `github.com/spf13/cobra`, `github.com/BurntSushi/toml`, `github.com/spf13/afero`, `github.com/charmbracelet/log`, `github.com/charmbracelet/huh`
- [x] 1.3 Create `cmd/stowmate/main.go` entrypoint that calls `internal/cli.Execute()`

## 2. Runner (Command Execution Layer)

- [x] 2.1 Define `CommandRunner` interface in `internal/runner/runner.go` with `Run(name string, args ...string) error` and `RunWithOutput(name string, args ...string) (string, error)`
- [x] 2.2 Implement `ExecRunner` struct that wraps `os/exec` and satisfies `CommandRunner`
- [x] 2.3 Implement `MockRunner` struct in `internal/runner/mock.go` that records invocations and returns configurable responses
- [x] 2.4 Add verbose mode support: log exact command string before execution when verbose flag is active
- [x] 2.5 Write unit tests for `ExecRunner` and `MockRunner`

## 3. Configuration (TOML Parsing)

- [x] 3.1 Define config structs in `internal/config/config.go`: `Config` (top-level), `PackageConfig` (per-package with `target`, `sys_package`, `sys_package_macos`, `sys_package_linux`, `sys_package_linux_apt`, etc., `pre_clean`, `post_install`, `post_remove`)
- [x] 3.2 Implement `LoadConfig(dotfilesPath string) (*Config, error)` that reads `.stowmate.toml` using `BurntSushi/toml` and returns parsed config. Return empty config if file doesn't exist.
- [x] 3.3 Implement `PackageConfig.ResolveSysPackage(os string, manager string) string` method with cascading hierarchy: `sys_package_<os>_<manager>` → `sys_package_<os>` → `sys_package` → folder name
- [x] 3.4 Implement `$HOME` expansion helper for `target` and `pre_clean` paths
- [x] 3.5 Write table-driven tests for TOML parsing, package name resolution, and `$HOME` expansion

## 4. System Detection

- [x] 4.1 Implement `DetectOS() (string, error)` in `internal/sysdetect/os.go` using `runtime.GOOS`, returning `darwin` or `linux`, error on unsupported
- [x] 4.2 Implement `DetectPackageManager(runner CommandRunner) (string, error)` in `internal/sysdetect/packagemanager.go` that checks `$PATH` for `brew`, `apt-get`, `dnf`, `pacman`, `zypper` in priority order
- [x] 4.3 Implement `InstallStow(runner CommandRunner, manager string) error` that generates the correct install command per package manager (with `sudo` prefix for apt/dnf/pacman/zypper)
- [x] 4.4 Implement `EnsureStow(runner CommandRunner, manager string) error` that checks if `stow` is in `$PATH` and calls `InstallStow` if missing
- [x] 4.5 Write table-driven tests for OS detection, package manager detection, and install command generation

## 5. Conflict Resolution

- [x] 5.1 Implement `DetectConflicts(fs afero.Fs, dotfilesPath string, targetDir string, packageName string) ([]string, error)` in `internal/orchestrator/conflict.go` that scans target directory for physical files matching package contents
- [x] 5.2 Implement `ResolveConflicts(conflicts []string, force bool, prompter Prompter) error` that prompts user for each conflict (or auto-deletes with `--force`)
- [x] 5.3 Write tests using `afero.MemMapFs` to verify conflict detection with physical files, symlinks, and empty directories

## 6. Package Orchestration (Core Loop)

- [x] 6.1 Implement `DiscoverPackages(fs afero.Fs, dotfilesPath string) ([]string, error)` in `internal/orchestrator/discover.go` that lists top-level directories excluding hidden ones
- [x] 6.2 Implement `InstallDependency(runner CommandRunner, manager string, sysPackage string) error` that checks if binary is in `$PATH` and installs via package manager if missing
- [x] 6.3 Implement `PreClean(fs afero.Fs, paths []string) error` that deletes configured pre-clean paths
- [x] 6.4 Implement `ExecuteStow(runner CommandRunner, dotfilesPath string, targetDir string, packageName string) error` that runs `stow -R -t <target> <package>`
- [x] 6.5 Implement `ExecuteHooks(runner CommandRunner, hooks []string) error` that runs hook scripts sequentially, logging errors but continuing on failure
- [x] 6.6 Implement `SetupPackage(...)` function that orchestrates the full loop: prompt → dependency check → pre-clean → conflict resolution → stow → hooks
- [x] 6.7 Implement `RemovePackage(...)` function that runs `stow -D -t <target> <package>` followed by `post_remove` hooks
- [x] 6.8 Implement error handling wrapper that catches per-package failures and prompts user to continue or abort

## 7. CLI Interface (Cobra Commands)

- [x] 7.1 Define root command in `internal/cli/root.go` with global flags: `--dotfiles`/`-d`, `--verbose`/`-v`
- [x] 7.2 Define `run` command in `internal/cli/run.go` that discovers and processes all packages
- [x] 7.3 Define `package` command in `internal/cli/package.go` that processes a single named package
- [x] 7.4 Define `remove` command in `internal/cli/remove.go` that un-stows a named package
- [x] 7.5 Add execution flags to applicable commands: `--dry-run`, `--yes`/`-y`, `--force`/`-f`, `--target`/`-t`
- [x] 7.6 Wire flag values into orchestrator calls

## 8. Terminal UI (Charmbracelet)

- [x] 8.1 Configure `charmbracelet/log` as the application logger with colorful output
- [x] 8.2 Implement `Prompter` interface and `HuhPrompter` using `charmbracelet/huh` for Y/N confirmations
- [x] 8.3 Implement `MockPrompter` for testing that returns configurable responses
- [x] 8.4 Add spinner display for long-running operations (package installations)
- [x] 8.5 Implement dry-run output formatting that prints planned operations without executing

## 9. Integration & CI/CD

- [x] 9.1 Write integration test that exercises the full `run` command with mock runner, mock filesystem, and mock prompter
- [x] 9.2 Create `.github/workflows/ci.yml` that runs `go test ./...`, `go vet ./...`, and `golangci-lint` on PRs to main
- [x] 9.3 Create `.goreleaser.yaml` with cross-compilation targets: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, with checksums and changelog
