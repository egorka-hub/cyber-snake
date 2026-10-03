package main

import "fmt"

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
		snake:   []Point{{cx, cy}},
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

func dirName(d Point) string {
	switch d {
	case Point{1, 0}:
		return "вправо"
	case Point{-1, 0}:
		return "влево"
	case Point{0, -1}:
		return "вверх"
	case Point{0, 1}:
		return "вниз"
	}
	return "неизвестно"
}

func main() {
	g := NewGame(40, 20)
	head := g.snake[0]
	fmt.Printf("Игра создана: поле %dx%d, змейка в (%d, %d), направление %s, уровень %d\n",
		g.width, g.height, head.x, head.y, dirName(g.dir), g.level)
}
