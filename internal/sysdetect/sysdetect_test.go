package sysdetect

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stowmate/stowmate/internal/runner"
)

func TestDetectOS(t *testing.T) {
	got, err := DetectOS()
	if err != nil {
		t.Fatalf("DetectOS() unexpected error: %v", err)
	}

	switch runtime.GOOS {
	case "darwin", "linux":
		if got != runtime.GOOS {
			t.Fatalf("DetectOS() = %q, want %q", got, runtime.GOOS)
		}
	default:
		t.Fatalf("DetectOS() = %q, want error on unsupported OS", got)
	}
}

func TestDetectPackageManager(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*runner.MockRunner)
		want      string
		wantError bool
	}{
		{
			name: "brew found",
			setup: func(m *runner.MockRunner) {
				m.SetOutput("which brew", "/opt/homebrew/bin/brew")
			},
			want: "brew",
		},
		{
			name: "apt found",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which brew", errors.New("not found"))
				m.SetOutput("which apt-get", "/usr/bin/apt-get")
			},
			want: "apt",
		},
		{
			name: "dnf found",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which brew", errors.New("not found"))
				m.SetResponse("which apt-get", errors.New("not found"))
				m.SetOutput("which dnf", "/usr/bin/dnf")
			},
			want: "dnf",
		},
		{
			name: "pacman found",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which brew", errors.New("not found"))
				m.SetResponse("which apt-get", errors.New("not found"))
				m.SetResponse("which dnf", errors.New("not found"))
				m.SetOutput("which pacman", "/usr/bin/pacman")
			},
			want: "pacman",
		},
		{
			name: "zypper found",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which brew", errors.New("not found"))
				m.SetResponse("which apt-get", errors.New("not found"))
				m.SetResponse("which dnf", errors.New("not found"))
				m.SetResponse("which pacman", errors.New("not found"))
				m.SetOutput("which zypper", "/usr/bin/zypper")
			},
			want: "zypper",
		},
		{
			name: "none found",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which brew", errors.New("not found"))
				m.SetResponse("which apt-get", errors.New("not found"))
				m.SetResponse("which dnf", errors.New("not found"))
				m.SetResponse("which pacman", errors.New("not found"))
				m.SetResponse("which zypper", errors.New("not found"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := runner.NewMockRunner()
			tt.setup(m)

			got, err := DetectPackageManager(m)
			if (err != nil) != tt.wantError {
				t.Fatalf("DetectPackageManager() error = %v, wantError %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Fatalf("DetectPackageManager() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInstallStow(t *testing.T) {
	tests := []struct {
		name    string
		manager string
		want    runner.Invocation
	}{
		{
			name:    "brew",
			manager: "brew",
			want:    runner.Invocation{Name: "brew", Args: []string{"install", "stow"}},
		},
		{
			name:    "apt",
			manager: "apt",
			want:    runner.Invocation{Name: "sudo", Args: []string{"apt-get", "install", "-y", "stow"}},
		},
		{
			name:    "dnf",
			manager: "dnf",
			want:    runner.Invocation{Name: "sudo", Args: []string{"dnf", "install", "-y", "stow"}},
		},
		{
			name:    "pacman",
			manager: "pacman",
			want:    runner.Invocation{Name: "sudo", Args: []string{"pacman", "-S", "--noconfirm", "stow"}},
		},
		{
			name:    "zypper",
			manager: "zypper",
			want:    runner.Invocation{Name: "sudo", Args: []string{"zypper", "install", "-y", "stow"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := runner.NewMockRunner()
			if err := InstallStow(m, tt.manager); err != nil {
				t.Fatalf("InstallStow() unexpected error: %v", err)
			}
			if len(m.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(m.Calls))
			}
			got := m.Calls[0]
			if got.Name != tt.want.Name || !slicesEqual(got.Args, tt.want.Args) {
				t.Fatalf("InstallStow() invoked %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsureStow(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*runner.MockRunner)
		wantCalls int
	}{
		{
			name: "already installed",
			setup: func(m *runner.MockRunner) {
				m.SetOutput("which stow", "/usr/bin/stow")
			},
			wantCalls: 1,
		},
		{
			name: "missing installs",
			setup: func(m *runner.MockRunner) {
				m.SetResponse("which stow", errors.New("not found"))
			},
			wantCalls: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := runner.NewMockRunner()
			tt.setup(m)

			if err := EnsureStow(m, "brew"); err != nil {
				t.Fatalf("EnsureStow() unexpected error: %v", err)
			}
			if len(m.Calls) != tt.wantCalls {
				t.Fatalf("expected %d calls, got %d", tt.wantCalls, len(m.Calls))
			}
		})
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
