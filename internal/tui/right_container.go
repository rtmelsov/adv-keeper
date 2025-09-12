package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) RightPane(inner string) string {
	title := ui.Title.Render(m.SelectedPage)
	return lipgloss.JoinVertical(lipgloss.Left, title, "", inner)
}
