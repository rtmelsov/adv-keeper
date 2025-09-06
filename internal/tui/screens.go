package tui

import (
	"fmt"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
)

func (m TuiModel) Menu() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}

func (m TuiModel) Register() string {
	s := "Login screen\n"

	if m.Cursor == 0 {
		s += "> Login:"
	} else {
		s += "  Login:"
	}
	s += m.login.View() + "\n"
	if m.Cursor == 1 {
		s += "> Password:"
	} else {
		s += "  Password:"
	}
	s += m.password.View() + "\n"

	s += "\n"
	btn := "[ Login ]"
	if m.Cursor == 2 {
		btn = "> " + btn
	} else {
		btn = "  " + btn
	}

	s += "\n" + btn + "\n"
	return s
}

func (m TuiModel) Login() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}

func (m TuiModel) Vault() string {
	s := fmt.Sprintf("Hello, %s!\n\n", m.Profile.Email)

	if m.CursorHor == 0 {
		if m.OpenFilePicker {
			s += m.FilePicker.View()
		} else if m.SelectedFile != "" {
			resp := make(chan string)

			go func(resp chan<- string) {
				var r string
				_, err := akclient.UploadFile(m.SelectedFile)
				if err != nil {
					r = fmt.Sprintf("----ERROR: %s", err.Error())
				} else {
					r = "----SUCCESS----"
				}

				resp <- r

			}(resp)

			s += <-resp

		}
		if m.Cursor == 0 {
			s += "[ADD FILE]"
		}
		if m.Cursor == 1 {
			s += "<[ADD FILE]>"
		}
	} else {
		s += baseStyle.Render(m.table.View()) + "\n"
	}
	return s
}
