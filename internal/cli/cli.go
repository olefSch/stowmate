package cli

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stowmate/stowmate/internal/orchestrator"
	"github.com/stowmate/stowmate/internal/runner"
)

// cliContext holds the runtime dependencies for the CLI. Tests construct this
// directly so commands can be exercised with doubles.
type cliContext struct {
	runner   runner.CommandRunner
	fs       afero.Fs
	prompter orchestrator.Prompter
	logger   *log.Logger
}

// Execute is the application entry point. It wires production dependencies and
// runs the Cobra command tree.
func Execute(version string) error {
	ctx := &cliContext{
		runner:   runner.NewExecRunner(false, nil),
		fs:       afero.NewOsFs(),
		prompter: HuhPrompter{},
		logger:   newLogger(false),
	}
	rootCmd := newRootCmd(ctx)
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		ctx.logger.Error(err.Error())
	}
	return err
}
