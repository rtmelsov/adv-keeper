package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/models"
	"github.com/rtmelsov/adv-keeper/internal/ui"

	"os"
)

type OpKind string

const (
	OpUpload   OpKind = "upload"
	OpDownload OpKind = "download"
	// OpIndex, OpSync и т.д.
)

type logoutFinishedMsg struct{ err error }
type loginFinishedMsg struct {
	err  error
	resp *commonv1.LoginResponse
}
type getListFinishedMsg struct {
	err  error
	list *filev1.GetFilesResponse
}
type deleteFileFinishedMsg struct{ err error }
type downloadFileFinishedMsg struct{ err error }

// сообщения
type progressChanReadyMsg struct {
	Kind OpKind
	ID   string
	ch   <-chan models.Prog
}
type progressMsg struct {
	Kind  OpKind
	ID    string
	Done  int64
	Total int64
}

type finishedMsg struct {
	Kind OpKind
	ID   string
	Err  error
}
type uploadFinishedMsg struct{ err error }

func listenNext(kind OpKind, id string, ch <-chan models.Prog) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return finishedMsg{Kind: kind, ID: id, Err: nil}
		}
		if p.Err != nil {
			return finishedMsg{Kind: kind, ID: id, Err: p.Err}
		}
		return progressMsg{Kind: kind, ID: id, Done: p.Done, Total: p.Total}
	}
}

func startUploadCmd(id, path string) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan models.Prog, 32)
		go akclient.UploadFile(path, ch)
		return progressChanReadyMsg{Kind: OpUpload, ID: id, ch: ch}
	}
}

func startDownloadCmd(id, fileID, outPath string) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan models.Prog, 32)
		go akclient.DownloadFile(fileID, ch)
		return progressChanReadyMsg{Kind: OpDownload, ID: id, ch: ch}
	}
}

func (m TuiModel) VaultAction(msg string) (tea.Model, tea.Cmd) {
	switch msg {
	case "esc":
		m.HorCursor = 0
		m.RightCursor = 0
		return m, tea.ClearScreen

	case "enter":
		if m.SelectedFile == "" {
			if home, err := os.UserHomeDir(); err == nil && home != "" {
				m.FilePicker.CurrentDirectory = home
				m.FilePicker.Path = home
			} else if wd, _ := os.Getwd(); wd != "" {
				m.FilePicker.CurrentDirectory = wd
				m.FilePicker.Path = wd
			} else {
				m.FilePicker.CurrentDirectory = "/"
				m.FilePicker.Path = "/"
			}

			m.FilePicker.SetHeight(14) // чтобы было видно список
			m.OpenFilePicker = true
			return m, m.FilePicker.Init() // ← важный момент
		} else if m.SelectedFile != "" {
			m.Loading = true
			return m, tea.Batch(
				m.Spin.Tick,
				func() tea.Msg {
					ch := make(chan models.Prog, 32)
					go akclient.UploadFile(m.SelectedFile, ch)
					return progressChanReadyMsg{ch: ch, Kind: OpUpload}
				},
			)
		}
	}
	return m, nil
}

func (m TuiModel) Vault() string {

	s := ui.Title.Render("try to add some file\n\n")
	if m.OpenFilePicker {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			m.FilePicker.View(),
		)
	} else if m.SelectedFile != "" {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			ui.Title.Render(fmt.Sprintf("\nВы выбрали файл: %s\n", m.SelectedFile)),
		)
	} else {
		btn := ui.ButtonInactive.Render("ADD FILE")
		if m.HorCursor == 1 {
			btn = ui.ButtonActive.Render("ADD FILE")
		}
		return lipgloss.JoinVertical(
			lipgloss.Top,
			s,
			btn,
		)
	}
}
