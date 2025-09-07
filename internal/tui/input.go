package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m TuiModel) InputActions(msg string) (tea.Model, tea.Cmd) {
	if msg == "esc" {
		m.login.Blur()
		m.password.Blur()
		m.InputFocused = false
		return m, nil
	}
	if msg == "enter" {
		if m.Cursor == 2 {
			m.InputFocused = false
			m.password.Blur()
			return m, nil
		}
		m.login.Blur()
		m.password.Focus()
		m.Cursor++
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.Cursor == 0 {
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Cursor == 1 {
		m.password, cmd = m.password.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
