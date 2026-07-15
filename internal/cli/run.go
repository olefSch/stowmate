package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/olefSch/stowmate/internal/config"
	"github.com/olefSch/stowmate/internal/orchestrator"
	"github.com/olefSch/stowmate/internal/runner"
	"github.com/olefSch/stowmate/internal/sysdetect"
)

// newRunCmd creates the `stowmate run` command.
func newRunCmd(ctx *cliContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Discover and process all packages in the dotfiles directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dotfiles := expandPath(dotfilesPath)
			target := expandPath(targetDir)

			cfg, err := config.LoadConfig(dotfiles)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			packages, err := orchestrator.DiscoverPackages(ctx.fs, dotfiles)
			if err != nil {
				return fmt.Errorf("discover packages: %w", err)
			}
			if len(packages) == 0 {
				ctx.logger.Info("No packages found")
				return nil
			}

			applyTargetOverride(cfg, packages, target, cmd.Flags().Changed("target"))

			osName, manager, err := detectEnvironment(ctx.runner)
			if err != nil {
				return err
			}

			if dryRun {
				ctx.logger.Info("Would ensure GNU Stow is installed", "manager", manager)
			} else if err := sysdetect.EnsureStow(ctx.runner, manager); err != nil {
				return fmt.Errorf("ensure stow: %w", err)
			}

			o := &orchestrator.Orchestrator{
				Runner:   ctx.runner,
				Fs:       ctx.fs,
				Prompter: ctx.prompter,
				Config:   cfg,
				Verbose:  verbose,
				DryRun:   dryRun,
				Force:    force,
				Yes:      yes,
			}

			if dryRun || !yes || !force {
				return o.SetupPackages(packages, dotfiles, osName, manager)
			}

			return runWithSpinner("Running stowmate", func() error {
				return o.SetupPackages(packages, dotfiles, osName, manager)
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print planned operations without executing them")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Answer yes to all prompts automatically")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Delete physical file conflicts without prompting")
	cmd.Flags().StringVarP(&targetDir, "target", "t", "$HOME", "Base target directory for stow operations")

	return cmd
}

// detectEnvironment returns the current OS and package manager.
func detectEnvironment(r runner.CommandRunner) (string, string, error) {
	osName, err := sysdetect.DetectOS()
	if err != nil {
		return "", "", fmt.Errorf("detect OS: %w", err)
	}
	manager, err := sysdetect.DetectPackageManager(r)
	if err != nil {
		return "", "", fmt.Errorf("detect package manager: %w", err)
	}
	return osName, manager, nil
}

// applyTargetOverride sets the target directory on every discovered package so
// the CLI --target flag overrides the config default of $HOME.
func applyTargetOverride(cfg *config.Config, packages []string, target string, changed bool) {
	if !changed {
		return
	}
	for _, pkg := range packages {
		pc := cfg.Packages[pkg]
		pc.Target = target
		cfg.Packages[pkg] = pc
	}
}
