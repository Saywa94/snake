package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	score uint
	position
	grid   [][]rune
	width  int
	height int
	crumb  position
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

const (
	normalSpeed = 30
	slowSpeed   = 55
)

func main() {

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

type TickMsg time.Time

func doTick(ms time.Duration) tea.Cmd {
	return tea.Tick(time.Millisecond*ms, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick(normalSpeed)
}

func (m model) CheckCollision() bool {
	// Border collision
	if m.position.x == 0 || m.position.x == len(m.grid[0])-1 || m.position.y == 0 || m.position.y == len(m.grid)-1 {
		return true
	}
	return false
}

func (m *model) Advance() {
	if m.position.axis == "x" {
		m.position.x += m.position.direction
	}
	if m.position.axis == "y" {
		m.position.y += m.position.direction
	}
}

func (m *model) PlaceCrumb() {
	crumbX := rand.IntN(len(m.grid[0])) + 1
	crumbY := rand.IntN(len(m.grid)-2) + 1
	m.crumb = position{crumbX, crumbY, "", 0}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.position.axis == "y" && m.position.direction == 1 {
				return m, nil
			}
			m.position.axis = "y"
			m.position.direction = -1
		case "down":
			if m.position.axis == "y" && m.position.direction == -1 {
				return m, nil
			}
			m.position.axis = "y"
			m.position.direction = 1
		case "right":
			if m.position.axis == "x" && m.position.direction == -1 {
				return m, nil
			}
			m.position.axis = "x"
			m.position.direction = 1
		case "left":
			if m.position.axis == "x" && m.position.direction == 1 {
				return m, nil
			}
			m.position.axis = "x"
			m.position.direction = -1
			// case " ":
			// 	m.PlaceCrumb()
		}

	case TickMsg:

		m.Advance()

		if m.CheckCollision() {
			return m, tea.Quit
		}

		if m.position.x == m.crumb.x && m.position.y == m.crumb.y {
			m.score++
			m.PlaceCrumb()
		}

		if m.position.axis == "x" {
			return m, doTick(normalSpeed)
		} else {
			return m, doTick(slowSpeed)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Position the player in the center
		m.position.x = m.width / 2
		m.position.y = m.height/2 - 1

		// Create playlale grid
		m.grid = make([][]rune, m.height-3)
		for i := range m.grid {
			m.grid[i] = make([]rune, m.width)
		}

		// Add one crumb
		m.PlaceCrumb()
	}

	return m, nil
}

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "No intialized"
	}

	title := fmt.Sprintf("Score: %d", m.score)
	title = getCenteredTitle(title, m.width)

	grid := m.grid

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
				grid[row][col] = '@'
			}
			if col == m.crumb.x && row == m.crumb.y {
				grid[row][col] = 'x'
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
