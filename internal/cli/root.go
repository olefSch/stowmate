package cli

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/olefSch/stowmate/internal/config"
	"github.com/olefSch/stowmate/internal/runner"
)

var (
	dotfilesPath string
	targetDir    string
	verbose      bool
	dryRun       bool
	yes          bool
	force        bool
)

// newRootCmd creates the root stowmate command and attaches all subcommands.
func newRootCmd(ctx *cliContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "stowmate",
		Short: "A friendly dotfiles manager powered by GNU Stow",
		Long: `Stowmate discovers packages in your dotfiles directory and uses GNU Stow
to symlink them into your home directory.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				ctx.logger.SetLevel(log.DebugLevel)
				if r, ok := ctx.runner.(*runner.ExecRunner); ok {
					r.Verbose = true
					r.Logger = func(msg string) { ctx.logger.Debug(msg) }
				}
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.PersistentFlags().StringVarP(&dotfilesPath, "dotfiles", "d", "$HOME/dotfiles", "Path to the dotfiles directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable debug logging and print shell commands before execution")

	rootCmd.AddCommand(newRunCmd(ctx))
	rootCmd.AddCommand(newPackageCmd(ctx))
	rootCmd.AddCommand(newRemoveCmd(ctx))

	return rootCmd
}

// expandPath expands leading $HOME or ~ in a path and falls back to the raw
// value if home directory resolution fails.
func expandPath(path string) string {
	expanded := config.ExpandHome(path)
	if expanded != "" {
		return expanded
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return home
}
