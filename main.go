package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

type model struct {
	score     uint
	positionY int
	positionX int
	width     int
	height    int
}

func NewModel() model {
	return model{
		score:     0,
		positionY: 1,
		positionX: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.score++
		case "down":
			m.score--
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Score: %d", m.score)

	s += "\n"
	s += fmt.Sprintf("Canvass size: (%d, %d)", m.width, m.height) + fmt.Sprintf("Position: (%d, %d)", m.positionY, m.positionX)

	return s

}
