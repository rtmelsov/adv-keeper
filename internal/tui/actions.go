package tui

import (
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
)

func (m *TuiModel) MenuAction() {
	m.Selected = m.Choices[m.Cursor]
	m.Cursor = 0
}

func (m *TuiModel) VaultActions() {

}

func (m *TuiModel) RegisterAction() {
	switch m.Cursor {
	case 0:
		m.InputFocused = true
		m.login.Focus()
		m.password.Blur()
	case 1:
		m.InputFocused = true
		m.password.Focus()
		m.login.Blur()
	default:
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
			m.Selected = "Register"
			m.Cursor = 0
			return
		}

		m.Profile.UserID = resp.UserId
		m.Profile.DeviceID = resp.DeviceId
		m.Profile.Email = resp.Email
		err = helpers.SaveSession(&helpers.Session{
			AccessToken: resp.Tokens.AccessToken,
			ExpiresAt:   resp.Tokens.ExpiresAt.AsTime(),
		})
		if err != nil {
			m.Error = err.Error()
			m.Selected = "Register"
			m.Cursor = 0
			return
		}
		m.Profile.Auth = true
		m.Selected = "Vault"
		m.Cursor = 1
	}
}
