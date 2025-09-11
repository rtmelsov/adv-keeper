package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/ui"

	"os"
)

type uploadFinishedMsg struct{ err error }

func (m TuiModel) VaultAction(msg string) (tea.Model, tea.Cmd) {
	switch msg {
	case "esc":
		m.HorCursor = 0
		m.RightCursor = 0
		return m, tea.ClearScreen

	case "enter":
		if m.SelectedFile == "" {
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
	return m, nil
}

func (m TuiModel) Vault() string {

	s := ui.Title.Render("try to add some file\n\n")
	if m.OpenFilePicker {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			m.FilePicker.View(),
		)
	} else if m.SelectedFile != "" {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			ui.Title.Render(fmt.Sprintf("\nВы выбрали файл: %s\n", m.SelectedFile)),
		)
	} else {
		btn := ui.ButtonInactive.Render("ADD FILE")
		if m.HorCursor == 1 {
			btn = ui.ButtonActive.Render("ADD FILE")
		}
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			btn,
		)
	}
}
