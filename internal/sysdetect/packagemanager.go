package sysdetect

import (
	"errors"
	"fmt"
	"strings"

	"github.com/olefSch/stowmate/internal/runner"
)

// managerChecks lists the supported package managers in priority order. The
// map value is the binary name used to detect the manager.
var managerChecks = []struct {
	name   string
	binary string
}{
	{"brew", "brew"},
	{"apt", "apt-get"},
	{"dnf", "dnf"},
	{"pacman", "pacman"},
	{"zypper", "zypper"},
}

// ErrNoPackageManager is returned when no supported package manager is found.
var ErrNoPackageManager = errors.New("no supported package manager found")

// DetectPackageManager returns the active system package manager by checking
// for known binaries in $PATH. Managers are checked in priority order: brew,
// apt-get, dnf, pacman, zypper.
func DetectPackageManager(runner runner.CommandRunner) (string, error) {
	for _, check := range managerChecks {
		if inPath(runner, check.binary) {
			return check.name, nil
		}
	}
	return "", fmt.Errorf("%w (supported: brew, apt-get, dnf, pacman, zypper)", ErrNoPackageManager)
}

func inPath(r runner.CommandRunner, name string) bool {
	out, err := r.RunWithOutput("which", name)
	return err == nil && strings.TrimSpace(out) != ""
}
