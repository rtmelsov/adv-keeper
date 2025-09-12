package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	// "github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) Main() string {
	title := ui.Title.Render(fmt.Sprintf("Hello: %s", m.Profile.Email))
	var s string
	s += title + "\n\n"

	btn := lipgloss.JoinVertical(
		lipgloss.Top,
		ui.ButtonInactive.Render("Посмотреть файлы"),
		ui.ButtonInactive.Render("Загрузить файл"),
	)
	if m.HorCursor == 0 {
		return s + btn
	}
	if m.RightCursor == 0 {
		btn = lipgloss.JoinVertical(
			lipgloss.Top,
			ui.ButtonActive.Render("Посмотреть файлы"),
			ui.ButtonInactive.Render("Загрузить файл"),
		)
	} else {
		btn = lipgloss.JoinVertical(
			lipgloss.Top,
			ui.ButtonInactive.Render("Посмотреть файлы"),
			ui.ButtonActive.Render("Загрузить файл"),
		)
	}
	return s + btn
}

func (m TuiModel) MainAction(msg string) (tea.Model, tea.Cmd) {
	switch msg {
	case "esc":
		m.HorCursor = 0
		m.RightCursor = 0
		return m, tea.ClearScreen
	case "enter":
		if m.RightCursor == 0 {
			m.Loading = true
			return m, tea.Batch(
				m.Spin.Tick,
				func() tea.Msg {
					list, err := akclient.GetFiles()
					return getListFinishedMsg{err: err, list: list}
				},
			)
		} else {
			m.SelectedPage = "Vault"
		}
		m.RightCursor = 0
		return m, tea.ClearScreen

	default:
		return m, nil
	}
}
