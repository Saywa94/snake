package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	score uint
	position
	width  int
	height int
}

type position struct {
	x         int
	y         int
	axis      string
	direction int
}

func NewModel() model {
	return model{
		score:  0,
		width:  30,
		height: 30,
		position: position{
			x:         0,
			y:         0,
			axis:      "x",
			direction: 1,
		},
	}
}

func main() {

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		if m.position.axis == "x" {
			m.position.x += m.position.direction
		}
		if m.position.axis == "y" {
			m.position.y += m.position.direction
		}
		return m, doTick()

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.position.axis = "y"
			m.position.direction = -1
		case "down":
			m.position.axis = "y"
			m.position.direction = 1
		case "right":
			m.position.axis = "x"
			m.position.direction = 1
		case "left":
			m.position.axis = "x"
			m.position.direction = -1
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.position.x = m.width / 2
		m.position.y = m.height/2 - 1
	}

	return m, nil
}

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "No intialized"
	}

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

			// Borders
			if (row == 0 && col == 0) || (row == 0 && col == m.width-1) || (row == len(grid)-1 && col == 0) || (row == len(grid)-1 && col == m.width-1) {
				grid[row][col] = '+'
			} else if row == 0 || row == len(grid)-1 {
				grid[row][col] = '='
			} else if col == 0 || col == len(grid[row])-1 {
				grid[row][col] = '|'
			}

			if col == m.position.x && row == m.position.y {
				grid[row][col] = 'O'
			}

			canvass += string(grid[row][col])
		}
		canvass += "\n"
	}

	s := title
	s += "\n"
	s += fmt.Sprintf("Canvass size: (%d, %d)", m.width, m.height) + " " + fmt.Sprintf("Position: (%d, %d)", m.position.x, m.position.y)
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
