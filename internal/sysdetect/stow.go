package sysdetect

import (
	"fmt"
	"strings"

	"github.com/olefSch/stowmate/internal/runner"
)

// InstallStow installs GNU Stow using the detected package manager.
func InstallStow(r runner.CommandRunner, manager string) error {
	name, args, err := stowInstallCommand(manager)
	if err != nil {
		return err
	}
	return r.Run(name, args...)
}

// EnsureStow installs GNU Stow if it is not already present in $PATH.
func EnsureStow(r runner.CommandRunner, manager string) error {
	if stowInstalled(r) {
		return nil
	}
	return InstallStow(r, manager)
}

func stowInstalled(r runner.CommandRunner) bool {
	out, err := r.RunWithOutput("which", "stow")
	return err == nil && strings.TrimSpace(out) != ""
}

func stowInstallCommand(manager string) (string, []string, error) {
	switch manager {
	case "brew":
		return "brew", []string{"install", "stow"}, nil
	case "apt", "apt-get":
		return "sudo", []string{"apt-get", "install", "-y", "stow"}, nil
	case "dnf":
		return "sudo", []string{"dnf", "install", "-y", "stow"}, nil
	case "pacman":
		return "sudo", []string{"pacman", "-S", "--noconfirm", "stow"}, nil
	case "zypper":
		return "sudo", []string{"zypper", "install", "-y", "stow"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}
