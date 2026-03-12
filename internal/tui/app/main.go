package app

import tea "github.com/charmbracelet/bubbletea"

type mainScreen struct{}

func NewMainScreen() tea.Model {
	return mainScreen{}
}

func (m mainScreen) Init() tea.Cmd {
	return nil
}

func (m mainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mainScreen) View() string {
	return "XYTER MAIN SCREEN\n\nPress q to quit"
}
