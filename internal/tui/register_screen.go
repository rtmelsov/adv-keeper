package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/ui"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func (m TuiModel) Register() string {
	var s string
	s += ui.Title.Render("Для регистрации введите данные") + "\n\n"

	// Login label
	loginLabel := ui.FieldLabel.Render("Логин:")
	if m.RightCursor == 0 {
		loginLabel = ui.FieldLabelActive.Render("Логин:")
	}
	s += loginLabel + " " + m.login.View() + "\n"

	// Password label
	passLabel := ui.FieldLabel.Render("Пароль:")
	if m.RightCursor == 1 {
		passLabel = ui.FieldLabelActive.Render("Пароль:")
	}
	s += passLabel + " " + m.password.View() + "\n\n"

	// Button
	btn := ui.ButtonInactive.Render("Регистрация")
	if m.RightCursor == 2 {
		btn = ui.ButtonActive.Render("Регистрация")
	}
	s += btn + "\n"
	return s
}

func (m TuiModel) RegisterAction(msg string) (tea.Model, tea.Cmd) {

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
			resp, err := akclient.Register(&commonv1.RegisterRequest{
				Email:    m.login.Value(),
				Password: m.password.Value(),
			})

			m.Loading = false
			if err != nil {
				m.Error = err.Error()
				return m, nil
			}
			m.Profile.Email = resp.Email
			m.token = resp.Tokens.AccessToken
			err = helpers.SaveSession(&helpers.Session{
				AccessToken: resp.Tokens.AccessToken,
				ExpiresAt:   resp.Tokens.ExpiresAt.AsTime(),
			})
			if err != nil {
				m.Error = err.Error()
				m.RightCursor = 0
				return m, nil
			}
			m.Profile.Auth = true
			m.SelectedPage = "Vault"
			m.RightCursor = 0
			m.LeftCursor = 0
		default:
			return m, nil
		}
	}

	return m, nil
}
