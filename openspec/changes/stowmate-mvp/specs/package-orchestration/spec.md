## ADDED Requirements

### Requirement: Discover packages from dotfiles directory
The system SHALL scan the dotfiles directory for top-level folders. Each folder SHALL be treated as a stow package. Hidden directories (starting with `.`) and non-directory entries SHALL be excluded.

#### Scenario: Standard dotfiles directory
- **WHEN** the dotfiles directory contains folders `nvim/`, `tmux/`, and `.git/`
- **THEN** the system SHALL identify `nvim` and `tmux` as packages and exclude `.git`

#### Scenario: Dotfiles directory does not exist
- **WHEN** the specified dotfiles path does not exist
- **THEN** the system SHALL print an error and exit with a non-zero code

### Requirement: Prompt user before setting up each package
The system SHALL prompt the user with "Do you want to setup <package>? (Y/n)" before processing each package. The prompt SHALL be bypassed when `--yes` flag is set.

#### Scenario: User confirms package setup
- **WHEN** the user answers "yes" or presses Enter (default yes)
- **THEN** the system SHALL proceed with the package setup loop

#### Scenario: User declines package setup
- **WHEN** the user answers "no"
- **THEN** the system SHALL skip the package and move to the next one

#### Scenario: --yes flag is set
- **WHEN** the `--yes` flag is provided
- **THEN** the system SHALL skip all prompts and process all packages

### Requirement: Install missing dependencies via package manager
The system SHALL check if the resolved system package binary is available. If not found, the system SHALL invoke the detected package manager to install it.

#### Scenario: Dependency is already installed
- **WHEN** the resolved system package binary is found in `$PATH`
- **THEN** the system SHALL skip installation

#### Scenario: Dependency is not installed
- **WHEN** the resolved system package binary is not found in `$PATH`
- **THEN** the system SHALL invoke the package manager install command for the resolved package name

#### Scenario: Dependency installation fails
- **WHEN** the package manager install command returns a non-zero exit code
- **THEN** the system SHALL log the error and prompt the user to continue with remaining packages

### Requirement: Execute pre-clean before stowing
The system SHALL delete all paths listed in the package's `pre_clean` configuration before executing stow.

#### Scenario: Pre-clean paths exist
- **WHEN** pre-clean paths are defined and exist on disk
- **THEN** the system SHALL delete them before stowing

#### Scenario: No pre-clean paths defined
- **WHEN** a package has no `pre_clean` configuration
- **THEN** the system SHALL skip the pre-clean step

### Requirement: Execute stow command
The system SHALL execute `stow -R -t <target_dir> <package>` from within the dotfiles directory for each package being set up.

#### Scenario: Successful stow
- **WHEN** the stow command completes with exit code 0
- **THEN** the system SHALL proceed to post-install hooks

#### Scenario: Stow command fails
- **WHEN** the stow command returns a non-zero exit code
- **THEN** the system SHALL log the error and prompt the user to continue with remaining packages

### Requirement: Execute post-install hooks
The system SHALL execute all `post_install` hook scripts defined in the package's TOML configuration after successful stow.

#### Scenario: Post-install hooks defined
- **WHEN** post-install hooks are defined and stow succeeded
- **THEN** the system SHALL execute each hook sequentially

#### Scenario: No post-install hooks defined
- **WHEN** no post-install hooks are defined for a package
- **THEN** the system SHALL skip the hooks step

### Requirement: Execute package removal
The system SHALL execute `stow -D -t <target_dir> <package>` from within the dotfiles directory when removing a package, followed by any `post_remove` hooks.

#### Scenario: Successful removal
- **WHEN** `stow -D` completes with exit code 0
- **THEN** the system SHALL execute post-remove hooks if defined

#### Scenario: Removal command fails
- **WHEN** `stow -D` returns a non-zero exit code
- **THEN** the system SHALL log the error and report the failure

### Requirement: Graceful error handling with continue prompt
The system SHALL catch errors during package processing and prompt the user: "Setting up '<package>' failed. Continue with remaining packages? (y/N)". In `--yes` mode, the system SHALL automatically continue.

#### Scenario: Package fails and user chooses to continue
- **WHEN** a package setup fails and the user answers "yes" to continue
- **THEN** the system SHALL proceed to the next package

#### Scenario: Package fails and user chooses to stop
- **WHEN** a package setup fails and the user answers "no"
- **THEN** the system SHALL abort and exit with a non-zero code

#### Scenario: Package fails in --yes mode
- **WHEN** a package setup fails and `--yes` flag is set
- **THEN** the system SHALL automatically continue to the next package
