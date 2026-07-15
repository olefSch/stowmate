package cli

import (
	"github.com/charmbracelet/huh"
	"github.com/olefSch/stowmate/internal/orchestrator"
)

// HuhPrompter implements orchestrator.Prompter using charmbracelet/huh for
// interactive Y/N confirmations.
type HuhPrompter struct{}

// Confirm asks the user a yes/no question and returns the response.
func (HuhPrompter) Confirm(message string) (bool, error) {
	var confirmed bool
	err := huh.NewConfirm().
		Title(message).
		Value(&confirmed).
		Run()
	return confirmed, err
}

var _ orchestrator.Prompter = HuhPrompter{}
