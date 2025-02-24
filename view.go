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

		var scoreStyle = lipgloss.NewStyle().
			MarginTop(10).
			Padding(1, 3).
			Align(lipgloss.Center).
			Bold(true).
			Background(m.crumb.prevStyle.GetForeground()).
			Foreground(lipgloss.Color("#353535"))

		var scoreListStyle = lipgloss.NewStyle().
			Width(50).
			Height(15).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("63")).
			Align(lipgloss.Center, lipgloss.Top).
			MarginTop(2).
			Padding(1)

		// Show current score in big letters
		score := scoreStyle.Render(fmt.Sprintf("SCORE: %d", m.score))
		s := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, score)

		// TODO: Show top 10 scores

		// TODO: If current score in top 10 show text input in correct place

		endString := scoreListStyle.SetString("Game Over!!")
		s += lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, endString.Render())

		return s
	}

	title := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, fmt.Sprintf("Score: %d", m.score))

	canvass := ""

	for row := range m.grid {
		for col := range m.grid[row] {
			canvass += string(m.grid[row][col])
		}
		canvass += "\n"
	}

	infoString := style2.Render(fmt.Sprintf("Position: (%d, %d)", m.snake.Head.X, m.snake.Head.Y)) + " "
	infoString += style2.Render(fmt.Sprintf("Game State: %s", m.gameState)) + " "
	infoString += fmt.Sprintf("First cell %s", m.grid[0][0]) + " "
	infoString += fmt.Sprintf("Body length %d", len(m.snake.Body))

	s := m.crumb.prevStyle.Render(title)
	s += "\n"
	s += infoString
	s += "\n"
	s += canvass

	return s

}
