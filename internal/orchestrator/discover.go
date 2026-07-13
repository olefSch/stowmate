package orchestrator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

// DiscoverPackages returns the top-level directory names in the dotfiles path.
// Hidden directories (names starting with ".") and non-directory entries are
// excluded.
func DiscoverPackages(fs afero.Fs, dotfilesPath string) ([]string, error) {
	info, err := fs.Stat(dotfilesPath)
	if err != nil {
		return nil, fmt.Errorf("dotfiles path %q: %w", dotfilesPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("dotfiles path %q is not a directory", dotfilesPath)
	}

	entries, err := afero.ReadDir(fs, dotfilesPath)
	if err != nil {
		return nil, err
	}

	packages := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		packages = append(packages, entry.Name())
	}

	sort.Strings(packages)
	return packages, nil
}
