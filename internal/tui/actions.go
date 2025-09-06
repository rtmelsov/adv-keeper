package tui

import (
	"fmt"
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
)

func (m *TuiModel) MenuAction() {
	m.Selected = m.Choices[m.Cursor]
	m.Cursor = 0
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

		fmt.Println("\nresp", resp)

		m.Profile.UserID = resp.UserId
		m.Profile.DeviceID = resp.DeviceId
		m.Profile.Email = resp.Email
		m.Profile.Auth = true
		m.Selected = "Vault"
		m.Cursor = 0
	}
}
