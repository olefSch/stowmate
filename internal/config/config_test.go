package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		toml      string
		wantErr   bool
		wantPkg   string
		wantField string
	}{
		{
			name: "valid config",
			toml: `[packages.nvim]
target = "$HOME/.config"
sys_package = "neovim"
pre_clean = ["$HOME/.cache/nvim"]`,
			wantPkg:   "nvim",
			wantField: "neovim",
		},
		{
			name:    "invalid toml",
			toml:    `[packages.nvim\ntarget = "bad`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			dir := "/dotfiles"
			_ = fs.MkdirAll(dir, 0o755)
			_ = afero.WriteFile(fs, filepath.Join(dir, ".stowmate.toml"), []byte(tt.toml), 0o644)

			cfg, err := loadConfig(fs, dir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			pkg, ok := cfg.Packages[tt.wantPkg]
			if !ok {
				t.Fatalf("expected package %q in config", tt.wantPkg)
			}
			if pkg.SysPackage != tt.wantField {
				t.Fatalf("SysPackage = %q, want %q", pkg.SysPackage, tt.wantField)
			}
		})
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg, err := loadConfig(fs, "/missing")
	if err != nil {
		t.Fatalf("loadConfig() unexpected error: %v", err)
	}
	if cfg == nil || len(cfg.Packages) != 0 {
		t.Fatalf("expected empty config, got %+v", cfg)
	}
}

func TestResolveSysPackage(t *testing.T) {
	tests := []struct {
		name       string
		pkg        PackageConfig
		os         string
		manager    string
		folderName string
		want       string
	}{
		{
			name:       "os and manager specific",
			pkg:        PackageConfig{SysPackageLinuxApt: "fd-find"},
			os:         "linux",
			manager:    "apt",
			folderName: "fd",
			want:       "fd-find",
		},
		{
			name:       "os specific",
			pkg:        PackageConfig{SysPackageMacos: "fd"},
			os:         "darwin",
			manager:    "brew",
			folderName: "fd",
			want:       "fd",
		},
		{
			name:       "generic override",
			pkg:        PackageConfig{SysPackage: "neovim"},
			os:         "linux",
			manager:    "apt",
			folderName: "nvim",
			want:       "neovim",
		},
		{
			name:       "fallback to folder name",
			pkg:        PackageConfig{},
			os:         "linux",
			manager:    "apt",
			folderName: "nvim",
			want:       "nvim",
		},
		{
			name:       "apt-get normalised to apt",
			pkg:        PackageConfig{SysPackageLinuxApt: "fd-find"},
			os:         "linux",
			manager:    "apt-get",
			folderName: "fd",
			want:       "fd-find",
		},
		{
			name:       "pacman uses pacman field",
			pkg:        PackageConfig{SysPackageLinuxPacman: "fd"},
			os:         "linux",
			manager:    "pacman",
			folderName: "fd",
			want:       "fd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pkg.ResolveSysPackage(tt.os, tt.manager, tt.folderName)
			if got != tt.want {
				t.Fatalf("ResolveSysPackage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("cannot determine home dir: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{"tilde prefix", "~/.config", filepath.Join(home, ".config")},
		{"$HOME prefix", "$HOME/.local", filepath.Join(home, ".local")},
		{"plain path", "/etc/foo", "/etc/foo"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandHome(tt.path)
			if !strings.HasPrefix(got, home) && tt.want != got {
				t.Fatalf("ExpandHome(%q) = %q, want %q", tt.path, got, tt.want)
			}
			if tt.want != "" && got != tt.want {
				t.Fatalf("ExpandHome(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
