//go:build e2e

package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stowmate/stowmate/internal/orchestrator"
	"github.com/stowmate/stowmate/internal/runner"
)

// TestE2E_RunCommand exercises the full command tree with a real filesystem and
// a mocked shell runner. It verifies the stow invocation produced for a real
// dotfiles package.
func TestE2E_RunCommand(t *testing.T) {
	dotfiles := t.TempDir()
	target := t.TempDir()

	nvimPkg := filepath.Join(dotfiles, "nvim", ".config", "nvim")
	if err := os.MkdirAll(nvimPkg, 0o755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}
	initLua := filepath.Join(nvimPkg, "init.lua")
	if err := os.WriteFile(initLua, []byte("-- nvim config\n"), 0o644); err != nil {
		t.Fatalf("create init.lua: %v", err)
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which brew", "/opt/homebrew/bin/brew")
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")
	mockRunner.SetOutput("which nvim", "/usr/local/bin/nvim")

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   mockRunner,
		fs:       afero.NewOsFs(),
		prompter: &orchestrator.MockPrompter{},
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"run", "--dotfiles", dotfiles, "--target", target, "--yes", "--force"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("run command failed: %v", err)
	}

	var stowCalls [][]string
	for _, call := range mockRunner.Calls {
		if call.Name == "stow" {
			stowCalls = append(stowCalls, call.Args)
		}
	}
	if len(stowCalls) != 1 {
		t.Fatalf("expected exactly one stow call, got %d: %v", len(stowCalls), mockRunner.Calls)
	}
	want := []string{"-d", dotfiles, "-R", "-t", target, "nvim"}
	for i, arg := range want {
		if stowCalls[0][i] != arg {
			t.Fatalf("stow arg %d: want %q, got %q", i, arg, stowCalls[0][i])
		}
	}
}
