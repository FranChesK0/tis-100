package tui

import tea "github.com/charmbracelet/bubbletea"

func ProgramRun() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
