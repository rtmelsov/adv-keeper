// Package ui
package ui

import "github.com/charmbracelet/lipgloss"

var (
	Fg     = lipgloss.AdaptiveColor{Light: "#1F2328", Dark: "#C7D0D9"}
	Bg     = lipgloss.AdaptiveColor{Light: "#F6F8FA", Dark: "#0B141A"}
	Accent = lipgloss.AdaptiveColor{Light: "#0969DA", Dark: "#58A6FF"}
	Muted  = lipgloss.AdaptiveColor{Light: "#57606A", Dark: "#8B949E"}
	Danger = lipgloss.AdaptiveColor{Light: "#D1242F", Dark: "#FF6B6B"}

	Title   = lipgloss.NewStyle().Bold(true).Foreground(Accent)
	Label   = lipgloss.NewStyle().Foreground(Muted)
	LabelOn = Label.Foreground(Accent)

	Input = lipgloss.NewStyle().Padding(0, 1).
		Border(lipgloss.RoundedBorder()).BorderForeground(Muted)
	InputFocused = Input.BorderForeground(Accent)

	Btn = lipgloss.NewStyle().Padding(0, 2).
		Border(lipgloss.RoundedBorder()).BorderForeground(Muted)
	BtnFocused = Btn.Foreground(Accent).BorderForeground(Accent).Bold(true)

	ErrorText = lipgloss.NewStyle().Foreground(Danger)

	// Вот он — статус-бар: одна строка со спокойным цветом
	StatusBar = lipgloss.NewStyle().Foreground(Muted).MarginTop(1)

	Nav       = lipgloss.NewStyle().Foreground(Muted).Padding(0, 1)
	NavActive = lipgloss.NewStyle().Foreground(Accent).Bold(true).Padding(0, 1)

	TableBase = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(Muted).Padding(0, 1)
)
