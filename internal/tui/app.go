package tui

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/akclient"
	"github.com/rtmelsov/adv-keeper/internal/models"
	"time"
)

type ProfileModel struct {
	Email string
	Auth  bool
}

type TuiModel struct {
	uploadCh    <-chan models.Prog
	Uploaded    int64
	UploadTotal int64
	UploadStart time.Time
	Uploading   bool

	// download
	downloadCh    <-chan models.Prog
	Downloaded    int64
	DownloadTotal int64
	DownloadStart time.Time
	Downloading   bool

	LoaderCount      models.LoaderType
	W                int
	H                int
	SideBarW         int
	Files            *filev1.GetFilesResponse
	SelectedFile     string
	token            string
	Busy             int           // >0 — идёт фоновая операция
	Spin             spinner.Model // bubbles/spinner
	Err              error
	Notice           string
	SelectedFileInfo *filev1.FileItem
	FilePicker       filepicker.Model
	OpenFilePicker   bool
	table            table.Model
	Loading          bool
	Error            string
	Profile          *ProfileModel
	Choices          []string // элементы списка
	LoginChoices     []string
	MainChoices      []string
	InputFocused     bool
	LeftCursor       int // индекс строки под курсором
	RightCursor      int // индекс строки под курсором
	HorCursor        int // индекс строки под курсором
	login            textinput.Model
	password         textinput.Model
	SelectedPage     string
	MaxPageSize      map[string]int
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

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	selected := "Main"
	profile := &ProfileModel{Auth: false}
	mainChoices := []string{"Vault", "FileList", "Main", "Logout"}
	choices := mainChoices
	loginChoices := []string{"Register", "Login"}

	resp, err := akclient.GetProfile()
	if err != nil {
		selected = "Register"
		choices = []string{"Register", "Login"}
	} else {
		profile.Email = resp.Email
		profile.Auth = true
	}

	MaxPageSize := map[string]int{
		"Vault": 0, "Main": 1, "Register": 2, "Login": 2, "FileList": 0, "FileDetails": 1,
	}

	return TuiModel{
		Files:        &filev1.GetFilesResponse{},
		Choices:      choices,
		MainChoices:  mainChoices,
		LoginChoices: loginChoices,
		W:            88,
		H:            44,
		HorCursor:    1,
		SideBarW:     24,
		table:        InitTable(),
		FilePicker:   fp,
		Spin:         sp,
		SelectedPage: selected,
		MaxPageSize:  MaxPageSize,
		Profile:      profile,
		login:        login,
		password:     password,
	}
}

func (m TuiModel) Init() tea.Cmd { return tea.ClearScreen }
