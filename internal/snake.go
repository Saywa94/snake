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

func (s *Snake) Start(width int, height int) {

	s.Head = Position{
		X:         width / 2,
		Y:         height/2 - 1,
		Axis:      "x",
		Direction: 1,
		Content:   "o",
	}
	s.Body = []Position{
		{X: 2, Y: 2, Content: "o"},
	}
}
