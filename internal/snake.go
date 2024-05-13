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
