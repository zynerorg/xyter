package token

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Msg emitted when token is saved
type TokenSavedMsg struct{}

type Model struct {
	TokenInput textinput.Model
	PassInput  textinput.Model
	Encrypt    bool
	Focus      int
}

func New() Model {
	token := textinput.New()
	token.Placeholder = "Paste API token"
	token.Focus()
	token.Width = 40

	pass := textinput.New()
	pass.Placeholder = "Encryption password"
	pass.EchoMode = textinput.EchoPassword
	pass.Width = 40

	return Model{
		TokenInput: token,
		PassInput:  pass,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.Focus = (m.Focus + 1) % 2
			if m.Focus == 0 {
				m.TokenInput.Focus()
				m.PassInput.Blur()
			} else {
				m.TokenInput.Blur()
				m.PassInput.Focus()
			}
		case "e":
			m.Encrypt = !m.Encrypt
		case "enter":
			// Here you would save the token
			return m, func() tea.Msg { return TokenSavedMsg{} }
		}
	}

	m.TokenInput, cmd = m.TokenInput.Update(msg)
	if m.Encrypt {
		m.PassInput, _ = m.PassInput.Update(msg)
	}

	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	boxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("63"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	checkbox := "[ ] Encrypt token"
	if m.Encrypt {
		checkbox = "[x] Encrypt token"
	}

	s := labelStyle.Render("Token") + "\n" + m.TokenInput.View() + "\n\n"
	s += labelStyle.Render("Encrypt token") + "\n" + checkbox + "\n\n"
	if m.Encrypt {
		s += labelStyle.Render("Password") + "\n" + m.PassInput.View() + "\n\n"
	}
	s += "Press Enter to save • e toggle encryption • tab switch input"

	return boxStyle.Render(s)
}
