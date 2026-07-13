package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

// Config is the top-level structure parsed from .stowmate.toml.
type Config struct {
	Packages map[string]PackageConfig `toml:"packages"`
}

// PackageConfig holds per-package overrides and hooks.
type PackageConfig struct {
	Target                string   `toml:"target"`
	SysPackage            string   `toml:"sys_package"`
	SysPackageMacos       string   `toml:"sys_package_macos"`
	SysPackageLinux       string   `toml:"sys_package_linux"`
	SysPackageLinuxApt    string   `toml:"sys_package_linux_apt"`
	SysPackageLinuxDnf    string   `toml:"sys_package_linux_dnf"`
	SysPackageLinuxPacman string   `toml:"sys_package_linux_pacman"`
	SysPackageLinuxZypper string   `toml:"sys_package_linux_zypper"`
	PreClean              []string `toml:"pre_clean"`
	PostInstall           []string `toml:"post_install"`
	PostRemove            []string `toml:"post_remove"`
}

// LoadConfig reads .stowmate.toml from the dotfiles directory. If the file does
// not exist, an empty Config is returned. Invalid TOML is returned as an error.
func LoadConfig(dotfilesPath string) (*Config, error) {
	return loadConfig(afero.NewOsFs(), dotfilesPath)
}

func loadConfig(fs afero.Fs, dotfilesPath string) (*Config, error) {
	path := filepath.Join(dotfilesPath, ".stowmate.toml")
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &Config{Packages: map[string]PackageConfig{}}, nil
	}

	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	if cfg.Packages == nil {
		cfg.Packages = map[string]PackageConfig{}
	}

	return &cfg, nil
}

// ResolveSysPackage returns the system package name using the cascading
// hierarchy: sys_package_<os>_<manager> → sys_package_<os> → sys_package →
// folderName.
func (p PackageConfig) ResolveSysPackage(os string, manager string, folderName string) string {
	managerKey := manager
	if manager == "apt-get" {
		managerKey = "apt"
	}

	candidates := []string{
		p.fieldForOSManager(os, managerKey),
		p.fieldForOS(os),
		p.SysPackage,
		folderName,
	}

	for _, c := range candidates {
		if strings.TrimSpace(c) != "" {
			return c
		}
	}
	return folderName
}

func (p PackageConfig) fieldForOSManager(os, manager string) string {
	if os == "linux" {
		switch manager {
		case "apt":
			return p.SysPackageLinuxApt
		case "dnf":
			return p.SysPackageLinuxDnf
		case "pacman":
			return p.SysPackageLinuxPacman
		case "zypper":
			return p.SysPackageLinuxZypper
		}
	}
	return ""
}

func (p PackageConfig) fieldForOS(os string) string {
	switch os {
	case "darwin":
		return p.SysPackageMacos
	case "linux":
		return p.SysPackageLinux
	}
	return ""
}

// ExpandHome replaces a leading "$HOME/" or "~/" with the user's home
// directory. If the home directory cannot be determined, the input is returned
// unchanged.
func ExpandHome(path string) string {
	if path == "" {
		return ""
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	if strings.HasPrefix(path, "$HOME/") || path == "$HOME" {
		rest := ""
		if len(path) > len("$HOME") {
			rest = path[len("$HOME/"):]
		}
		if rest == "" {
			return home
		}
		return filepath.Join(home, rest)
	}
	return path
}
