// Package tui
package tui

import (
	"fmt"
)

func (m TuiModel) View() string {
	if m.Loading {
		return "LOADING..."
	}
	s := "CAN'T FIND THE PAGE"

	switch m.Selected {
	case "Vault":
		s = m.Vault()
	case "Menu":
		s = m.Menu()
	case "Register":
		s = m.Register()
	case "Login":
		s = "Login"
	}

	s += "\nm.Selected: "
	s += m.Selected + "\n"
	s += "\n↑/k ↓/j — навигация, ␣/Enter — выбрать, a — выделить всё, x — удалить, q — выход.\n"

	before := ""
	if m.Error != "" {
		before = fmt.Sprintf("Error: %s\n", m.Error)
	}
	if m.Profile.Auth {
		before += fmt.Sprintf("Profile: %s\n", m.Profile.Email)
	} else {
		before += "Profile: ---\n"
	}

	nav := "<[MAIN]>  [FAILS] \n"
	if m.CursorHor == 1 {
		nav = " [MAIN]  <[FAILS]>\n"
	}
	return before + nav + s
}
