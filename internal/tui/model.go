package tui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/FranChesK0/tis-100/internal/emu"
	"github.com/FranChesK0/tis-100/internal/types"
)

type model struct {
	keys keyMap
	help help.Model

	running bool

	program *emu.Program
	puzzle  *types.Puzzle
}

func NewModel() *model {
	return &model{
		keys:    keys,
		help:    help.New(),
		program: emu.NewProgram(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("TIS-100")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd)

func (m model) View() string
