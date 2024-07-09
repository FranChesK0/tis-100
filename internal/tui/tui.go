package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/FranChesK0/tis-100/internal/tui/models"
)

func ProgramRun() error {
	p := tea.NewProgram(models.NewMainModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
