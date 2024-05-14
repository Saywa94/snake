package internal

type Snake struct {
	Head Position
	Body []Position
}

type Position struct {
	X         int
	Y         int
	Axis      string
	Direction int
	Content   string
}

func CreateSnake(width int, height int) Snake {
	return Snake{
		Head: Position{
			X:         width / 2,
			Y:         height/2 - 1,
			Axis:      "x",
			Direction: 1,
			Content:   "@",
		},
		Body: []Position{},
	}
}

func (s *Snake) HasColided(width, height int) bool {
	// Border collision
	colision := false
	if s.Head.X == 0 || s.Head.X == width-1 || s.Head.Y == 0 || s.Head.Y == height-1 {
		colision = true
	}

	for _, p := range s.Body {
		if p.X == s.Head.X && p.Y == s.Head.Y {
			colision = true
		}
	}
	return colision
}
