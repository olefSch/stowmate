package main

import (
	"fmt"
	"os"

	"github.com/olefSch/stowmate/internal/cli"
)

var (
	version = "dev"
	commit  string
	date    string
)

func main() {
	if err := cli.Execute(formatVersion(version, commit, date)); err != nil {
		os.Exit(1)
	}
}

func formatVersion(version, commit, date string) string {
	v := version
	if commit != "" && date != "" {
		v = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	} else if commit != "" {
		v = fmt.Sprintf("%s (%s)", version, commit)
	} else if date != "" {
		v = fmt.Sprintf("%s (%s)", version, date)
	}
	return v
}
