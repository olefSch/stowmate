package cli

import (
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
)

// runWithSpinner runs fn while displaying a simple terminal spinner. If stdout
// is not a TTY, the spinner is skipped and fn is executed directly so tests and
// non-interactive environments are not affected.
func runWithSpinner(label string, fn func() error) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return fn()
	}

	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	p := tea.NewProgram(spinnerModel{spinner: s, label: label, fn: fn})
	m, err := p.Run()
	if err != nil {
		return err
	}
	return m.(spinnerModel).err
}

type spinnerModel struct {
	spinner spinner.Model
	label   string
	fn      func() error
	done    bool
	err     error
}

type runDoneMsg struct{ err error }

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		return runDoneMsg{err: m.fn()}
	})
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case runDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m spinnerModel) View() string {
	if m.done {
		return ""
	}
	return m.spinner.View() + " " + m.label
}
