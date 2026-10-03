package main

type Point struct {
	x, y int
}

type Game struct {
	snake         []Point
	food          Point
	malware       []Point
	dir           Point
	score         int
	level         int
	gameOver      bool
	width, height int
	quit          chan struct{}
}

func NewGame(width, height int) *Game {
	cx, cy := width/2, height/2
	return &Game{
		snake:   []Point{{cx, cy}, {cx - 1, cy}, {cx - 2, cy}},
		food:    Point{cx + 5, cy},
		malware: []Point{},
		dir:     Point{1, 0},
		score:   0,
		level:   1,
		width:   width,
		height:  height,
		quit:    make(chan struct{}),
	}
}

func main() {

}
