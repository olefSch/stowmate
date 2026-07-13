package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandRunner executes external commands. All external command execution in
// Stowmate flows through this interface so tests can substitute a mock.
type CommandRunner interface {
	Run(name string, args ...string) error
	RunWithOutput(name string, args ...string) (string, error)
}

// ExecRunner is the production implementation of CommandRunner using os/exec.
type ExecRunner struct {
	Verbose bool
	Logger  func(string)
}

// NewExecRunner returns an ExecRunner. If logger is nil, verbose output is
// written to stdout.
func NewExecRunner(verbose bool, logger func(string)) *ExecRunner {
	if logger == nil {
		logger = func(msg string) { fmt.Println(msg) }
	}
	return &ExecRunner{Verbose: verbose, Logger: logger}
}

// Run executes a command and returns an error on non-zero exit.
func (r *ExecRunner) Run(name string, args ...string) error {
	r.logCommand(name, args)
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return err
	}
	return nil
}

// RunWithOutput executes a command and returns its stdout.
func (r *ExecRunner) RunWithOutput(name string, args ...string) (string, error) {
	r.logCommand(name, args)
	cmd := exec.Command(name, args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return out.String(), nil
}

func (r *ExecRunner) logCommand(name string, args []string) {
	if r.Verbose {
		r.Logger(strings.Join(append([]string{name}, args...), " "))
	}
}
