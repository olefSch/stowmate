package orchestrator

import (
	"fmt"
	"strings"

	"github.com/stowmate/stowmate/internal/runner"
)

// InstallDependency checks whether the system package binary is available in
// $PATH and installs it via the package manager when missing.
func InstallDependency(r runner.CommandRunner, manager string, sysPackage string) error {
	if isInPATH(r, sysPackage) {
		return nil
	}

	name, args, err := installCommand(manager, sysPackage)
	if err != nil {
		return err
	}
	return r.Run(name, args...)
}

func installCommand(manager string, sysPackage string) (string, []string, error) {
	switch manager {
	case "brew":
		return "brew", []string{"install", sysPackage}, nil
	case "apt", "apt-get":
		return "sudo", []string{"apt-get", "install", "-y", sysPackage}, nil
	case "dnf":
		return "sudo", []string{"dnf", "install", "-y", sysPackage}, nil
	case "pacman":
		return "sudo", []string{"pacman", "-S", "--noconfirm", sysPackage}, nil
	case "zypper":
		return "sudo", []string{"zypper", "install", "-y", sysPackage}, nil
	default:
		return "", nil, fmt.Errorf("unsupported package manager: %s", manager)
	}
}

func isInPATH(r runner.CommandRunner, name string) bool {
	out, err := r.RunWithOutput("which", name)
	return err == nil && strings.TrimSpace(out) != ""
}
