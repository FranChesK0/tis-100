package tui

import tea "github.com/charmbracelet/bubbletea"

func ProgramRun() error {
	m, err := NewModel()
	if err != nil {
		return err
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
