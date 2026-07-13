## ADDED Requirements

### Requirement: CommandRunner interface
The system SHALL define a `CommandRunner` interface with methods for executing external commands. All external command execution SHALL flow through this interface.

#### Scenario: Interface definition
- **WHEN** the runner package is imported
- **THEN** it SHALL expose a `CommandRunner` interface with `Run(name string, args ...string) error` and `RunWithOutput(name string, args ...string) (string, error)` methods

### Requirement: Production ExecRunner implementation
The system SHALL provide an `ExecRunner` struct that implements `CommandRunner` using `os/exec`. `Run` SHALL execute the command and return an error on non-zero exit. `RunWithOutput` SHALL capture and return stdout.

#### Scenario: Successful command execution
- **WHEN** `Run("echo", "hello")` is called on ExecRunner
- **THEN** the command SHALL execute and return nil error

#### Scenario: Failed command execution
- **WHEN** `Run("false")` is called on ExecRunner
- **THEN** the command SHALL return a non-nil error

#### Scenario: Command with output capture
- **WHEN** `RunWithOutput("echo", "hello")` is called on ExecRunner
- **THEN** the method SHALL return `"hello\n"` and nil error

### Requirement: Mock runner for testing
The system SHALL provide a `MockRunner` struct for use in tests that records all command invocations and returns configurable responses.

#### Scenario: Mock records commands
- **WHEN** `Run("brew", "install", "stow")` is called on MockRunner
- **THEN** the mock SHALL record the command name and arguments for later assertion

#### Scenario: Mock returns configured error
- **WHEN** MockRunner is configured to return an error for a specific command
- **THEN** calling that command SHALL return the configured error

### Requirement: Verbose mode command logging
The system SHALL log the exact command string before execution when verbose mode is enabled.

#### Scenario: Verbose mode enabled
- **WHEN** verbose mode is active and `Run("stow", "-R", "-t", "/home/user/.config", "nvim")` is called
- **THEN** the system SHALL print `stow -R -t /home/user/.config nvim` before executing
