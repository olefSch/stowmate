package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/olefSch/stowmate/internal/orchestrator"
	"github.com/olefSch/stowmate/internal/runner"
)

// TestRunCommandIntegration exercises the full `stowmate run` flow using a mock
// runner, in-memory filesystem, and mock prompter. It verifies that all
// packages are discovered and stowed with the expected arguments.
func TestRunCommandIntegration(t *testing.T) {
	dotfiles := "/mock/dotfiles"
	target := "/mock/target"
	pkgName := "nvim"

	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(filepath.Join(dotfiles, pkgName), 0o755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which brew", "/opt/homebrew/bin/brew")
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")
	mockRunner.SetOutput("which nvim", "/usr/local/bin/nvim")

	prompter := &orchestrator.MockPrompter{Responses: []bool{true}}
	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})

	ctx := &cliContext{
		runner:   mockRunner,
		fs:       fs,
		prompter: prompter,
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"run", "--dotfiles", dotfiles, "--target", target, "--yes"})
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
		t.Fatalf("expected exactly one stow call, got %d", len(stowCalls))
	}
	want := []string{"-d", dotfiles, "-R", "-t", target, pkgName}
	for i, arg := range want {
		if stowCalls[0][i] != arg {
			t.Fatalf("stow arg %d: want %q, got %q", i, arg, stowCalls[0][i])
		}
	}
}

// TestRunCommandIntegration_DryRun verifies that the dry-run path reports the
// planned operation without invoking stow.
func TestRunCommandIntegration_DryRun(t *testing.T) {
	dotfiles := "/mock/dryrun"
	target := "/mock/target"

	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(filepath.Join(dotfiles, "zsh"), 0o755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which brew", "/opt/homebrew/bin/brew")
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")
	mockRunner.SetOutput("which zsh", "/bin/zsh")

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   mockRunner,
		fs:       fs,
		prompter: &orchestrator.MockPrompter{},
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"run", "--dotfiles", dotfiles, "--target", target, "--dry-run", "--yes"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("run command failed: %v", err)
	}

	for _, call := range mockRunner.Calls {
		if call.Name == "stow" {
			t.Fatalf("dry-run should not invoke stow, got %v", call.Args)
		}
	}
}

// TestPackageCommandIntegration exercises `stowmate package nvim` and verifies
// stow is invoked with the expected arguments.
func TestPackageCommandIntegration(t *testing.T) {
	dotfiles := "/mock/pkg/dotfiles"
	target := "/mock/pkg/target"
	pkgName := "nvim"

	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(filepath.Join(dotfiles, pkgName, ".config", pkgName), 0o755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which brew", "/opt/homebrew/bin/brew")
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")
	mockRunner.SetOutput("which nvim", "/usr/local/bin/nvim")

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   mockRunner,
		fs:       fs,
		prompter: &orchestrator.MockPrompter{},
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"package", pkgName, "--dotfiles", dotfiles, "--target", target, "--yes", "--force"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("package command failed: %v", err)
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
	want := []string{"-d", dotfiles, "-R", "-t", target, pkgName}
	for i, arg := range want {
		if stowCalls[0][i] != arg {
			t.Fatalf("stow arg %d: want %q, got %q", i, arg, stowCalls[0][i])
		}
	}
}

// TestRemoveCommandIntegration exercises `stowmate remove nvim` and verifies
// stow -D and post-remove hooks are executed.
func TestRemoveCommandIntegration(t *testing.T) {
	dotfiles := t.TempDir()
	target := "/mock/rm/target"
	pkgName := "nvim"

	if err := os.MkdirAll(filepath.Join(dotfiles, pkgName), 0o755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}
	tomlPath := filepath.Join(dotfiles, ".stowmate.toml")
	tomlData := `[packages.nvim]
target = "/mock/rm/target"
post_remove = ["echo removed"]
`
	if err := os.WriteFile(tomlPath, []byte(tomlData), 0o644); err != nil {
		t.Fatalf("create config: %v", err)
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   mockRunner,
		fs:       afero.NewOsFs(),
		prompter: &orchestrator.MockPrompter{},
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"remove", pkgName, "--dotfiles", dotfiles, "--target", target})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("remove command failed: %v", err)
	}

	var stowCalls [][]string
	var hookCalls [][]string
	for _, call := range mockRunner.Calls {
		switch call.Name {
		case "stow":
			stowCalls = append(stowCalls, call.Args)
		case "sh":
			hookCalls = append(hookCalls, call.Args)
		}
	}
	if len(stowCalls) != 1 {
		t.Fatalf("expected exactly one stow call, got %d: %v", len(stowCalls), mockRunner.Calls)
	}
	want := []string{"-d", dotfiles, "-D", "-t", target, pkgName}
	for i, arg := range want {
		if stowCalls[0][i] != arg {
			t.Fatalf("stow arg %d: want %q, got %q", i, arg, stowCalls[0][i])
		}
	}
	if len(hookCalls) != 1 {
		t.Fatalf("expected exactly one hook call, got %d", len(hookCalls))
	}
}

// TestRunCommandIntegration_MultiplePackages verifies that `stowmate run`
// processes discovered packages in alphabetical order.
func TestRunCommandIntegration_MultiplePackages(t *testing.T) {
	dotfiles := "/mock/multi/dotfiles"
	target := "/mock/multi/target"

	fs := afero.NewMemMapFs()
	for _, pkg := range []string{"nvim", "tmux"} {
		if err := fs.MkdirAll(filepath.Join(dotfiles, pkg), 0o755); err != nil {
			t.Fatalf("create package dir %s: %v", pkg, err)
		}
	}

	mockRunner := runner.NewMockRunner()
	mockRunner.SetOutput("which brew", "/opt/homebrew/bin/brew")
	mockRunner.SetOutput("which stow", "/usr/local/bin/stow")
	mockRunner.SetOutput("which nvim", "/usr/local/bin/nvim")
	mockRunner.SetOutput("which tmux", "/usr/local/bin/tmux")

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   mockRunner,
		fs:       fs,
		prompter: &orchestrator.MockPrompter{Responses: []bool{true, true}},
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
	if len(stowCalls) != 2 {
		t.Fatalf("expected 2 stow calls, got %d: %v", len(stowCalls), mockRunner.Calls)
	}
	want := [][]string{
		{"-d", dotfiles, "-R", "-t", target, "nvim"},
		{"-d", dotfiles, "-R", "-t", target, "tmux"},
	}
	for i, args := range want {
		for j, arg := range args {
			if stowCalls[i][j] != arg {
				t.Fatalf("stow call %d arg %d: want %q, got %q", i, j, arg, stowCalls[i][j])
			}
		}
	}
}

// TestRunCommandIntegration_PackageNotFound verifies that `stowmate package`
// returns an error when the requested package does not exist.
func TestRunCommandIntegration_PackageNotFound(t *testing.T) {
	dotfiles := "/mock/missing/dotfiles"

	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(dotfiles, 0o755); err != nil {
		t.Fatalf("create dotfiles dir: %v", err)
	}

	logger := log.NewWithOptions(io.Discard, log.Options{Level: log.InfoLevel})
	ctx := &cliContext{
		runner:   runner.NewMockRunner(),
		fs:       fs,
		prompter: &orchestrator.MockPrompter{},
		logger:   logger,
	}

	cmd := newRootCmd(ctx)
	cmd.SetArgs([]string{"package", "nonexistent", "--dotfiles", dotfiles, "--yes"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing package")
	}
}
