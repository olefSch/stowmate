package orchestrator

import "github.com/stowmate/stowmate/internal/runner"

// ExecuteStow runs "stow -d <dotfilesPath> -R -t <target> <package>" so that
// the stow directory is set without requiring working-directory support from
// the runner.
func ExecuteStow(r runner.CommandRunner, dotfilesPath string, targetDir string, packageName string) error {
	return r.Run("stow", "-d", dotfilesPath, "-R", "-t", targetDir, packageName)
}

func executeStowDelete(r runner.CommandRunner, dotfilesPath string, targetDir string, packageName string) error {
	return r.Run("stow", "-d", dotfilesPath, "-D", "-t", targetDir, packageName)
}
