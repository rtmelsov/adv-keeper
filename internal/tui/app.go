package tui

import (
	filepicker "github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
)

type ProfileModel struct {
	UserID   string
	Email    string
	DeviceID string
	Auth     bool
}

type TuiModel struct {
	SelectedFile   string
	FilePicker     filepicker.Model
	OpenFilePicker bool
	table          table.Model
	Loading        bool
	History        []string
	Error          string
	Profile        *ProfileModel
	Choices        []string // элементы списка
	InputFocused   bool
	Cursor         int // индекс строки под курсором
	CursorHor      int // индекс строки под курсором
	login          textinput.Model
	password       textinput.Model
	Selected       string
}

func (m TuiModel) PickerInit() tea.Cmd {
	m.FilePicker = filepicker.New()
	// пример: показывать скрытые файлы и разрешить выбирать только файлы
	m.FilePicker.ShowHidden = false
	m.FilePicker.DirAllowed = false
	m.FilePicker.FileAllowed = true
	// если хотите ограничить по расширениям:
	// m.FilePicker.AllowedTypes = []string{".txt", ".json"}

	return m.FilePicker.Init()
}

func InitialModel() TuiModel {
	login := textinput.New()
	login.Placeholder = "your login"
	login.CharLimit = 64
	login.Prompt = ""
	login.Focus()

	password := textinput.New()
	password.Placeholder = "your password"
	password.CharLimit = 64
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.Prompt = ""

	fp := filepicker.New()
	fp.SetHeight(14)
	fp.ShowHidden = false
	fp.DirAllowed = false
	fp.FileAllowed = true

	return TuiModel{
		Choices:    []string{"Vault", "Menu", "Register", "Login"},
		History:    []string{},
		table:      InitTable(),
		FilePicker: fp,
		Selected:   "Menu",
		Profile:    &ProfileModel{Auth: false},
		login:      login,
		password:   password,
	}
}

func (m TuiModel) Init() tea.Cmd { return tea.ClearScreen }
