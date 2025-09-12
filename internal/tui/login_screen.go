package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/ui"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func (m TuiModel) Login() string {
	var s string
	s += ui.Title.Render("Login screen") + "\n\n"

	// Login label
	loginLabel := ui.FieldLabel.Render("Login:")
	if m.RightCursor == 0 {
		loginLabel = ui.FieldLabelActive.Render("Login:")
	}
	s += loginLabel + " " + m.login.View() + "\n"

	// Password label
	passLabel := ui.FieldLabel.Render("Password:")
	if m.RightCursor == 1 {
		passLabel = ui.FieldLabelActive.Render("Password:")
	}
	s += passLabel + " " + m.password.View() + "\n\n"

	// Button
	btn := ui.ButtonInactive.Render("Login")
	if m.RightCursor == 2 {
		btn = ui.ButtonActive.Render("Login")
	}
	s += btn + "\n"
	return s
}

func (m TuiModel) LoginAction(msg string) (tea.Model, tea.Cmd) {

	if msg == "esc" {
		if m.password.Focused() || m.login.Focused() {
			m.password.Blur()
			m.login.Blur()
		} else {
			m.HorCursor = 0
			m.RightCursor = 0
			return m, tea.ClearScreen
		}
	}
	if msg == "enter" {
		switch m.RightCursor {
		case 0:
			m.InputFocused = true
			m.login.Focus()
			m.password.Blur()
		case 1:
			m.InputFocused = true
			m.password.Focus()
			m.login.Blur()
		case 2:
			m.Loading = true
			m.InputFocused = false
			m.login.Blur()
			m.password.Blur()
			m.Loading = true
			return m, tea.Batch(
				m.Spin.Tick,
				func() tea.Msg {
					resp, err := akclient.Login(&commonv1.LoginRequest{
						Email:    m.login.Value(),
						Password: m.password.Value(),
					})
					return loginFinishedMsg{err: err, resp: resp}
				},
			)
		default:
			return m, nil
		}
	}

	return m, nil
}
