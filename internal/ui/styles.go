// Package ui
package ui

import "github.com/charmbracelet/lipgloss"

// Палитра (адаптивная под светлую/тёмную тему)
var (
	ColBorder = lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#3B3F51"}

	// Outer app box (the centered 700×700-ish container)
	AppBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colBorder).
		Padding(1, 2)

	// Left pane (menu) — only a right border, so it looks like a splitter
	Sidebar = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colBorder).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(true).
		Padding(0, 1)

	// Right pane (content)
	Content = lipgloss.NewStyle().
		Padding(0, 1)

	// Optional: table box look
	Table = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colBorder).
		Padding(0, 1)

	colText   = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#E6E6E6"}
	colMuted  = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#A1A1AA"} // серый
	colAccent = lipgloss.AdaptiveColor{Light: "#0EA5E9", Dark: "#7DD3FC"} // голубой
	colPanel  = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#0F1117"} // фон «карточки»
	colBorder = lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#3B3F51"}
	colError  = lipgloss.AdaptiveColor{Light: "#B91C1C", Dark: "#F87171"}
)

// Заголовки/текст
var (
	Title = lipgloss.NewStyle().Bold(true).Foreground(colText)
	Muted = lipgloss.NewStyle().Faint(true).Foreground(colMuted)
)

// Вкладки (табы) — верхняя навигация
var (
	TabActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1).
			MarginRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colBorder)

	TabInactive = lipgloss.NewStyle().
			Foreground(colMuted).
			Padding(0, 1).
			MarginRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colBorder)
)

// Пункты навигации/меню (в списках)
var (
	NavActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)

	NavInactive = lipgloss.NewStyle().
			Foreground(colAccent).
			Underline(true)
)

// Формы и кнопки
var (
	FieldLabel       = lipgloss.NewStyle().Foreground(colMuted)
	FieldLabelActive = lipgloss.NewStyle().Foreground(colAccent).Bold(true)
	ButtonInactive   = lipgloss.NewStyle().Foreground(colText).Padding(0, 2).Border(lipgloss.NormalBorder()).BorderForeground(colBorder)
	ButtonActive     = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 2).Border(lipgloss.RoundedBorder()).BorderForeground(colBorder)
)

// Контейнеры/таблицы/статус
var (
	Box       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colBorder).Padding(1, 2)
	Error     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(colError).Padding(0, 1)
	StatusBar = lipgloss.NewStyle().Background(colAccent).Foreground(lipgloss.Color("#000000")).Padding(0, 1)
	MetaStyle = lipgloss.NewStyle().Faint(true)
)

// Контент справа (disabled)
var ContentDisabled = Content.Copy().
	Foreground(colMuted).
	Faint(true)
