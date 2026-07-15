package orchestrator

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/olefSch/stowmate/internal/config"
	"github.com/olefSch/stowmate/internal/runner"
)

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
}

func TestDetectConflicts_NoConflicts(t *testing.T) {
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")

	conflicts, err := DetectConflicts(fs, "/dotfiles", "/home/user", "nvim")
	if err != nil {
		t.Fatalf("DetectConflicts() error = %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", conflicts)
	}
}

func TestDetectConflicts_PhysicalFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")
	_ = afero.WriteFile(fs, "/home/user/.config/nvim/init.lua", []byte("-- old"), 0o644)

	conflicts, err := DetectConflicts(fs, "/dotfiles", "/home/user", "nvim")
	if err != nil {
		t.Fatalf("DetectConflicts() error = %v", err)
	}
	want := "/home/user/.config/nvim/init.lua"
	if len(conflicts) != 1 || conflicts[0] != want {
		t.Fatalf("conflicts = %v, want [%s]", conflicts, want)
	}
}

func makeNvimPackage(fs afero.Fs, dotfilesPath string) {
	_ = fs.MkdirAll(filepath.Join(dotfilesPath, "nvim/.config/nvim"), 0o755)
	_ = afero.WriteFile(fs, filepath.Join(dotfilesPath, "nvim/.config/nvim/init.lua"), []byte("-- nvim"), 0o644)
}

func TestDetectConflicts_SymlinkIgnored(t *testing.T) {
	fs := newSymlinkFs(afero.NewMemMapFs())
	makeNvimPackage(fs, "/dotfiles")
	_ = afero.WriteFile(fs, "/home/user/real.lua", []byte("x"), 0o644)
	_ = fs.SymlinkIfPossible("/home/user/real.lua", "/home/user/.config/nvim/init.lua")

	conflicts, err := DetectConflicts(fs, "/dotfiles", "/home/user", "nvim")
	if err != nil {
		t.Fatalf("DetectConflicts() error = %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts for symlink, got %v", conflicts)
	}
}

func TestDetectConflicts_EmptyDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")
	_ = fs.MkdirAll("/home/user/.config/nvim", 0o755)

	conflicts, err := DetectConflicts(fs, "/dotfiles", "/home/user", "nvim")
	if err != nil {
		t.Fatalf("DetectConflicts() error = %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %v", conflicts)
	}
	if !strings.HasSuffix(conflicts[0], "/.config/nvim") {
		t.Fatalf("unexpected conflict path: %s", conflicts[0])
	}
}

func TestDetectConflicts_PackageNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	_, err := DetectConflicts(fs, "/dotfiles", "/home/user", "missing")
	if err == nil {
		t.Fatal("expected error for missing package")
	}
}

func TestResolveConflicts_ForceDeletes(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/conflict1", []byte("x"), 0o644)
	_ = afero.WriteFile(fs, "/conflict2", []byte("y"), 0o644)

	err := ResolveConflicts(fs, []string{"/conflict1", "/conflict2"}, true, nil)
	if err != nil {
		t.Fatalf("ResolveConflicts() error = %v", err)
	}
	if exists, _ := afero.Exists(fs, "/conflict1"); exists {
		t.Fatal("conflict1 should have been deleted")
	}
	if exists, _ := afero.Exists(fs, "/conflict2"); exists {
		t.Fatal("conflict2 should have been deleted")
	}
}

func TestResolveConflicts_PromptDecline(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/conflict", []byte("x"), 0o644)

	prompter := &MockPrompter{Responses: []bool{false}}
	err := ResolveConflicts(fs, []string{"/conflict"}, false, prompter)
	if err == nil {
		t.Fatal("expected error when user declines")
	}
	if exists, _ := afero.Exists(fs, "/conflict"); !exists {
		t.Fatal("conflict should not have been deleted")
	}
}

func TestResolveConflicts_PromptAccept(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/conflict", []byte("x"), 0o644)

	prompter := &MockPrompter{Responses: []bool{true}}
	err := ResolveConflicts(fs, []string{"/conflict"}, false, prompter)
	if err != nil {
		t.Fatalf("ResolveConflicts() error = %v", err)
	}
	if exists, _ := afero.Exists(fs, "/conflict"); exists {
		t.Fatal("conflict should have been deleted")
	}
}

func TestDiscoverPackages(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/dotfiles/nvim/.config/nvim", 0o755)
	_ = fs.MkdirAll("/dotfiles/tmux", 0o755)
	_ = fs.MkdirAll("/dotfiles/.git", 0o755)
	_ = afero.WriteFile(fs, "/dotfiles/README.md", []byte("# dotfiles"), 0o644)

	got, err := DiscoverPackages(fs, "/dotfiles")
	if err != nil {
		t.Fatalf("DiscoverPackages() error = %v", err)
	}
	want := []string{"nvim", "tmux"}
	if !slicesEqual(got, want) {
		t.Fatalf("DiscoverPackages() = %v, want %v", got, want)
	}
}

func TestDiscoverPackages_MissingDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	_, err := DiscoverPackages(fs, "/missing")
	if err == nil {
		t.Fatal("expected error for missing dotfiles directory")
	}
}

func TestInstallDependency_AlreadyInstalled(t *testing.T) {
	m := runner.NewMockRunner()
	m.SetOutput("which fd", "/usr/bin/fd")

	if err := InstallDependency(m, "brew", "fd"); err != nil {
		t.Fatalf("InstallDependency() error = %v", err)
	}
	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(m.Calls))
	}
	if m.Calls[0].Name != "which" {
		t.Fatalf("expected which call, got %v", m.Calls[0])
	}
}

func TestInstallDependency_NotInstalled(t *testing.T) {
	m := runner.NewMockRunner()
	m.SetResponse("which fd", errors.New("not found"))

	if err := InstallDependency(m, "apt", "fd"); err != nil {
		t.Fatalf("InstallDependency() error = %v", err)
	}
	if len(m.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(m.Calls))
	}

	want := runner.Invocation{Name: "sudo", Args: []string{"apt-get", "install", "-y", "fd"}}
	if m.Calls[1].Name != want.Name || !slicesEqual(m.Calls[1].Args, want.Args) {
		t.Fatalf("install invocation = %v, want %v", m.Calls[1], want)
	}
}

func TestInstallDependency_UnsupportedManager(t *testing.T) {
	m := runner.NewMockRunner()
	m.SetResponse("which x", errors.New("not found"))

	err := InstallDependency(m, "unknown", "x")
	if err == nil {
		t.Fatal("expected error for unsupported manager")
	}
}

func TestPreClean(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/a", []byte("a"), 0o644)
	_ = afero.WriteFile(fs, "/b/c", []byte("c"), 0o644)

	if err := PreClean(fs, []string{"/a", "/b"}); err != nil {
		t.Fatalf("PreClean() error = %v", err)
	}
	if exists, _ := afero.Exists(fs, "/a"); exists {
		t.Fatal("/a should have been deleted")
	}
	if exists, _ := afero.Exists(fs, "/b"); exists {
		t.Fatal("/b should have been deleted")
	}
}

func TestExecuteStow(t *testing.T) {
	m := runner.NewMockRunner()
	if err := ExecuteStow(m, "/dotfiles", "/home/user/.config", "nvim"); err != nil {
		t.Fatalf("ExecuteStow() error = %v", err)
	}
	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(m.Calls))
	}

	want := runner.Invocation{Name: "stow", Args: []string{"-d", "/dotfiles", "-R", "-t", "/home/user/.config", "nvim"}}
	if m.Calls[0].Name != want.Name || !slicesEqual(m.Calls[0].Args, want.Args) {
		t.Fatalf("ExecuteStow() invoked %v, want %v", m.Calls[0], want)
	}
}

func TestExecuteHooks_ContinuesOnFailure(t *testing.T) {
	m := runner.NewMockRunner()
	m.SetResponse("sh -c false", errors.New("exit status 1"))

	err := ExecuteHooks(m, []string{"echo ok", "false", "echo done"})
	if err != nil {
		t.Fatalf("ExecuteHooks() error = %v", err)
	}
	if len(m.Calls) != 3 {
		t.Fatalf("expected 3 hook calls, got %d", len(m.Calls))
	}
}

func TestOrchestrator_SetupPackage_Success(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")

	m := runner.NewMockRunner()
	m.SetResponse("which neovim", errors.New("not found"))

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{true}},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"nvim": {Target: "$HOME", SysPackage: "neovim"},
			},
		},
	}

	if err := o.SetupPackage("nvim", "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackage() error = %v", err)
	}

	wantCalls := []runner.Invocation{
		{Name: "which", Args: []string{"neovim"}},
		{Name: "sudo", Args: []string{"apt-get", "install", "-y", "neovim"}},
		{Name: "stow", Args: []string{"-d", "/dotfiles", "-R", "-t", "/home/user", "nvim"}},
	}
	if len(m.Calls) != len(wantCalls) {
		t.Fatalf("expected %d calls, got %d: %v", len(wantCalls), len(m.Calls), m.Calls)
	}
	for i, want := range wantCalls {
		if m.Calls[i].Name != want.Name || !slicesEqual(m.Calls[i].Args, want.Args) {
			t.Fatalf("call %d = %v, want %v", i, m.Calls[i], want)
		}
	}
}

func TestOrchestrator_SetupPackage_UserSkips(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/dotfiles/nvim", 0o755)

	m := runner.NewMockRunner()
	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{false}},
		Config:   &config.Config{Packages: map[string]config.PackageConfig{}},
	}

	if err := o.SetupPackage("nvim", "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackage() error = %v", err)
	}
	if len(m.Calls) != 0 {
		t.Fatalf("expected no calls when user skips, got %v", m.Calls)
	}
}

func TestOrchestrator_SetupPackage_DryRun(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")
	_ = afero.WriteFile(fs, "/home/user/.config/nvim/init.lua", []byte("old"), 0o644)

	m := runner.NewMockRunner()
	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{true}},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"nvim": {Target: "$HOME", SysPackage: "neovim"},
			},
		},
		DryRun: true,
	}

	if err := o.SetupPackage("nvim", "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackage() error = %v", err)
	}
	if len(m.Calls) != 0 {
		t.Fatalf("expected no runner calls in dry-run, got %v", m.Calls)
	}
	if exists, _ := afero.Exists(fs, "/home/user/.config/nvim/init.lua"); !exists {
		t.Fatal("dry-run should not delete conflicts")
	}
}

func TestOrchestrator_RemovePackage(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	m := runner.NewMockRunner()

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"nvim": {
					Target:     "$HOME",
					PostRemove: []string{"echo removed"},
				},
			},
		},
	}

	if err := o.RemovePackage("nvim", "/dotfiles"); err != nil {
		t.Fatalf("RemovePackage() error = %v", err)
	}

	if len(m.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d: %v", len(m.Calls), m.Calls)
	}
	wantStow := runner.Invocation{Name: "stow", Args: []string{"-d", "/dotfiles", "-D", "-t", "/home/user", "nvim"}}
	if m.Calls[0].Name != wantStow.Name || !slicesEqual(m.Calls[0].Args, wantStow.Args) {
		t.Fatalf("stow call = %v, want %v", m.Calls[0], wantStow)
	}
	wantHook := runner.Invocation{Name: "sh", Args: []string{"-c", "echo removed"}}
	if m.Calls[1].Name != wantHook.Name || !slicesEqual(m.Calls[1].Args, wantHook.Args) {
		t.Fatalf("hook call = %v, want %v", m.Calls[1], wantHook)
	}
}

func TestOrchestrator_SetupPackages_ContinueOnFailure(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/dotfiles/a", 0o755)
	_ = fs.MkdirAll("/dotfiles/b", 0o755)

	m := runner.NewMockRunner()
	m.SetResponse("stow -d /dotfiles -R -t /home/user a", errors.New("stow failed"))

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{true, true, true}},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"a": {Target: "$HOME"},
				"b": {Target: "$HOME"},
			},
		},
	}

	if err := o.SetupPackages([]string{"a", "b"}, "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackages() error = %v", err)
	}
	if len(m.Calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", len(m.Calls))
	}
}

func TestOrchestrator_SetupPackages_StopOnFailure(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/dotfiles/a", 0o755)
	_ = fs.MkdirAll("/dotfiles/b", 0o755)

	m := runner.NewMockRunner()
	m.SetResponse("stow -d /dotfiles -R -t /home/user a", errors.New("stow failed"))

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{true, false}},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"a": {Target: "$HOME"},
				"b": {Target: "$HOME"},
			},
		},
	}

	err := o.SetupPackages([]string{"a", "b"}, "/dotfiles", "linux", "apt")
	if err == nil {
		t.Fatal("expected error when user stops")
	}
}

func TestOrchestrator_SetupPackages_YesModeContinues(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/dotfiles/a", 0o755)
	_ = fs.MkdirAll("/dotfiles/b", 0o755)

	m := runner.NewMockRunner()
	m.SetResponse("stow -d /dotfiles -R -t /home/user a", errors.New("stow failed"))

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"a": {Target: "$HOME"},
				"b": {Target: "$HOME"},
			},
		},
		Yes: true,
	}

	if err := o.SetupPackages([]string{"a", "b"}, "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackages() error = %v", err)
	}
}

func TestOrchestrator_SetupPackage_ResolveConflicts(t *testing.T) {
	setHome(t, "/home/user")
	fs := afero.NewMemMapFs()
	makeNvimPackage(fs, "/dotfiles")
	_ = afero.WriteFile(fs, "/home/user/.config/nvim/init.lua", []byte("old"), 0o644)

	m := runner.NewMockRunner()
	m.SetOutput("which neovim", "/usr/bin/neovim")

	o := &Orchestrator{
		Runner:   m,
		Fs:       fs,
		Prompter: &MockPrompter{Responses: []bool{true, true}},
		Config: &config.Config{
			Packages: map[string]config.PackageConfig{
				"nvim": {Target: "$HOME"},
			},
		},
	}

	if err := o.SetupPackage("nvim", "/dotfiles", "linux", "apt"); err != nil {
		t.Fatalf("SetupPackage() error = %v", err)
	}
	if exists, _ := afero.Exists(fs, "/home/user/.config/nvim/init.lua"); exists {
		t.Fatal("conflict should have been deleted")
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

// symlinkFs wraps an afero.Fs with minimal symlink support for tests.
// It is used because afero.MemMapFs does not implement afero.Symlinker.
type symlinkFs struct {
	afero.Fs
	links map[string]string
}

func newSymlinkFs(fs afero.Fs) *symlinkFs {
	return &symlinkFs{Fs: fs, links: map[string]string{}}
}

func (s *symlinkFs) SymlinkIfPossible(oldname, newname string) error {
	s.links[newname] = oldname
	return nil
}

func (s *symlinkFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	if target, ok := s.links[name]; ok {
		fi, err := s.Stat(target)
		if err != nil {
			return nil, true, err
		}
		return &symlinkFileInfo{FileInfo: fi}, true, nil
	}
	fi, err := s.Stat(name)
	return fi, false, err
}

func (s *symlinkFs) ReadlinkIfPossible(name string) (string, bool, error) {
	if target, ok := s.links[name]; ok {
		return target, true, nil
	}
	return "", false, nil
}

type symlinkFileInfo struct {
	os.FileInfo
}

func (s *symlinkFileInfo) Mode() os.FileMode {
	return s.FileInfo.Mode() | os.ModeSymlink
}
