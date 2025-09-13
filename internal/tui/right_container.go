package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) RightPane(inner string) string {
	// грубо берём доступную ширину по контенту; можно заменить на реальную innerRightW
	w := lipgloss.Width(inner)
	if w < 24 {
		w = 24
	}

	// крупный заголовок по центру
	title := ui.Title.
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render(MenuList[m.SelectedPage])

	// тонкая черта под тайтлом
	rule := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Width(w).
		Render(strings.Repeat("─", w))

	return lipgloss.JoinVertical(lipgloss.Left, title, rule, "", inner)
}
