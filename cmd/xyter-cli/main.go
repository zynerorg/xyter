package main

import (
	"fmt"
	"os"

	"git.zyner.org/meta/xyter/internal/tui/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp <command>")
		fmt.Println("Commands: tui, save-token, show-token")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "tui":
		runTUI()
	case "save-token":
		// CLI logic here (later)
		fmt.Println("CLI save-token not implemented yet")
	case "show-token":
		fmt.Println("CLI show-token not implemented yet")
	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}

func runTUI() {
	p := tea.NewProgram(app.NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
