package cli

import (
	"os"

	"github.com/charmbracelet/log"
)

// newLogger returns a colorful charmbracelet/log logger writing to stderr.
func newLogger(verbose bool) *log.Logger {
	level := log.InfoLevel
	if verbose {
		level = log.DebugLevel
	}
	return log.NewWithOptions(os.Stderr, log.Options{
		Level: level,
	})
}
