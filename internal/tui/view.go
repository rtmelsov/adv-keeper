// Package tui
package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/rtmelsov/adv-keeper/internal/ui"
)

var (
	// Base box style (~700×700 px ≈ 88×44 cells; tweak for your font)
	SidebarW = 24
)

func (m TuiModel) LoadingFullScreen() string {
	// build the colored loader content (see next section)
	content := m.loadingViewStyled(min(m.W-8, 72)) // clamp width so it looks neat

	panel := ui.StPanel.
		Width(min(m.W-6, 76)). // panel width
		MaxWidth(m.W - 6).
		Render(content)

	// center panel within the whole terminal
	return lipgloss.Place(
		m.W, m.H,
		lipgloss.Center, lipgloss.Center,
		panel,
	)
}

func (m TuiModel) View() string {
	// Wait for a size before we try to center things
	if m.W == 0 || m.H == 0 {
		return ""
	}

	// --- Header/status (outside the box) ---
	status := ""
	if m.Loading {
		if m.StreamLoading {
			return m.LoadingFullScreen()
		} else {
			status = ui.StatusBar.Render(m.Spin.View()+" Загрузка…") +
				fmt.Sprintf(" chunk size: %v - %v", m.LoaderCount.FileSize, m.LoaderCount.ChankSize)
		}

	} else {
		status = ui.StatusBar.Render("Готово")
	}

	// --- Meta (error/profile) ---
	before := ""
	if m.Error != "" {
		before += ui.Error.Render(fmt.Sprintf("Error: %s", m.Error)) + "\n"
	}
	// --- Page content ---
	s := "CAN'T FIND THE PAGE"
	switch m.SelectedPage {
	case "Vault":
		s = m.Vault()
	case "Main":
		s = m.Main()
	case "FileDetails":
		s = m.FileDetails()
	case "FileList":
		s = m.table.View()
	case "Register":
		s = m.Register()
	case "Login":
		s = m.Login()
	}

	inner := lipgloss.JoinVertical(lipgloss.Left, before, s)

	// --- Box sizing & centering ---
	boxW, boxH := 88, 44 // ≈ 700×700 px at ~8×16px cells
	if boxW > m.W {
		boxW = m.W
	}
	// Reserve one line for the status bar above
	availH := m.H - 1
	if availH < 1 {
		availH = 1
	}
	if boxH > availH {
		boxH = availH
	}

	rightStyle := ui.Content
	if m.HorCursor == 0 || m.Loading {
		rightStyle = ui.ContentDisabled
	}

	// --- left/right widths ---
	leftW := SidebarW
	if leftW > boxW-10 {
		leftW = boxW / 3
	}
	rightW := boxW - leftW - 1

	frameW, _ := ui.Sidebar.GetFrameSize() // бордеры+паддинги
	innerLeftW := leftW - frameW
	if innerLeftW < 1 {
		innerLeftW = 1
	}

	frameRightW, _ := rightStyle.GetFrameSize() // бордеры/паддинги правой колонки
	innerRightW := rightW - frameRightW
	if innerRightW < 1 {
		innerRightW = 1
	}
	// --- left menu & right content ---
	left := ui.Sidebar.
		Width(leftW).
		Height(boxH - 2).
		Render(m.Menu(innerLeftW)) // <── передаём inner width
	rightContent := m.RightPane(inner) // table OR upload (see below)
	right := rightStyle.
		Width(rightW).
		Height(boxH - 2).
		Render(rightContent)

	// --- join panes and wrap into outer box ---
	row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	app := ui.AppBox.
		Width(boxW).
		Height(boxH).
		Render(row)

	centered := lipgloss.Place(m.W, availH, lipgloss.Center, lipgloss.Center, app)

	return status + "\n" + centered
}
