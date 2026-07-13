package cli

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stowmate/stowmate/internal/config"
	"github.com/stowmate/stowmate/internal/orchestrator"
	"github.com/stowmate/stowmate/internal/sysdetect"
)

// newPackageCmd creates the `stowmate package <name>` command.
func newPackageCmd(ctx *cliContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package <name>",
		Short: "Process a single named package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			dotfiles := expandPath(dotfilesPath)
			target := expandPath(targetDir)

			exists, err := afero.Exists(ctx.fs, dotfiles)
			if err != nil {
				return fmt.Errorf("check dotfiles path: %w", err)
			}
			if !exists {
				return fmt.Errorf("dotfiles path %q does not exist", dotfiles)
			}

			pkgPath := dotfiles + "/" + packageName
			exists, err = afero.DirExists(ctx.fs, pkgPath)
			if err != nil {
				return fmt.Errorf("check package path: %w", err)
			}
			if !exists {
				return fmt.Errorf("package %q not found at %s", packageName, pkgPath)
			}

			cfg, err := config.LoadConfig(dotfiles)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			applyTargetOverride(cfg, []string{packageName}, target, cmd.Flags().Changed("target"))

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
				return o.SetupPackage(packageName, dotfiles, osName, manager)
			}

			return runWithSpinner(fmt.Sprintf("Setting up %s", packageName), func() error {
				return o.SetupPackage(packageName, dotfiles, osName, manager)
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print planned operations without executing them")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Answer yes to all prompts automatically")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Delete physical file conflicts without prompting")
	cmd.Flags().StringVarP(&targetDir, "target", "t", "$HOME", "Base target directory for stow operations")

	return cmd
}
