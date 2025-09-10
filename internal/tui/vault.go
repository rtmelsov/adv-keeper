package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rtmelsov/adv-keeper/internal/akclient"

	"os"
)

type uploadFinishedMsg struct{ err error }

func (m TuiModel) ReturnVault(msg string) (tea.Model, tea.Cmd) {
	switch msg {
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
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < 1 {
			m.Cursor++
		}

	case "enter":
		if m.CursorHor == 0 {
			if m.Cursor == 1 && m.SelectedFile == "" {
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
				return m, tea.Batch(
					m.Spin.Tick,
					func() tea.Msg {
						_, err := akclient.UploadFile(m.SelectedFile)
						return uploadFinishedMsg{err: err}
					},
				)
			}
		}

		if m.CursorHor == 1 {
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	return m, nil
}
