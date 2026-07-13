package orchestrator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompter abstracts user confirmation prompts. Milestone 3 will replace the
// SimplePrompter implementation with a Charmbracelet/huh based prompter.
type Prompter interface {
	Confirm(message string) (bool, error)
}

// SimplePrompter reads Y/n confirmations from stdin. Empty input defaults to
// yes when the message contains "(Y/n)" and defaults to no for "(y/N)".
type SimplePrompter struct{}

// Confirm prompts the user and returns true for yes.
func (SimplePrompter) Confirm(message string) (bool, error) {
	fmt.Print(message + " ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response := strings.ToLower(strings.TrimSpace(line))
	if response == "" {
		return strings.Contains(message, "(Y/n)"), nil
	}
	return response == "y" || response == "yes", nil
}

// MockPrompter returns configured responses sequentially. It is intended for
// tests only.
type MockPrompter struct {
	Responses []bool
	idx       int
}

// Confirm returns the next configured response.
func (m *MockPrompter) Confirm(message string) (bool, error) {
	if m.idx >= len(m.Responses) {
		return false, fmt.Errorf("no more configured responses")
	}
	resp := m.Responses[m.idx]
	m.idx++
	return resp, nil
}
