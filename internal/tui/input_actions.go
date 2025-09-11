package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m TuiModel) InputActions(msg string) (tea.Model, tea.Cmd) {
	if msg == "esc" || msg == "enter" {
		m.login.Blur()
		m.password.Blur()
		m.InputFocused = false
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	if m.password.Focused() {
		m.password, cmd = m.password.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.login.Focused() {
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
