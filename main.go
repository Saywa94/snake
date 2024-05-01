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
	score  uint
	head   position
	body   []position
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
		head: position{
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
	hasColided := false
	if m.head.x == 0 || m.head.x == len(m.grid[0])-1 || m.head.y == 0 || m.head.y == len(m.grid)-1 {
		hasColided = true
	}

	for _, p := range m.body {
		if p.x == m.head.x && p.y == m.head.y {
			hasColided = true
		}
	}
	return hasColided
}

func (m *model) Advance() {
	// Empty previous position
	prevPos := m.head
	m.grid[prevPos.y][prevPos.x] = ' '

	// Restrict backwards movement
	if m.head.axis == "x" {
		m.head.x += m.head.direction
	}
	if m.head.axis == "y" {
		m.head.y += m.head.direction
	}

	// Fill new position
	m.grid[m.head.y][m.head.x] = '@'

	// Move body
	if len(m.body) > 0 {
		last := m.body[len(m.body)-1]
		m.grid[last.y][last.x] = ' '
	}
	var newBody []position
	for i := range m.body {
		if i == 0 {
			newBody = append(newBody, prevPos)
		} else {
			newBody = append(newBody, m.body[i-1])
		}
	}

	m.body = newBody
	for _, p := range m.body {
		m.grid[p.y][p.x] = 'o'
	}
}

func (m *model) PlaceCrumb() {
	crumbX := 0
	crumbY := 0
	for {
		crumbX = rand.IntN(len(m.grid[0])-2) + 1
		crumbY = rand.IntN(len(m.grid)-2) + 1

		if m.grid[crumbY][crumbX] == ' ' {
			break
		}
	}

	m.grid[crumbY][crumbX] = 'x'
	m.crumb = position{crumbX, crumbY, "", 0}
}

func (m *model) AddBodyPart() {
	p := position{}
	m.body = append(m.body, p)
}

func (m *model) FillGrid() {
	for row := range m.grid {
		for col := range m.grid[row] {
			m.grid[row][col] = ' '

			// Add walls
			if (row == 0 && col == 0) || (row == 0 && col == m.width-1) || (row == len(m.grid)-1 && col == 0) || (row == len(m.grid)-1 && col == m.width-1) {
				m.grid[row][col] = '+'
			} else if row == 0 || row == len(m.grid)-1 {
				m.grid[row][col] = '='
			} else if col == 0 || col == len(m.grid[row])-1 {
				m.grid[row][col] = '|'
			}

			if col == m.head.x && row == m.head.y {
				m.grid[row][col] = '@'
			}
			if col == m.crumb.x && row == m.crumb.y {
				m.grid[row][col] = 'x'
			}

		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.head.axis == "y" && m.head.direction == 1 {
				return m, nil
			}
			m.head.axis = "y"
			m.head.direction = -1
		case "down":
			if m.head.axis == "y" && m.head.direction == -1 {
				return m, nil
			}
			m.head.axis = "y"
			m.head.direction = 1
		case "right":
			if m.head.axis == "x" && m.head.direction == -1 {
				return m, nil
			}
			m.head.axis = "x"
			m.head.direction = 1
		case "left":
			if m.head.axis == "x" && m.head.direction == 1 {
				return m, nil
			}
			m.head.axis = "x"
			m.head.direction = -1
		case " ":
			m.PlaceCrumb()
		}

	case TickMsg:

		m.Advance()

		if m.CheckCollision() {
			return m, tea.Quit
		}

		// Check if the crumb has been eaten
		if m.head.x == m.crumb.x && m.head.y == m.crumb.y {
			m.score++
			m.PlaceCrumb()

			// Add new body part
			m.AddBodyPart()

		}

		// Check if we need to speed up
		if m.head.axis == "x" {
			return m, doTick(normalSpeed)
		} else {
			return m, doTick(slowSpeed)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Position the player in the center
		m.head.x = m.width / 2
		m.head.y = m.height/2 - 1

		// Create temporary body
		m.body = []position{
			{
				x:         m.head.x - 1,
				y:         m.head.y,
				axis:      m.head.axis,
				direction: m.head.direction,
			},
		}

		// Create playlale grid
		m.grid = make([][]rune, m.height-3)
		for i := range m.grid {
			m.grid[i] = make([]rune, m.width)
		}

		// Fill in the grid
		m.FillGrid()

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

	canvass := ""

	for row := range m.grid {
		for col := range m.grid[row] {
			canvass += string(m.grid[row][col])
		}
		canvass += "\n"
	}

	s := title
	s += "\n"
	s += fmt.Sprintf("Canvass size: (%d, %d)", m.width, m.height) + " " + fmt.Sprintf("Position: (%d, %d)", m.head.x, m.head.y) + " " + fmt.Sprintf("Parts: %d", len(m.body))
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
