package tui

import tea "github.com/charmbracelet/bubbletea"

type Model struct{}

func NewModel() *Model

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("TIS-100")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)

func (m Model) View() string
