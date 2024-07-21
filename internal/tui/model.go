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
	filepicker filepicker.Model
	help       help.Model

	puzzle  *types.Puzzle
	program *emu.Program

	keys       keyMap
	puzzlePath string
	running    bool

	filepickerErr  error
	fetchPuzzleErr error
}

func NewModel() (*model, error) {
	var err error

	fp := filepicker.New()
	fp.AllowedTypes = []string{".lua"}
	fp.CurrentDirectory, err = os.Getwd()
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
	} else if m.puzzle == nil {
		m.puzzle, m.fetchPuzzleErr = parser.FetchPuzzle(m.puzzlePath)
	}

	return m, nil
}

func (m model) View() string {
	if m.puzzlePath == "" {
		return m.viewFilepicker()
	} else if m.puzzle == nil {
		return "Loading puzzle..."
	} else if m.fetchPuzzleErr != nil {
		return "Error while fetching puzzle: " + m.fetchPuzzleErr.Error()
	} else {
		return "Puzzle title: " + m.puzzle.Title
	}
}
