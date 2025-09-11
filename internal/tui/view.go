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

func (m TuiModel) View() string {
	// Wait for a size before we try to center things
	if m.W == 0 || m.H == 0 {
		return ""
	}

	// --- Header/status (outside the box) ---
	status := ""
	if m.Loading {
		status = ui.StatusBar.Render(m.Spin.View()+" Загрузка…") +
			fmt.Sprintf(" chunk size: %v - %v", m.LoaderCount.FileSize, m.LoaderCount.ChankSize)
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

	// --- left/right widths ---
	leftW := SidebarW
	if leftW > boxW-10 { // keep content usable on small screens
		leftW = boxW / 3
	}
	rightW := boxW - leftW - 1 // -1 accounts for AppBox padding/borders visually; tweak if needed
	if rightW < 10 {
		rightW = 10
	}

	rightStyle := ui.Content
	if m.HorCursor == 0 {
		rightStyle = ui.ContentDisabled
	}

	// --- left menu & right content ---
	left := ui.Sidebar.
		Width(leftW).
		Height(boxH - 2). // compensate AppBox padding; remove if you prefer natural height
		Render(m.Menu())

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
