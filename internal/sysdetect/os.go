package sysdetect

import (
	"fmt"
	"runtime"
)

// DetectOS returns the normalized host operating system name. Only darwin and
// linux are supported.
func DetectOS() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "darwin", nil
	case "linux":
		return "linux", nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
