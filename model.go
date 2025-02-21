package main

import (
	"github.com/Saywa94/snake/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand/v2"
	"strings"
	"time"
)

type crumb struct {
	x         int
	y         int
	style     lipgloss.Style
	prevStyle lipgloss.Style
}

type GameState string

const (
	Start   GameState = "start"
	Playing GameState = "playing"
	Paused  GameState = "paused"
	End     GameState = "end"
)

type model struct {
	gameState GameState
	store     *Store
	score     uint
	grid      [][]string
	width     int
	height    int
	crumb
	snake game.Snake
}

type TickMsg time.Time

func NewModel(store *Store) model {
	return model{
		gameState: Start,
		store:     store,
		score:     0,
		width:     30,
		height:    30,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func doTick(ms time.Duration) tea.Cmd {
	return tea.Tick(time.Millisecond*ms, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		key := msg.String()
		switch m.gameState {
		case End:
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			case "r":
				m.gameState = Playing
				m.RestartGame()
				return m, doTick(normalSpeed)
			}
		case Start:
			switch key {
			case "enter":

				m.gameState = Playing
				m.snake = game.CreateSnake(m.width, m.height)

				// Create playlale grid
				m.grid = make([][]string, m.height)
				for i := range m.grid {
					m.grid[i] = make([]string, m.width)
				}
				// Fill in the grid
				m.FillGrid()
				// Add one crumb
				m.PlaceCrumb()
				return m, doTick(normalSpeed)
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				if m.gameState != End {
					return m, tea.Quit
				}
			}

		case Paused:
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			case " ":
				// resume game
				m.gameState = Playing
				return m, doTick(normalSpeed)
			case "r":
				m.gameState = Playing
				m.RestartGame()
				return m, doTick(normalSpeed)
			}

		case Playing:
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			case "up":
				if m.snake.Head.Axis == "y" && m.snake.Head.Direction == 1 {
					return m, nil
				}
				m.snake.Head.Axis = "y"
				m.snake.Head.Direction = -1
			case "down":
				if m.snake.Head.Axis == "y" && m.snake.Head.Direction == -1 {
					return m, nil
				}
				m.snake.Head.Axis = "y"
				m.snake.Head.Direction = 1
			case "right":
				if m.snake.Head.Axis == "x" && m.snake.Head.Direction == -1 {
					return m, nil
				}
				m.snake.Head.Axis = "x"
				m.snake.Head.Direction = 1
			case "left":
				if m.snake.Head.Axis == "x" && m.snake.Head.Direction == 1 {
					return m, nil
				}
				m.snake.Head.Axis = "x"
				m.snake.Head.Direction = -1
			case " ":
				// pause game
				m.gameState = Paused
				return m, nil
			}
		}

	case TickMsg:

		if m.gameState == Paused {
			return m, nil
		}

		m.Advance()

		// Game over
		if m.snake.HasColided(m.width, m.height) {
			m.gameState = End
			return m, nil
		}

		// Check if the crumb has been eaten
		if m.snake.Head.X == m.crumb.x && m.snake.Head.Y == m.crumb.y {
			m.score++

			// Add new body part
			m.AddBodyPart()

			// Place new crumb at random position
			m.PlaceCrumb()

		}

		// Check if we need to speed up
		if m.snake.Head.Axis == "x" {
			return m, doTick(normalSpeed)
		} else {
			return m, doTick(slowSpeed)
		}

	case tea.WindowSizeMsg:

		// NOTE: Here is the real start of the game
		m.width = msg.Width
		m.height = msg.Height - extraRowsUsed

	}

	return m, nil
}

func (m *model) Advance() {
	// Empty previous position
	prevPos := m.snake.Head
	m.grid[prevPos.Y][prevPos.X] = " "

	// Restrict backwards movement
	if m.snake.Head.Axis == "x" {
		m.snake.Head.X += m.snake.Head.Direction
	}
	if m.snake.Head.Axis == "y" {
		m.snake.Head.Y += m.snake.Head.Direction
	}

	// Fill new position
	m.grid[m.snake.Head.Y][m.snake.Head.X] = "@"

	// Move body
	if len(m.snake.Body) > 0 {
		last := m.snake.Body[len(m.snake.Body)-1]
		m.grid[last.Y][last.X] = " "
	}
	var newBody []game.Position
	for i, p := range m.snake.Body {
		if i == 0 {
			newBody = []game.Position{
				{
					X:         prevPos.X,
					Y:         prevPos.Y,
					Axis:      prevPos.Axis,
					Direction: prevPos.Direction,
					Content:   "o",
				},
			}
		} else {
			newBody = append(newBody, game.Position{
				X:       m.snake.Body[i-1].X,
				Y:       m.snake.Body[i-1].Y,
				Content: p.Content,
			})
		}
	}

	m.snake.Body = newBody
	for _, p := range m.snake.Body {
		m.grid[p.Y][p.X] = p.Content
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

	color := game.NextColor(m.score)
	var crumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	m.grid[crumbY][crumbX] = crumbStyle.Render("x")
	m.crumb = crumb{crumbX, crumbY, crumbStyle, m.crumb.style}
}

func (m *model) AddBodyPart() {
	style := m.crumb.style
	p := game.Position{
		Content: style.Render("o"),
	}
	m.snake.Body = append(m.snake.Body, p)
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

			if col == m.snake.Head.X && row == m.snake.Head.Y {
				m.grid[row][col] = "@"
			}

		}
	}
}

var style2 = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

func (m *model) RestartGame() {
	m.score = 0
	m.snake = game.CreateSnake(m.width, m.height)
	m.FillGrid()
	m.PlaceCrumb()
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
