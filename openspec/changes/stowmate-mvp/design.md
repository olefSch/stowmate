## Context

Stowmate is a greenfield Go CLI project. The repository is empty aside from LICENSE and README. The tool bridges GNU Stow (symlink management) and system package managers (brew, apt, dnf, pacman, zypper) to automate dotfiles setup across darwin and linux machines.

Constraints from the PRD:
- GNU Stow is the symlink backend — Stowmate never creates symlinks itself
- No templating, no secrets management, no state files
- Must use cobra/viper/charmbracelet ecosystem
- Must be testable via mock CommandRunner and in-memory filesystems

## Goals / Non-Goals

**Goals:**
- Ship a working CLI binary that can `run`, `package`, and `remove` stow packages
- Auto-detect OS and package manager, bootstrap GNU Stow if missing
- Parse `.stowmate.toml` for package overrides and hooks
- Detect and resolve physical file conflicts before stowing
- Execute pre/post hooks defined in TOML
- Provide interactive prompts (charmbracelet/huh) and beautiful output (charmbracelet/log)
- Support --dry-run, --yes, --force, --verbose flags
- Table-driven tests with mock runner and afero filesystem

**Non-Goals:**
- Windows support (darwin and linux only)
- Symlink management (delegated entirely to GNU Stow)
- File content templating or modification
- Secrets/encryption management
- State files, lockfiles, or databases
- Homebrew tap distribution (reserved for v1.0.0)
- Documentation site (separate concern, post-MVP)

## Decisions

### 1. Project Structure: Domain-Driven Layout

**Decision:** Use the PRD-specified structure:
```
cmd/stowmate/main.go
internal/cli/         — Cobra command definitions
internal/sysdetect/   — OS and package manager detection
internal/orchestrator/ — Core execution loop
internal/config/      — TOML parsing and struct definitions
internal/runner/      — os/exec wrapper with CommandRunner interface
```

**Rationale:** Matches PRD requirement. Clean separation of concerns. Each package has a single responsibility. The `internal/` prefix prevents external imports, enforcing module boundaries.

**Alternatives considered:** Flat structure (rejected — too many files in one package), cmd/internal/pkg (rejected — `pkg/` implies external consumers, but this is a CLI binary).

### 2. Configuration Parsing: BurntSushi/toml over Viper

**Decision:** Use `BurntSushi/toml` directly for TOML parsing instead of `spf13/viper`.

**Rationale:** The PRD mentions viper, but viper adds significant complexity (multiple config sources, watchers, key-value merging) that Stowmate doesn't need. Stowmate reads a single `.stowmate.toml` file. `BurntSushi/toml` is simpler, faster, and the de facto Go TOML library. The config struct is straightforward — no need for viper's abstraction layer.

**Alternatives considered:** spf13/viper (rejected — over-engineered for single-file TOML reading), pelletier/go-toml (rejected — BurntSushi/toml is more widely adopted).

### 3. CommandRunner Interface for Testability

**Decision:** Define `CommandRunner` interface in `internal/runner/`:
```go
type CommandRunner interface {
    Run(name string, args ...string) error
    RunWithOutput(name string, args ...string) (string, error)
}
```

Production uses `ExecRunner` wrapping `os/exec`. Tests use `MockRunner` that records calls and returns configured responses.

**Rationale:** Every external command (stow, brew install, apt-get install, hook scripts) flows through this interface. This is the single seam for testing without executing real shell commands.

### 4. Filesystem Abstraction: spf13/afero

**Decision:** Use `spf13/afero` for filesystem operations in conflict resolution and package discovery.

**Rationale:** Conflict resolution needs to scan directories for non-symlink files. Testing this requires an in-memory filesystem. afero provides `MemMapFs` for tests and `OsFs` for production with a clean interface.

### 5. Package Name Resolution: Cascading Lookup

**Decision:** Implement the PRD-specified resolution hierarchy as a simple cascading lookup:
```
sys_package_<os>_<manager> → sys_package_<os> → sys_package → folder_name
```

Implemented as a method on the PackageConfig struct that checks fields in order and returns the first non-empty value.

**Rationale:** Direct mapping from PRD. Simple string checks, no complex logic.

### 6. Error Handling: Per-Package Try/Catch with Continue Prompt

**Decision:** The orchestrator wraps each package setup in a function that returns an error. On error, log the failure and prompt the user to continue or abort. In --yes mode, automatically continue.

**Rationale:** Matches PRD requirement. Prevents one failed package from blocking the entire run. Simple error propagation pattern.

### 7. UI Layer: Charmbracelet Ecosystem

**Decision:** Use `charmbracelet/log` for structured colorful output, `charmbracelet/huh` for interactive prompts (Y/N, confirmations), and a simple spinner implementation for long-running operations.

**Rationale:** PRD requirement. These libraries are well-maintained and composable. In --yes mode, prompts are skipped. In --dry-run mode, operations are printed but not executed.

## Risks / Trade-offs

- **[Risk] GNU Stow not available in package manager** → Mitigation: Print clear error message with manual installation instructions. Don't attempt to compile from source.
- **[Risk] Package manager command requires sudo** → Mitigation: Detect if running as root. If not, prefix apt/dnf/pacman/zypper commands with sudo. Brew on macOS does not use sudo.
- **[Risk] Hook scripts fail silently** → Mitigation: Capture and log hook stdout/stderr. Report non-zero exit codes. In verbose mode, print the full command before execution.
- **[Risk] TOML parsing errors are cryptic** → Mitigation: Wrap TOML decode errors with context (file path, line number if available) and suggest fixes.
- **[Trade-off] No state file means no rollback** → Accepted: The PRD explicitly forbids state files. Users can use `stowmate remove` to undo, or manually manage symlinks via GNU Stow.
