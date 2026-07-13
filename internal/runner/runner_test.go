package runner

import (
	"errors"
	"strings"
	"testing"
)

func TestExecRunner_Run(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		wantErr bool
	}{
		{"successful command", "echo", []string{"hello"}, false},
		{"failed command", "false", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewExecRunner(false, nil)
			err := r.Run(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecRunner_RunWithOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		want    string
		wantErr bool
	}{
		{"capture stdout", "echo", []string{"hello"}, "hello\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewExecRunner(false, nil)
			got, err := r.RunWithOutput(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunWithOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("RunWithOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExecRunner_Verbose(t *testing.T) {
	var logged strings.Builder
	r := NewExecRunner(true, func(msg string) { logged.WriteString(msg) })

	if err := r.Run("echo", "verbose", "test"); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if got := logged.String(); got != "echo verbose test" {
		t.Fatalf("verbose log = %q, want %q", got, "echo verbose test")
	}
}

func TestMockRunner_RecordsCommands(t *testing.T) {
	m := NewMockRunner()
	_ = m.Run("brew", "install", "stow")

	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 recorded call, got %d", len(m.Calls))
	}
	if m.Calls[0].Name != "brew" || strings.Join(m.Calls[0].Args, " ") != "install stow" {
		t.Fatalf("unexpected recorded call: %+v", m.Calls[0])
	}
}

func TestMockRunner_ReturnsConfiguredError(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*MockRunner)
		wantError error
	}{
		{
			name: "specific command error",
			setup: func(m *MockRunner) {
				m.SetResponse("brew install stow", ErrSimulated)
			},
			wantError: ErrSimulated,
		},
		{
			name: "binary-wide error",
			setup: func(m *MockRunner) {
				m.SetResponse("brew", errors.New("always fails"))
			},
			wantError: errors.New("always fails"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockRunner()
			tt.setup(m)
			err := m.Run("brew", "install", "stow")
			if !errors.Is(err, tt.wantError) && err.Error() != tt.wantError.Error() {
				t.Fatalf("Run() error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

func TestMockRunner_Verbose(t *testing.T) {
	var logged strings.Builder
	m := NewMockRunner()
	m.SetVerbose(func(msg string) { logged.WriteString(msg) })

	_ = m.Run("stow", "-R", "-t", "/home/user/.config", "nvim")

	want := "stow -R -t /home/user/.config nvim"
	if got := logged.String(); got != want {
		t.Fatalf("verbose log = %q, want %q", got, want)
	}
}

func TestMockRunner_RunWithOutput(t *testing.T) {
	m := NewMockRunner()
	m.SetOutput("echo hello", "hello\n")

	got, err := m.RunWithOutput("echo", "hello")
	if err != nil {
		t.Fatalf("RunWithOutput() unexpected error: %v", err)
	}
	if got != "hello\n" {
		t.Fatalf("RunWithOutput() = %q, want %q", got, "hello\n")
	}
}
