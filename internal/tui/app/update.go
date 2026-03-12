package app

import (
	"git.zyner.org/meta/xyter/internal/tui/components/token"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// global keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// terminal resize
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width = sizeMsg.Width
		m.Height = sizeMsg.Height
	}
	var cmd tea.Cmd

	if m.ShowToken {
		m.Token, cmd = m.Token.Update(msg)
		if _, ok := msg.(token.TokenSavedMsg); ok {
			m.ShowToken = false
		}
	} else {
		switch m.Screen {
		case MainScreen:
			m.Main, cmd = m.Main.Update(msg)
		}
	}

	return m, cmd
}
