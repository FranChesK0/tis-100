package tui

import (
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type clearErrMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg { return clearErrMsg{} })
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
