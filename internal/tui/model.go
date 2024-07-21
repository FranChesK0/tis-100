package tui

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/FranChesK0/tis-100/internal/emu"
	"github.com/FranChesK0/tis-100/internal/parser"
	"github.com/FranChesK0/tis-100/internal/types"
)

type model struct {
	keys          keyMap
	help          help.Model
	filepicker    filepicker.Model
	filepickerErr error
	puzzlePath    string
	running       bool
	program       *emu.Program
	puzzle        *types.Puzzle
}

type clearErrMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg { return clearErrMsg{} })
}

func NewModel() (*model, error) {
	var err error

	fp := filepicker.New()
	fp.AllowedTypes = []string{".lua"}
	fp.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		return &model{}, nil
	}
	return &model{
		keys:       keys,
		help:       help.New(),
		filepicker: fp,
		program:    emu.NewProgram(),
	}, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("TIS-100"), m.filepicker.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.puzzlePath == "" {
		return m.updateFilepicker(msg)
	}
	if m.puzzle == nil {
		m.puzzle, _ = parser.FetchPuzzle(m.puzzlePath) // TODO: add checks for errors
	}

	return m, nil
}

func (m model) View() string {
	if m.puzzlePath == "" {
		return m.viewFilepicker()
	}
	if m.puzzle == nil {
		return "Loading puzzle..."
	}
	return "Puzzle title: " + m.puzzle.Title
}

func (m model) updateFilepicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case clearErrMsg:
		m.filepickerErr = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.puzzlePath = path
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.filepickerErr = errors.New(path + " is not valid")
		m.puzzlePath = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) viewFilepicker() string {
	view := "\n "
	if m.filepickerErr != nil {
		view += m.filepicker.Styles.DisabledFile.Render(m.filepickerErr.Error())
	} else if m.puzzlePath == "" {
		view += "Pick a puzzle:"
	} else {
		view += "Selected puzzle: " + m.filepicker.Styles.Selected.Render(m.puzzlePath)
	}
	view += "\n\n" + m.filepicker.View() + "\n"
	return view
}
