package cli

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/olefSch/stowmate/internal/config"
	"github.com/olefSch/stowmate/internal/orchestrator"
)

// newRemoveCmd creates the `stowmate remove <name>` command.
func newRemoveCmd(ctx *cliContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Un-stow a named package and run post-remove hooks",
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

			if dryRun {
				return o.RemovePackage(packageName, dotfiles)
			}

			return runWithSpinner(fmt.Sprintf("Removing %s", packageName), func() error {
				return o.RemovePackage(packageName, dotfiles)
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print planned operations without executing them")
	cmd.Flags().StringVarP(&targetDir, "target", "t", "$HOME", "Base target directory for stow operations")

	return cmd
}
