## ADDED Requirements

### Requirement: Detect operating system
The system SHALL detect the host operating system at startup and expose it as a normalized value: `darwin` or `linux`. Detection SHALL use `runtime.GOOS`.

#### Scenario: Running on macOS
- **WHEN** the binary executes on a macOS host
- **THEN** the detected OS value SHALL be `darwin`

#### Scenario: Running on Linux
- **WHEN** the binary executes on a Linux host
- **THEN** the detected OS value SHALL be `linux`

#### Scenario: Running on unsupported OS
- **WHEN** the binary executes on an OS other than darwin or linux
- **THEN** the system SHALL print an error message and exit with a non-zero code

### Requirement: Detect system package manager
The system SHALL detect the active system package manager by checking for the presence of known binaries in `$PATH`. Supported managers: `brew`, `apt-get`, `dnf`, `pacman`, `zypper`.

#### Scenario: Homebrew is installed
- **WHEN** `brew` is found in `$PATH`
- **THEN** the detected package manager SHALL be `brew`

#### Scenario: apt-get is installed
- **WHEN** `apt-get` is found in `$PATH` and `brew` is not found
- **THEN** the detected package manager SHALL be `apt`

#### Scenario: dnf is installed
- **WHEN** `dnf` is found in `$PATH` and neither `brew` nor `apt-get` are found
- **THEN** the detected package manager SHALL be `dnf`

#### Scenario: pacman is installed
- **WHEN** `pacman` is found in `$PATH` and no higher-priority manager is found
- **THEN** the detected package manager SHALL be `pacman`

#### Scenario: zypper is installed
- **WHEN** `zypper` is found in `$PATH` and no higher-priority manager is found
- **THEN** the detected package manager SHALL be `zypper`

#### Scenario: No supported package manager found
- **WHEN** none of the supported package manager binaries are found in `$PATH`
- **THEN** the system SHALL print an error listing the supported managers and exit with a non-zero code

### Requirement: Bootstrap GNU Stow if missing
The system SHALL check if `stow` is available in `$PATH`. If not found, the system SHALL invoke the detected package manager to install GNU Stow before proceeding.

#### Scenario: Stow is already installed
- **WHEN** `stow` is found in `$PATH`
- **THEN** the system SHALL skip installation and proceed normally

#### Scenario: Stow is not installed and brew is detected
- **WHEN** `stow` is not in `$PATH` and the detected package manager is `brew`
- **THEN** the system SHALL execute `brew install stow`

#### Scenario: Stow is not installed and apt is detected
- **WHEN** `stow` is not in `$PATH` and the detected package manager is `apt`
- **THEN** the system SHALL execute `sudo apt-get install -y stow`

#### Scenario: Stow is not installed and dnf is detected
- **WHEN** `stow` is not in `$PATH` and the detected package manager is `dnf`
- **THEN** the system SHALL execute `sudo dnf install -y stow`

#### Scenario: Stow is not installed and pacman is detected
- **WHEN** `stow` is not in `$PATH` and the detected package manager is `pacman`
- **THEN** the system SHALL execute `sudo pacman -S --noconfirm stow`

#### Scenario: Stow is not installed and zypper is detected
- **WHEN** `stow` is not in `$PATH` and the detected package manager is `zypper`
- **THEN** the system SHALL execute `sudo zypper install -y stow`

#### Scenario: Stow installation fails
- **WHEN** the package manager command to install stow returns a non-zero exit code
- **THEN** the system SHALL print the error output and exit with a non-zero code
