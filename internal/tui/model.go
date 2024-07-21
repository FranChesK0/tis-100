package tui

import (
	"os"

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
