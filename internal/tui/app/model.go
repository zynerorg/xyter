package app

import (
	"git.zyner.org/meta/xyter/internal/tui/components/token"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	MainScreen Screen = iota
	OtherScreen
)

type Model struct {
	Screen    Screen
	Main      tea.Model
	Token     tea.Model
	ShowToken bool
	Width     int
	Height    int
}

func NewModel() Model {
	// width, height, err := term.GetSize(os.Stdout.Fd())
	// if err != nil {
	// 	width = 80
	// 	height = 24
	// }
	return Model{
		Screen:    MainScreen,
		Main:      NewMainScreen(),
		Token:     token.New(),
		ShowToken: true, // true for first-time launch
		Width:     80,
		Height:    24,
	}
}

func (m Model) Init() tea.Cmd {
	if m.ShowToken {
		return m.Token.Init()
	}
	return m.Main.Init()
}
