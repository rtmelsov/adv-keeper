package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"path/filepath"
)

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.OpenFilePicker {
		var cmd tea.Cmd
		m.FilePicker, cmd = m.FilePicker.Update(msg)

		if ok, path := m.FilePicker.DidSelectFile(msg); ok {
			m.SelectedFile = filepath.Clean(path)
			m.OpenFilePicker = false
			return m, tea.Batch(cmd, tea.ClearScreen)
		}
		if km, ok := msg.(tea.KeyMsg); ok && (km.String() == "esc" || km.String() == "q") {
			m.OpenFilePicker = false
			return m, tea.Batch(cmd, tea.ClearScreen)
		}
		return m, cmd
	}

	// 2) сначала ловим завершение фоновой задачи
	switch msg := msg.(type) {
	case uploadFinishedMsg:
		m.Loading = false
		if msg.err != nil {
			m.Error = msg.err.Error()
		} else {
			m.Error = ""
			m.SelectedFile = ""
		}
		return m, nil

	// 3) тики спиннера: только когда Loading=true
	case spinner.TickMsg:
		if m.Loading {
			var cmd tea.Cmd
			m.Spin, cmd = m.Spin.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.Selected == "Vault" {
			return m.ReturnVault(msg.String())
		}
		if m.Selected == "Register" && m.InputFocused {
			if msg.String() == "esc" {
				m.login.Blur()
				m.password.Blur()
				m.InputFocused = false
				return m, nil
			}
			if msg.String() == "enter" {
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
		switch msg.String() {
		case "esc":
			if len(m.History) <= 0 {
				return m, tea.Quit
			}
			m.Selected = m.History[len(m.History)-1]
			m.History = m.History[:len(m.History)-1]
			m = *selectedScreen(&m)
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter": // пробел тоже работает
			m.History = append(m.History, m.Selected)
			m = *selectedScreen(&m)

			return m, tea.ClearScreen
		}
	}
	return m, nil
}

func selectedScreen(m *TuiModel) *TuiModel {
	switch m.Selected {
	case "Menu":
		m.MenuAction()
	case "Register":
		m.RegisterAction()
	case "Vault":
		m.Vault()
	}
	return m
}
