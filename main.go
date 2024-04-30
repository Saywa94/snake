package main

import (
	"fmt"
	"os"
	"strings"

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
	canvass   [][]rune
}

func NewModel() model {
	return model{
		score:     0,
		positionY: 1,
		positionX: 1,
		width:     30,
		height:    30,
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

	title := fmt.Sprintf("Score: %d", m.score)
	title = getCenteredTitle(title, m.width)

	grid := make([][]rune, m.height-3)

	for i := range grid {
		grid[i] = make([]rune, m.width)
	}
	canvass := ""

	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = ' '
			if row == 0 || row == len(grid)-1 {
				grid[row][col] = '='
			}
			if col == 0 || col == len(grid[row])-1 {
				grid[row][col] = '|'
			}
			canvass += string(grid[row][col])
		}
		canvass += "\n"
	}

	s := title
	s += "\n"
	s += fmt.Sprintf("Canvass size: (%d, %d)", m.width, m.height) + fmt.Sprintf("Position: (%d, %d)", m.positionY, m.positionX)
	s += "\n"
	s += canvass

	return s

}

func getPaddingLeft(title string, width int) int {
	spaces := width/2 - len(title)/2
	if spaces < 0 {
		spaces = 0
	}
	return spaces
}
func getCenteredTitle(title string, width int) string {
	return strings.Repeat(" ", getPaddingLeft(title, width)) + title
}
