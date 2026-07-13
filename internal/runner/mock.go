package runner

import (
	"errors"
	"fmt"
)

// Invocation records a single command executed by MockRunner.
type Invocation struct {
	Name string
	Args []string
}

// MockRunner is a test double that records command invocations and returns
// configurable responses.
type MockRunner struct {
	Calls     []Invocation
	Responses map[string]error
	Outputs   map[string]string
	Verbose   bool
	Logger    func(string)
}

// NewMockRunner returns a MockRunner with initialized maps.
func NewMockRunner() *MockRunner {
	return &MockRunner{
		Calls:     []Invocation{},
		Responses: map[string]error{},
		Outputs:   map[string]string{},
		Logger:    func(msg string) { fmt.Println(msg) },
	}
}

// Run records the invocation and returns the configured response.
func (m *MockRunner) Run(name string, args ...string) error {
	m.logCommand(name, args)
	m.Calls = append(m.Calls, Invocation{Name: name, Args: append([]string(nil), args...)})

	if err, ok := m.Responses[commandKey(name, args)]; ok {
		return err
	}
	if err, ok := m.Responses[name]; ok {
		return err
	}
	return nil
}

// RunWithOutput records the invocation and returns the configured output.
func (m *MockRunner) RunWithOutput(name string, args ...string) (string, error) {
	m.logCommand(name, args)
	m.Calls = append(m.Calls, Invocation{Name: name, Args: append([]string(nil), args...)})

	key := commandKey(name, args)
	if out, ok := m.Outputs[key]; ok {
		if err, hasErr := m.Responses[key]; hasErr {
			return out, err
		}
		return out, nil
	}
	if out, ok := m.Outputs[name]; ok {
		if err, hasErr := m.Responses[name]; hasErr {
			return out, err
		}
		return out, nil
	}
	if err, ok := m.Responses[key]; ok {
		return "", err
	}
	if err, ok := m.Responses[name]; ok {
		return "", err
	}
	return "", nil
}

// commandKey returns a deterministic string key for a command and its args.
func commandKey(name string, args []string) string {
	if len(args) == 0 {
		return name
	}
	return fmt.Sprintf("%s %s", name, joinArgs(args))
}

func joinArgs(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}
	return out
}

// SetResponse configures the error returned when a command matching the key is
// executed. The key may be a binary name or a full "name arg1 arg2" string.
func (m *MockRunner) SetResponse(key string, err error) {
	if m.Responses == nil {
		m.Responses = map[string]error{}
	}
	m.Responses[key] = err
}

// SetOutput configures the stdout returned when a command matching the key is
// executed.
func (m *MockRunner) SetOutput(key string, output string) {
	if m.Outputs == nil {
		m.Outputs = map[string]string{}
	}
	m.Outputs[key] = output
}

// SetVerbose enables command logging through the configured logger.
func (m *MockRunner) SetVerbose(logger func(string)) {
	m.Verbose = true
	m.Logger = logger
}

func (m *MockRunner) logCommand(name string, args []string) {
	if m.Verbose && m.Logger != nil {
		m.Logger(joinArgs(append([]string{name}, args...)))
	}
}

// ErrSimulated is a sentinel error returned by mocks when no specific error is
// configured.
var ErrSimulated = errors.New("simulated command failure")
