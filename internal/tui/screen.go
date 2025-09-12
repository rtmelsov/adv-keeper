package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"path/filepath"
	"time"

	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
)

func up(num int) int {
	if num > 0 {
		return num - 1
	}
	return num
}

func down(num int, max int) int {
	if num < max {
		return num + 1
	}
	return num
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.OpenFilePicker {
		var cmd tea.Cmd
		m.FilePicker, cmd = m.FilePicker.Update(msg)

		if ok, path := m.FilePicker.DidSelectFile(msg); ok {
			m.SelectedFile = filepath.Clean(path)
			m.OpenFilePicker = false
			return m, tea.Batch(cmd, tea.ClearScreen)
		}
		if km, ok := msg.(tea.KeyMsg); ok && (km.String() == "esc" || km.String() == "q") {
			m.OpenFilePicker = false
			return m, tea.Batch(cmd, tea.ClearScreen)
		}
		return m, cmd
	}

	if m.Loading {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.Spin, cmd = m.Spin.Update(msg)
			return m, cmd

		case tea.KeyMsg:
			// (опционально) дать пользователю отменить загрузку
			if msg.String() == "esc" {
				m.Loading = false
				m.Uploaded = 0
				m.UploadTotal = 0
				m.UploadStart = time.Time{}
				m.Uploading = false

				// download
				m.Downloaded = 0
				m.DownloadTotal = 0
				m.DownloadStart = time.Time{}
				m.Downloading = false
				m.Error = "Отменено"
				return m, nil
			}
			return m, nil

		case progressChanReadyMsg:
			switch msg.Kind {
			case OpUpload:
				m.uploadCh = msg.ch
				m.Uploading = true
				m.UploadStart = time.Now()
				return m, tea.Batch(m.Spin.Tick, listenNext(OpUpload, msg.ID, msg.ch))
			case OpDownload:
				m.downloadCh = msg.ch
				m.Downloading = true
				m.DownloadStart = time.Now()
				return m, tea.Batch(m.Spin.Tick, listenNext(OpDownload, msg.ID, msg.ch))
			}
		case progressMsg:
			switch msg.Kind {
			case OpUpload:
				m.Uploaded, m.UploadTotal = msg.Done, msg.Total
				return m, listenNext(OpUpload, msg.ID, m.uploadCh)
			case OpDownload:
				m.Downloaded, m.DownloadTotal = msg.Done, msg.Total
				return m, listenNext(OpDownload, msg.ID, m.downloadCh)
			}
		case finishedMsg:
			m.Loading = false
			switch msg.Kind {
			case OpUpload:
				m.Uploading = false
				if msg.Err != nil {
					m.Error = msg.Err.Error()
					m.uploadCh = nil
				} else {
					m.Error = ""
					m.SelectedPage = "FileList"
					m.Loading = true
					m.uploadCh = nil
					return m, tea.Batch(
						m.Spin.Tick,
						func() tea.Msg {
							list, err := akclient.GetFiles()
							return getListFinishedMsg{err: err, list: list}
						},
					)
				}
			case OpDownload:
				m.Downloading = false
				if msg.Err != nil {
					m.Error = msg.Err.Error()
				} else {
					m.Error = ""
					m.SelectedPage = "Main"
				}
				m.downloadCh = nil
			}
			return m, nil
			// продолжаем слушать

		case logoutFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
				return m, tea.ClearScreen
			}
			helpers.SaveSession(&helpers.Session{
				AccessToken: "",
				ExpiresAt:   time.Now(),
			})
			m.LeftCursor = 0
			m.Profile = &ProfileModel{}
			m.login.Reset()
			m.password.Reset()
			m.Choices = m.LoginChoices
			m.SelectedPage = "Register"
			return m, tea.ClearScreen
		case loginFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
				return m, nil
			}
			m.Profile.Email = msg.resp.Email
			err := helpers.SaveSession(&helpers.Session{
				AccessToken: msg.resp.Tokens.AccessToken,
				ExpiresAt:   msg.resp.Tokens.ExpiresAt.AsTime(),
			})
			if err != nil {
				m.Error = err.Error()
				m.RightCursor = 0
				return m, nil
			}
			m.Profile.Auth = true
			m.SelectedPage = "Main"
			m.RightCursor = 0
			m.LeftCursor = 0
			return m, nil

		case getListFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
				return m, tea.ClearScreen
			}
			m.Files = msg.list
			m.SelectedPage = "FileList"
			m.HorCursor = 1
			m.table.SetRows([]table.Row(helpers.FilesToRows(msg.list)))
			return m, nil

		case deleteFileFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
			} else {
				m.Error = ""
				return m, tea.Batch(
					m.Spin.Tick,
					func() tea.Msg {
						list, err := akclient.GetFiles()
						return getListFinishedMsg{err: err, list: list}
					},
				)
			}
			return m, nil
		case downloadFileFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
			} else {
				m.Error = ""
				m.SelectedFile = ""
				m.SelectedPage = "FileList"
			}
			return m, nil

		case uploadFinishedMsg:
			m.Loading = false
			if msg.err != nil {
				m.Error = msg.err.Error()
			} else {
				m.Error = ""
				m.SelectedPage = "FileList"
				m.SelectedFile = ""
			}
			return m, nil

		default:
			return m, nil
		}
	}

	if m.SelectedPage == "FileList" && m.HorCursor == 1 {
		// 1) Сначала свои хоткеи
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "enter":
				i := m.table.Cursor()
				if i >= 0 && i < len(m.Files.Files) {
					m.SelectedFileInfo = m.Files.Files[i]
					m.SelectedPage = "FileDetails"
					return m, nil
					// return m, func() tea.Msg { return fileChosenMsg{ID: id} }
				}
			case "esc":
				// выход из списка файлов
				m.RightCursor = 0
				m.HorCursor = 0
				return m, nil
			}
		}

		// 2) Затем отдаём событие таблице (стрелки вверх/вниз и т.п.)
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	// 2) сначала ловим завершение фоновой задачи
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.W, m.H = msg.Width, msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.InputFocused {
			if msg.String() == "esc" || msg.String() == "enter" {
				m.login.Blur()
				m.password.Blur()
				m.InputFocused = false
			}

			var cmds []tea.Cmd
			var cmd tea.Cmd
			if m.password.Focused() {
				m.password, cmd = m.password.Update(msg)
				cmds = append(cmds, cmd)
			}
			if m.login.Focused() {
				m.login, cmd = m.login.Update(msg)
				cmds = append(cmds, cmd)
			}

			return m, tea.Batch(cmds...)
		}
		if msg.String() == "up" || msg.String() == "k" {
			if m.HorCursor == 0 {
				m.LeftCursor = up(m.LeftCursor)
			} else {
				m.RightCursor = up(m.RightCursor)
			}

			return m, nil
		}
		if msg.String() == "down" || msg.String() == "j" {
			if m.HorCursor == 0 {
				length := len(m.LoginChoices) - 1
				if m.Profile.Auth {
					length = len(m.MainChoices) - 1
				}
				m.LeftCursor = down(m.LeftCursor, length)
			} else {
				m.RightCursor = down(m.RightCursor, m.MaxPageSize[m.SelectedPage])
			}
			return m, nil
		}

		if m.HorCursor == 0 {
			if msg.String() == "enter" {
				if m.Profile.Auth {

					if m.MainChoices[m.LeftCursor] == "Vault" {
						m.SelectedFile = ""
					}
					if m.MainChoices[m.LeftCursor] == "Logout" {
						m.Loading = true
						return m, tea.Batch(
							m.Spin.Tick,
							func() tea.Msg {
								_, err := akclient.Logout()
								return logoutFinishedMsg{err: err}
							},
						)

					}
					if m.MainChoices[m.LeftCursor] == "FileList" {
						m.Loading = true
						return m, tea.Batch(
							m.Spin.Tick,
							func() tea.Msg {
								list, err := akclient.GetFiles()
								return getListFinishedMsg{err: err, list: list}
							},
						)
					}
					m.SelectedPage = m.MainChoices[m.LeftCursor]
				} else {
					m.SelectedPage = m.LoginChoices[m.LeftCursor]
				}
				m.password.Reset()
				m.login.Reset()
				m.HorCursor = 1
				return m, tea.ClearScreen
			}
			if msg.String() == "esc" {
				m.login.Blur()
				m.password.Blur()
				m.InputFocused = false
				m.HorCursor = 0
				m.RightCursor = 0
				return m, nil
			}
		}

		if m.SelectedPage == "FileDetails" {
			return m.FileDetailsAction(msg.String())
		}
		if m.SelectedPage == "Vault" {
			return m.VaultAction(msg.String())
		}
		if m.SelectedPage == "Login" {
			return m.LoginAction(msg.String())
		}
		if m.SelectedPage == "Logout" {
			return m.LoginAction(msg.String())
		}
		if m.SelectedPage == "Register" {
			return m.RegisterAction(msg.String())
		}
		if m.SelectedPage == "Main" {
			return m.MainAction(msg.String())
		}
	}
	return m, nil
}
