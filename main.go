package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

// TODO: Figure out how to make this configurable
const (
	normalSpeed   = 30
	slowSpeed     = 55
	extraRowsUsed = 3
)

func main() {

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
