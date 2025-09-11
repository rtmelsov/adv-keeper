package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"os"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.Selected == "Vault" {
			switch msg.String() {
			case "esc":
				if m.table.Focused() {
					m.table.Blur()
				} else {
					m.table.Focus()
				}
			case "right", "l":
				if m.CursorHor < 1 {
					m.CursorHor++
				} else {
					m.CursorHor = 0
				}
			case "left", "h":
				if m.CursorHor > 0 {
					m.CursorHor--
				}
			case "enter":
				if m.CursorHor == 0 && m.SelectedFile == "" {
					if home, err := os.UserHomeDir(); err == nil && home != "" {
						m.FilePicker.CurrentDirectory = home
						m.FilePicker.Path = home
					} else if wd, _ := os.Getwd(); wd != "" {
						m.FilePicker.CurrentDirectory = wd
						m.FilePicker.Path = wd
					} else {
						m.FilePicker.CurrentDirectory = "/"
						m.FilePicker.Path = "/"
					}
					m.FilePicker.SetHeight(14) // чтобы было видно список
					// m.FilePicker.AllowedTypes = nil // не ставь []string{"*"}
					m.OpenFilePicker = true
					return m, m.FilePicker.Init() // ← важный момент
				} else if m.SelectedFile != "" {
					m.Loading = true
					// запускаем асинхронную команду
					return m, tea.Batch(
						m.Spinner.Tick,
						func() tea.Msg {
							_, err := akclient.UploadFile(m.SelectedFile)
							return struct{ err error }{err: err}
						},
					)
				} else {
					return m, tea.Batch(
						tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
					)
				}

			}
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
