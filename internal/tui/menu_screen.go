package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) Menu() string {
	title := ui.Title.Render("todo: надо добавить инфо клиента")

	if m.Profile.Auth {
		m.Choices = m.MainChoices
	} else {
		m.Choices = m.LoginChoices
	}

	items := make([]string, 0, len(m.Choices))
	for i, choice := range m.Choices {
		item := ui.NavInactive.Render(fmt.Sprintf("[%s]", choice))
		if m.LeftCursor == i && m.HorCursor == 0 {
			item = ui.NavActive.Render(fmt.Sprintf("[%s]", choice))
		}
		items = append(items, item)
	}

	list := lipgloss.JoinVertical(lipgloss.Top, items...)
	return lipgloss.JoinVertical(lipgloss.Left, title, "", list)
}
