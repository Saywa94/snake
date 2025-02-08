package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "No intialized"
	}

	// TODO: Render different view according to game state

	if m.gameState == "start" {
		// Render start screen
		var startStyle = lipgloss.NewStyle().
			SetString("Press [Enter] to Start Game").
			Width(40).
			Height(5).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Align(lipgloss.Center, lipgloss.Center)

		var dialog = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, startStyle.Render())

		return dialog
	}

	title := fmt.Sprintf("Score: %d", m.score)
	title = getCenteredTitle(title, m.width)

	canvass := ""

	for row := range m.grid {
		for col := range m.grid[row] {
			canvass += string(m.grid[row][col])
		}
		canvass += "\n"
	}

	s := style.Render(title)
	s += "\n"
	s += style2.Render(fmt.Sprintf("Position: (%d, %d)", m.snake.Head.X, m.snake.Head.Y) + " " + fmt.Sprintf("Paused: %t", Paused))
	s += "\n"
	s += canvass

	return s

}
