## ADDED Requirements

### Requirement: Run command processes all packages
The `stowmate run` command SHALL discover all packages in the dotfiles directory and execute the setup loop for each one.

#### Scenario: Run with default dotfiles path
- **WHEN** the user executes `stowmate run` without flags
- **THEN** the system SHALL use `$HOME/dotfiles` as the dotfiles directory and process all packages

#### Scenario: Run with custom dotfiles path
- **WHEN** the user executes `stowmate run --dotfiles /path/to/dotfiles`
- **THEN** the system SHALL use `/path/to/dotfiles` as the dotfiles directory

### Requirement: Package command processes a single package
The `stowmate package <name>` command SHALL execute the setup loop for the specified package only.

#### Scenario: Package exists
- **WHEN** the user executes `stowmate package nvim` and the `nvim` folder exists
- **THEN** the system SHALL run the setup loop for `nvim` only

#### Scenario: Package does not exist
- **WHEN** the user executes `stowmate package nonexistent` and no such folder exists
- **THEN** the system SHALL print an error and exit with a non-zero code

### Requirement: Remove command un-stows a package
The `stowmate remove <name>` command SHALL execute `stow -D` for the specified package and run any `post_remove` hooks.

#### Scenario: Remove existing package
- **WHEN** the user executes `stowmate remove nvim`
- **THEN** the system SHALL execute `stow -D -t <target> nvim` and run post-remove hooks

### Requirement: Support --dotfiles flag
The system SHALL accept `--dotfiles` (short: `-d`) to specify the dotfiles directory path. Default: `$HOME/dotfiles`.

#### Scenario: --dotfiles flag provided
- **WHEN** `--dotfiles /custom/path` is provided
- **THEN** the system SHALL use `/custom/path` as the dotfiles directory

### Requirement: Support --verbose flag
The system SHALL accept `--verbose` (short: `-v`) to enable debug logging. When enabled, the system SHALL print exact shell commands before execution.

#### Scenario: --verbose flag enabled
- **WHEN** `--verbose` is provided
- **THEN** the system SHALL print each shell command before executing it

### Requirement: Support --dry-run flag
The system SHALL accept `--dry-run` to print planned operations without modifying the filesystem. No package manager installs, file deletions, or stow commands SHALL be executed.

#### Scenario: --dry-run flag enabled
- **WHEN** `--dry-run` is provided
- **THEN** the system SHALL print what would happen (installs, deletions, stow commands) without executing them

### Requirement: Support --yes flag
The system SHALL accept `--yes` (short: `-y`) to run in non-interactive mode. All setup/install prompts SHALL be automatically confirmed.

#### Scenario: --yes flag enabled
- **WHEN** `--yes` is provided
- **THEN** the system SHALL skip all user prompts and answer "yes" to all

### Requirement: Support --force flag
The system SHALL accept `--force` (short: `-f`) to automatically delete physical file conflicts without prompting.

#### Scenario: --force flag enabled
- **WHEN** `--force` is provided and physical file conflicts are detected
- **THEN** the system SHALL delete conflicting files without prompting the user

### Requirement: Support --target flag
The system SHALL accept `--target` (short: `-t`) to override the base target directory. Default: `$HOME`.

#### Scenario: --target flag provided
- **WHEN** `--target /custom/target` is provided
- **THEN** the system SHALL use `/custom/target` as the default stow target directory
