package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

// DetectConflicts scans the target directory for physical files or directories
// that match paths inside the package. Symlinks are not treated as conflicts.
func DetectConflicts(fs afero.Fs, dotfilesPath string, targetDir string, packageName string) ([]string, error) {
	pkgPath := filepath.Join(dotfilesPath, packageName)
	exists, err := afero.Exists(fs, pkgPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("package %q not found at %s", packageName, pkgPath)
	}

	var conflicts []string
	err = afero.Walk(fs, pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == pkgPath {
			return nil
		}
		if info.IsDir() {
			hasSubdir, err := dirHasSubdir(fs, path)
			if err != nil {
				return err
			}
			if hasSubdir {
				return nil
			}
		}

		rel, err := filepath.Rel(pkgPath, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, rel)
		if isPhysicalPath(fs, targetPath) {
			conflicts = append(conflicts, targetPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return filterAncestorConflicts(conflicts), nil
}

// ResolveConflicts prompts the user for each conflict and deletes it when
// confirmed. If force is true, all conflicts are deleted without prompting.
func ResolveConflicts(fs afero.Fs, conflicts []string, force bool, prompter Prompter) error {
	for _, conflict := range conflicts {
		if !force {
			ok, err := prompter.Confirm(fmt.Sprintf("Delete conflicting path %q? (y/N)", conflict))
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("conflict not resolved: %s", conflict)
			}
		}
		if err := fs.RemoveAll(conflict); err != nil {
			return err
		}
	}
	return nil
}

func filterAncestorConflicts(conflicts []string) []string {
	// Process deepest paths first so a directory is skipped when one of its
	// descendant files is already a conflict.
	sort.Slice(conflicts, func(i, j int) bool {
		return len(conflicts[i]) > len(conflicts[j])
	})

	var filtered []string
	for _, c := range conflicts {
		ancestor := false
		for _, existing := range filtered {
			if strings.HasPrefix(existing, c+string(filepath.Separator)) {
				ancestor = true
				break
			}
		}
		if !ancestor {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func dirHasSubdir(fs afero.Fs, path string) (bool, error) {
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func isPhysicalPath(fs afero.Fs, path string) bool {
	fi, err := lstat(fs, path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink == 0
}

func lstat(fs afero.Fs, path string) (os.FileInfo, error) {
	if lfs, ok := fs.(afero.Lstater); ok {
		fi, _, err := lfs.LstatIfPossible(path)
		return fi, err
	}
	return fs.Stat(path)
}
