## ADDED Requirements

### Requirement: Detect physical file conflicts
The system SHALL scan the target directory for each package and identify physical files (non-symlinks) that would conflict with stow operations. A conflict exists when a regular file or directory exists at a path where stow would create a symlink.

#### Scenario: No conflicts exist
- **WHEN** the target directory contains only symlinks or is empty
- **THEN** the system SHALL report no conflicts and proceed with stow

#### Scenario: Physical file conflict detected
- **WHEN** a regular file exists at `$HOME/.config/nvim/init.lua` and the `nvim` package would create a symlink at that path
- **THEN** the system SHALL report the conflict and prompt the user for resolution

#### Scenario: Physical directory conflict detected
- **WHEN** a regular directory exists at `$HOME/.config/nvim/` and the `nvim` package would create symlinks inside it
- **THEN** the system SHALL report the directory conflict

### Requirement: Prompt user to resolve conflicts
The system SHALL present each conflict to the user and prompt for deletion. The prompt SHALL be bypassed when `--force` flag is set.

#### Scenario: User confirms conflict deletion
- **WHEN** a conflict is detected and the user confirms deletion
- **THEN** the system SHALL delete the conflicting file or directory

#### Scenario: User declines conflict deletion
- **WHEN** a conflict is detected and the user declines deletion
- **THEN** the system SHALL skip stowing that package and continue to the next

#### Scenario: --force flag is set
- **WHEN** `--force` is provided and conflicts are detected
- **THEN** the system SHALL automatically delete all conflicting files without prompting

### Requirement: Dry-run conflict reporting
The system SHALL report all detected conflicts in dry-run mode without modifying the filesystem.

#### Scenario: Dry-run with conflicts
- **WHEN** `--dry-run` is provided and conflicts exist
- **THEN** the system SHALL print the list of conflicting paths that would be deleted
