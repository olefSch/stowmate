package orchestrator

import "github.com/spf13/afero"

// PreClean deletes every path listed in the package's pre_clean configuration.
// Paths are expected to be absolute and expanded before calling PreClean.
func PreClean(fs afero.Fs, paths []string) error {
	for _, p := range paths {
		if err := fs.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}
