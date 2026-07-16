package orchestrator

import (
	"fmt"
	"os"

	"github.com/olefSch/stowmate/internal/runner"
)

// ExecuteHooks runs hook scripts sequentially. Errors are logged but do not
// stop execution of remaining hooks.
func ExecuteHooks(r runner.CommandRunner, hooks []string) error {
	for _, hook := range hooks {
		if err := r.Run("sh", "-c", hook); err != nil {
			fmt.Fprintf(os.Stderr, "hook failed: %v\n", err)
		}
	}
	return nil
}
