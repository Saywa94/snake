package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var startStyle = lipgloss.NewStyle().
	Width(40).
	Height(5).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")).
	Align(lipgloss.Center, lipgloss.Center)

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "No intialized"
	}

	// TODO: Render different view according to game state

	if m.gameState == "start" {
		// Render start screen
		start := startStyle.SetString("Press [Enter] to Start Game, \n or [q] to Quit")

		var dialog = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, start.Render())

		return dialog
	}

	if m.gameState == End {

		endString := startStyle.SetString("Game Over!!")
		s := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, endString.Render())

		return s
	}

	title := getCenteredTitle(fmt.Sprintf("Score: %d", m.score), m.width)

	canvass := ""

	for row := range m.grid {
		for col := range m.grid[row] {
			canvass += string(m.grid[row][col])
		}
		canvass += "\n"
	}

	s := style.Render(title)
	s += "\n"
	s += style2.Render(fmt.Sprintf("Position: (%d, %d)", m.snake.Head.X, m.snake.Head.Y) + " " + fmt.Sprintf("Game State: %s", m.gameState))
	s += "\n"
	s += canvass

	return s

}
