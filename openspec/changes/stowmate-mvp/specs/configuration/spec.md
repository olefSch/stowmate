## ADDED Requirements

### Requirement: Parse .stowmate.toml configuration file
The system SHALL read and parse a `.stowmate.toml` file located in the root of the dotfiles directory. If the file does not exist, the system SHALL proceed with default behavior (folder name = package name, target = $HOME).

#### Scenario: TOML file exists and is valid
- **WHEN** a `.stowmate.toml` file exists in the dotfiles root and contains valid TOML
- **THEN** the system SHALL parse it into the configuration struct

#### Scenario: TOML file does not exist
- **WHEN** no `.stowmate.toml` file exists in the dotfiles root
- **THEN** the system SHALL use default configuration with no package overrides

#### Scenario: TOML file contains syntax errors
- **WHEN** the `.stowmate.toml` file contains invalid TOML syntax
- **THEN** the system SHALL print an error message including the file path and exit with a non-zero code

### Requirement: Resolve package system name via cascading hierarchy
The system SHALL resolve the system package name for each package using the hierarchy: `sys_package_<os>_<manager>` → `sys_package_<os>` → `sys_package` → folder name. The first non-empty value SHALL be used.

#### Scenario: OS and manager specific override exists
- **WHEN** OS is `linux`, manager is `apt`, and `sys_package_linux_apt = "fd-find"` is set
- **THEN** the resolved system package name SHALL be `fd-find`

#### Scenario: Only OS-specific override exists
- **WHEN** OS is `darwin`, `sys_package_macos = "fd"` is set, and no `sys_package_darwin_brew` exists
- **THEN** the resolved system package name SHALL be `fd`

#### Scenario: Only generic override exists
- **WHEN** `sys_package = "neovim"` is set and no OS-specific overrides exist
- **THEN** the resolved system package name SHALL be `neovim`

#### Scenario: No override exists
- **WHEN** no `sys_package` fields are set for a package
- **THEN** the resolved system package name SHALL be the folder name

### Requirement: Support target directory override
The system SHALL allow each package to override the default target directory via the `target` field in TOML. The value SHALL support `$HOME` expansion.

#### Scenario: Target override is specified
- **WHEN** a package has `target = "$HOME/.config"` in TOML
- **THEN** the system SHALL use `$HOME/.config` (with `$HOME` expanded) as the stow target directory

#### Scenario: No target override is specified
- **WHEN** a package has no `target` field in TOML
- **THEN** the system SHALL use the global `--target` flag value (default: `$HOME`)

### Requirement: Support pre-clean paths
The system SHALL allow each package to define a list of `pre_clean` paths in TOML. These paths SHALL be deleted before stowing. Paths SHALL support `$HOME` expansion.

#### Scenario: Pre-clean paths are defined
- **WHEN** a package has `pre_clean = ["$HOME/.cache/nvim", "$HOME/.local/state/nvim"]`
- **THEN** the system SHALL delete those paths (if they exist) before executing stow

#### Scenario: Pre-clean path does not exist
- **WHEN** a pre-clean path points to a non-existent location
- **THEN** the system SHALL silently skip that path without error

### Requirement: Support hook scripts
The system SHALL allow each package to define `post_install` and `post_remove` hook scripts in TOML. Hooks SHALL be executed sequentially in the order defined.

#### Scenario: Post-install hooks are defined
- **WHEN** a package has `post_install = ["git clone https://example.com/repo ~/.config/plugin || true"]`
- **THEN** the system SHALL execute each script via shell after successful stow

#### Scenario: Post-remove hooks are defined
- **WHEN** a package has `post_remove = ["rm -rf ~/.config/plugin"]`
- **THEN** the system SHALL execute each script via shell after successful stow -D

#### Scenario: Hook script fails
- **WHEN** a hook script returns a non-zero exit code
- **THEN** the system SHALL log the error output and continue processing remaining hooks
