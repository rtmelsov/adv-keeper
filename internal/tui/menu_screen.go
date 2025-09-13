package tui

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var MenuList = map[string]string{
	"Vault":    "Новый файл",
	"Main":     "Основное",
	"FileList": "Список файлов",
	"Logout":   "Выйти",
	"Login":    "Войти",
	"Register": "Регистрация",
}

var (
	stWelcomeTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E5E7EB"))
	stWelcomeSub   = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#9CA3AF"))
	stAvatar       = lipgloss.NewStyle().Width(3).Align(lipgloss.Center).Bold(true).
			Foreground(lipgloss.Color("#063E3B")).Background(lipgloss.Color("#34D399")).MarginRight(1)
	stSep = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
)

func firstInitial(name string) string {
	if name == "" {
		return "?"
	}
	r, _ := utf8.DecodeRuneInString(name)
	return strings.ToUpper(string(r))
}

func (m TuiModel) Menu(innerW int) string {
	if innerW < 10 {
		innerW = 10
	}

	// --- Header: аватар + «Добро пожаловать»
	name := "гость"
	if m.Profile.Auth && m.Profile.Email != "" {
		name = m.Profile.Email
	}
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		stAvatar.Render(firstInitial(name)),
		lipgloss.JoinVertical(lipgloss.Left,
			stWelcomeSub.Render("Добро пожаловать,"),
			stWelcomeTitle.Render(name),
		),
	)

	// --- Точная линия по ширине сайдбара
	hr := stSep.Width(innerW).Render(strings.Repeat("─", innerW))

	// --- Пункты меню: зелёная заливка на всю ширину, БЕЗ рамки
	var choices []string
	if m.Profile.Auth {
		choices = m.MainChoices
	} else {
		choices = m.LoginChoices
	}

	rowBase := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Left)
	rowInactive := rowBase.PaddingLeft(1).Foreground(lipgloss.Color("#9CA3AF"))
	rowActive := rowBase.PaddingLeft(1).Bold(true).
		Foreground(lipgloss.Color("#0B1021")).Background(lipgloss.Color("#A7F3D0")) // полный зелёный фон

	items := make([]string, 0, len(choices))
	for i, choice := range choices {
		active := m.LeftCursor == i && m.HorCursor == 0
		label := "› " + MenuList[choice]
		if active {
			items = append(items, rowActive.Render(label))
		} else {
			items = append(items, rowInactive.Render(label))
		}
	}
	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	// --- Сборка
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header, // вернули хедер
		hr,     // линия ровно innerW
		"",
		list, // пункты
	)
}
