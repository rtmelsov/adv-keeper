package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	// "github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) FileDetails() string {
	title := ui.Title.Render(fmt.Sprintf("Информация файла: %s", m.SelectedFileInfo.Filename))
	var s string
	s += title + "\n\n"

	btn := lipgloss.JoinVertical(
		lipgloss.Top,
		ui.ButtonInactive.Render("Скачать файл"),
		ui.ButtonInactive.Render("Удалить файл"),
	)
	if m.HorCursor == 0 {
		return s + btn
	}
	if m.RightCursor == 0 {
		btn = lipgloss.JoinVertical(
			lipgloss.Top,
			ui.ButtonActive.Render("Скачать файл"),
			ui.ButtonInactive.Render("Удалить файл"),
		)
	} else {
		btn = lipgloss.JoinVertical(
			lipgloss.Top,
			ui.ButtonInactive.Render("Скачать файл"),
			ui.ButtonActive.Render("Удалить файл"),
		)
	}
	return s + btn
	// return "FILE ID: " +

}

func (m TuiModel) FileDetailsAction(msg string) (tea.Model, tea.Cmd) {
	switch msg {
	case "esc":
		m.HorCursor = 0
		m.RightCursor = 0
		return m, tea.ClearScreen
	case "enter":
		if m.RightCursor == 0 {
			_, err := akclient.DownloadFile(m.SelectedFileInfo.Fileid)
			if err != nil {
				m.Error = err.Error()
			}
		}
		if m.RightCursor == 1 {
			err := akclient.DeleteFile(m.SelectedFileInfo.Fileid)
			if err != nil {
				m.Error = err.Error()
			}
		}
		return m, tea.ClearScreen

	default:
		return m, nil
	}
}
