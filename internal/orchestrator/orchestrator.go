package orchestrator

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/stowmate/stowmate/internal/config"
	"github.com/stowmate/stowmate/internal/runner"
)

// Orchestrator coordinates package setup and removal. It is intentionally
// independent of the CLI layer so the same core logic can be exercised from
// tests and from Cobra commands in Milestone 3.
type Orchestrator struct {
	Runner   runner.CommandRunner
	Fs       afero.Fs
	Prompter Prompter
	Config   *config.Config
	Verbose  bool
	DryRun   bool
	Force    bool
	Yes      bool
}

// SetupPackage runs the full setup loop for a single package: prompt,
// dependency installation, pre-clean, conflict resolution, stow, and hooks.
func (o *Orchestrator) SetupPackage(packageName string, dotfilesPath string, osName string, manager string) error {
	pkgCfg := o.Config.Packages[packageName]

	if !o.Yes {
		ok, err := o.Prompter.Confirm(fmt.Sprintf("Do you want to setup %s? (Y/n)", packageName))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	targetDir, err := resolveTargetDir(pkgCfg.Target)
	if err != nil {
		return err
	}

	sysPackage := pkgCfg.ResolveSysPackage(osName, manager, packageName)

	if o.DryRun {
		o.reportDryRun(packageName, targetDir, sysPackage)
		conflicts, err := DetectConflicts(o.Fs, dotfilesPath, targetDir, packageName)
		if err != nil {
			return err
		}
		if len(conflicts) > 0 {
			fmt.Println("[dry-run] Conflicts detected:")
			for _, c := range conflicts {
				fmt.Printf("  - %s\n", c)
			}
		}
		return nil
	}

	if err := InstallDependency(o.Runner, manager, sysPackage); err != nil {
		return fmt.Errorf("install dependency %s: %w", sysPackage, err)
	}

	if err := PreClean(o.Fs, expandPaths(pkgCfg.PreClean)); err != nil {
		return fmt.Errorf("pre-clean: %w", err)
	}

	conflicts, err := DetectConflicts(o.Fs, dotfilesPath, targetDir, packageName)
	if err != nil {
		return fmt.Errorf("detect conflicts: %w", err)
	}
	if len(conflicts) > 0 {
		if err := ResolveConflicts(o.Fs, conflicts, o.Force, o.Prompter); err != nil {
			return fmt.Errorf("resolve conflicts: %w", err)
		}
	}

	if err := ExecuteStow(o.Runner, dotfilesPath, targetDir, packageName); err != nil {
		return fmt.Errorf("stow: %w", err)
	}

	if err := ExecuteHooks(o.Runner, pkgCfg.PostInstall); err != nil {
		return fmt.Errorf("post-install hooks: %w", err)
	}

	return nil
}

// RemovePackage unstows a package and runs post-remove hooks.
func (o *Orchestrator) RemovePackage(packageName string, dotfilesPath string) error {
	pkgCfg := o.Config.Packages[packageName]

	targetDir, err := resolveTargetDir(pkgCfg.Target)
	if err != nil {
		return err
	}

	if o.DryRun {
		fmt.Printf("[dry-run] Would remove %s from %s\n", packageName, targetDir)
		return nil
	}

	if err := executeStowDelete(o.Runner, dotfilesPath, targetDir, packageName); err != nil {
		return fmt.Errorf("stow -D: %w", err)
	}

	if err := ExecuteHooks(o.Runner, pkgCfg.PostRemove); err != nil {
		return fmt.Errorf("post-remove hooks: %w", err)
	}

	return nil
}

// SetupPackages processes each package in order. When a package fails, the user
// is prompted to continue or abort. In --yes mode failures are logged and
// processing continues automatically.
func (o *Orchestrator) SetupPackages(packages []string, dotfilesPath string, osName string, manager string) error {
	for _, pkg := range packages {
		err := o.SetupPackage(pkg, dotfilesPath, osName, manager)
		if err == nil {
			continue
		}

		fmt.Printf("Setting up '%s' failed: %v\n", pkg, err)
		if o.Yes {
			continue
		}

		ok, promptErr := o.Prompter.Confirm(fmt.Sprintf("Setting up '%s' failed. Continue with remaining packages? (y/N)", pkg))
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) reportDryRun(packageName string, targetDir string, sysPackage string) {
	fmt.Printf("[dry-run] Would setup %s (target: %s", packageName, targetDir)
	if sysPackage != "" {
		fmt.Printf(", sys_package: %s", sysPackage)
	}
	fmt.Println(")")
}

func resolveTargetDir(target string) (string, error) {
	targetDir := config.ExpandHome(target)
	if targetDir != "" {
		return targetDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve target dir: %w", err)
	}
	return home, nil
}

func expandPaths(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = config.ExpandHome(p)
	}
	return out
}
