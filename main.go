package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	score  uint
	head   position
	body   []position
	grid   [][]string
	width  int
	height int
	crumb
}

type crumb struct {
	x     int
	y     int
	style lipgloss.Style
}

type position struct {
	x         int
	y         int
	axis      string
	direction int
	content   string
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
			content:   "o",
		},
		body: []position{
			{x: 2, y: 2, content: "o"},
		},
	}
}

// TODO:Figure out how to make this configurable
const (
	normalSpeed = 30
	slowSpeed   = 55
)

var Paused = true

func main() {

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

type TickMsg time.Time

func doTick(ms time.Duration) tea.Cmd {
	Paused = false
	return tea.Tick(time.Millisecond*ms, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	// TODO: Have a start screen
	// before starting the game

	// return doTick(normalSpeed)
	return nil
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
	m.grid[prevPos.y][prevPos.x] = " "

	// Restrict backwards movement
	if m.head.axis == "x" {
		m.head.x += m.head.direction
	}
	if m.head.axis == "y" {
		m.head.y += m.head.direction
	}

	// Fill new position
	m.grid[m.head.y][m.head.x] = "@"

	// Move body
	if len(m.body) > 0 {
		last := m.body[len(m.body)-1]
		m.grid[last.y][last.x] = " "
	}
	var newBody []position
	for i, p := range m.body {
		if i == 0 {
			newBody = append(newBody, prevPos)
		} else {
			newBody = append(newBody, position{
				x:       m.body[i-1].x,
				y:       m.body[i-1].y,
				content: p.content,
			})
		}
	}

	m.body = newBody
	for _, p := range m.body {
		m.grid[p.y][p.x] = p.content
	}
}

func (m *model) PlaceCrumb() {
	crumbX := 0
	crumbY := 0
	for {
		crumbX = rand.IntN(len(m.grid[0])-2) + 1
		crumbY = rand.IntN(len(m.grid)-2) + 1

		if m.grid[crumbY][crumbX] == " " {
			break
		}
	}

	color := NextColor(m.score)
	var crumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	m.grid[crumbY][crumbX] = crumbStyle.Render("x")
	m.crumb = crumb{crumbX, crumbY, crumbStyle}
}

func (m *model) AddBodyPart() {
	style := m.crumb.style
	p := position{
		content: style.Render("o"),
	}
	m.body = append(m.body, p)
}

func (m *model) FillGrid() {
	for row := range m.grid {
		for col := range m.grid[row] {
			m.grid[row][col] = " "

			// Add walls
			if (row == 0 && col == 0) || (row == 0 && col == m.width-1) || (row == len(m.grid)-1 && col == 0) || (row == len(m.grid)-1 && col == m.width-1) {
				m.grid[row][col] = "+"
			} else if row == 0 || row == len(m.grid)-1 {
				m.grid[row][col] = "="
			} else if col == 0 || col == len(m.grid[row])-1 {
				m.grid[row][col] = "|"
			}

			if col == m.head.x && row == m.head.y {
				m.grid[row][col] = "@"
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
			// Implement pause/resume
			if Paused == false {
				Paused = true
				return m, nil
			} else {
				Paused = true
				return m, doTick(normalSpeed)
			}
		case "r":
			if Paused == true {
				Paused = false
				m.RestartGame()
				return m, doTick(normalSpeed)
			}

		}

	case TickMsg:

		if Paused == true {
			return m, nil
		}

		m.Advance()

		if m.CheckCollision() {
			// Game over
			// TODO: Show score + option to start new game

			// return m, tea.Quit
			Paused = true
			return m, nil
		}

		// Check if the crumb has been eaten
		if m.head.x == m.crumb.x && m.head.y == m.crumb.y {
			m.score++

			// Add new body part
			m.AddBodyPart()

			// Place new crumb at random position
			m.PlaceCrumb()

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
		m.grid = make([][]string, m.height-3)
		for i := range m.grid {
			m.grid[i] = make([]string, m.width)
		}

		// Fill in the grid
		m.FillGrid()

		// Add one crumb
		m.PlaceCrumb()
	}

	return m, nil
}

var style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7bdff2"))
var style2 = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "No intialized"
	}

	// TODO: Render different view according to game state

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
	s += style2.Render(fmt.Sprintf("Position: (%d, %d)", m.head.x, m.head.y) + " " + fmt.Sprintf("Paused: %t", Paused))
	s += "\n"
	s += canvass

	return s

}

func (m *model) RestartGame() {
	m.score = 0
	m.head = position{
		x:         m.width / 2,
		y:         m.height/2 - 1,
		axis:      "x",
		direction: 1,
		content:   "o",
	}
	m.body = []position{
		{x: 2, y: 2, content: "o"},
	}

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
