package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/decentralized-chat/internal/ui"
)

func main() {
	p := tea.NewProgram(
		ui.NewMainModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
