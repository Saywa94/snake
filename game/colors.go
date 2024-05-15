package game

var colors = []string{
	"#7bdff2",
	"#b2f7ef",
	"#eff7f6",
	"#f7d6e0",
	"#f2b5d4",
	"#cfbaf0",
}

func NextColor(score uint) string {
	index := 0
	if score > 30 {
		score = score - 30
	}
	if score > 5 {
		index = 1
	}
	if score > 10 {
		index = 2
	}
	if score > 15 {
		index = 3
	}
	if score > 20 {
		index = 4
	}
	if score > 25 {
		index = 5
	}
	return colors[index]
}
