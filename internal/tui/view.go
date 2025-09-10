// Package tui
package tui

import (
	"fmt"
	"github.com/rtmelsov/adv-keeper/internal/ui"
)

func (m TuiModel) View() string {

	if m.Loading {
		return ui.StatusBar.Render(m.Spin.View()+" Загрузка…") + fmt.Sprintf("chunk size: %v - %v", m.LoaderCount.FileSize, m.LoaderCount.ChankSize) + "\n"
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

	before := ""
	if m.Error != "" {
		before = fmt.Sprintf("Error: %s\n", m.Error)
	}
	if m.Profile.Auth {
		before += fmt.Sprintf("Profile: %s\n", m.Profile.Email)
	} else {
		before += "Profile: ---\n"
	}

	nav := " [MAIN]   [FAILS] \n"

	if m.Cursor == 0 {
		nav = "<[MAIN]>  [FAILS] \n"
	}
	if m.CursorHor == 1 {
		nav = " [MAIN]  <[FAILS]>\n"
	}
	return before + nav + s
}
