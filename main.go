package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

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

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	right := g.width + 1
	bottom := g.height + 1
	fg, bg := termbox.ColorDefault, termbox.ColorDefault

	for x := 1; x < right; x++ {
		termbox.SetCell(x, 0, '─', fg, bg)
		termbox.SetCell(x, bottom, '─', fg, bg)
	}

	for y := 1; y < bottom; y++ {
		termbox.SetCell(0, y, '│', fg, bg)
		termbox.SetCell(right, y, '│', fg, bg)
	}

	termbox.SetCell(0, 0, '┌', fg, bg)
	termbox.SetCell(right, 0, '┐', fg, bg)
	termbox.SetCell(0, bottom, '└', fg, bg)
	termbox.SetCell(right, bottom, '┘', fg, bg)

	head := g.snake[0]
	termbox.SetCell(head.x+1, head.y+1, '@', termbox.ColorGreen|termbox.AttrBold, bg)

	info := fmt.Sprintf(" Score: %d Level: %d ", g.score, g.level)
	for i, r := range []rune(info) {
		termbox.SetCell(i+2, 0, r, termbox.ColorYellow, bg)
	}

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	g := NewGame(40, 20)
	g.draw()
	for {
		e := termbox.PollEvent()
		if e.Type == termbox.EventKey && e.Key == termbox.KeyEsc {
			return
		}
	}
}
