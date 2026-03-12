package app

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	view := ""
	if m.ShowToken {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			m.Token.View(),
		)

	} else {
		view = m.Main.View()
	}

	return view
}
