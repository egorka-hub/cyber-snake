package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

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
	g := &Game{
		snake:   []Point{{cx, cy}},
		food:    Point{},
		malware: []Point{},
		dir:     Point{1, 0},
		score:   0,
		level:   1,
		width:   width,
		height:  height,
		quit:    make(chan struct{}),
	}
	g.placeFood()
	g.placeMalware()
	return g
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
	right := g.width - 1
	bottom := g.height - 1
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
	termbox.SetCell(head.x, head.y, g.dir.ToRune(), termbox.ColorGreen|termbox.AttrBold, bg)
	for _, seg := range g.snake[1:] {
		termbox.SetCell(seg.x, seg.y, '○', termbox.ColorGreen|termbox.AttrBold, bg)
	}

	termbox.SetCell(g.food.x, g.food.y, '●', termbox.ColorCyan, bg)

	for _, m := range g.malware {
		termbox.SetCell(m.x, m.y, '✗', termbox.ColorRed, bg)
	}

	info := fmt.Sprintf(" Score: %d Level: %d ", g.score, g.level)
	for i, r := range []rune(info) {
		termbox.SetCell(i+2, 0, r, termbox.ColorYellow, bg)
	}

	termbox.Flush()
}

func (p Point) ToRune() rune {
	switch p {
	case Point{1, 0}:
		return '▶'
	case Point{-1, 0}:
		return '◀'
	case Point{0, -1}:
		return '▲'
	case Point{0, 1}:
		return '▼'
	}
	return '●'
}

func (p Point) oppositeDir() Point {
	return Point{-p.x, -p.y}
}

func (g *Game) handleInput(ev termbox.Event) {
	var newDir Point
	if ev.Type != termbox.EventKey {
		return
	}
	switch {
	case ev.Key == termbox.KeyArrowUp || ev.Ch == 'w' || ev.Ch == 'W' || ev.Ch == 'ц' || ev.Ch == 'Ц':
		newDir = Point{0, -1}
	case ev.Key == termbox.KeyArrowDown || ev.Ch == 's' || ev.Ch == 'S' || ev.Ch == 'ы' || ev.Ch == 'Ы':
		newDir = Point{0, 1}
	case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a' || ev.Ch == 'A' || ev.Ch == 'ф' || ev.Ch == 'Ф':
		newDir = Point{-1, 0}
	case ev.Key == termbox.KeyArrowRight || ev.Ch == 'd' || ev.Ch == 'D' || ev.Ch == 'в' || ev.Ch == 'В':
		newDir = Point{1, 0}
	case ev.Key == termbox.KeyEsc || ev.Ch == 'q' || ev.Ch == 'Q' || ev.Ch == 'й' || ev.Ch == 'Й':
		select {
		case <-g.quit:
		default:
			close(g.quit)
		}
		return
	default:
		return
	}
	if newDir == g.dir.oppositeDir() {
		return
	}
	g.dir = newDir
}

func (g *Game) isOnSnake(p Point) bool {
	return slices.Contains(g.snake, p)
}

func (g *Game) isOnMalware(p Point) bool {
	return slices.Contains(g.malware, p)
}

func (g *Game) isOutOfBounds(p Point) bool {
	x, y := p.x, p.y
	return x < 1 || x > g.width-2 || y < 1 || y > g.height-2
}

func (g *Game) placeFood() {
	for {
		p := Point{rand.Intn(g.width-2) + 1, rand.Intn(g.height-2) + 1}
		if !g.isOnSnake(p) && !g.isOnMalware(p) {
			g.food = p
			return
		}
	}
}

func (g *Game) placeMalware() {
	for {
		p := Point{rand.Intn(g.width-2) + 1, rand.Intn(g.height-2) + 1}
		if !g.isOnSnake(p) && !g.isOnMalware(p) && p != g.food {
			g.malware = append(g.malware, p)
			return
		}
	}
}

func (g *Game) move() {
	head := g.snake[0]
	newHead := Point{head.x + g.dir.x, head.y + g.dir.y}
	if g.isOnMalware(newHead) || g.isOutOfBounds(newHead) || g.isOnSnake(newHead) {
		g.gameOver = true
		return
	}
	g.snake = append([]Point{newHead}, g.snake...)
	if g.food == newHead {
		g.score++
		if g.score%5 == 0 {
			g.level++
			g.placeMalware()
		}
		g.placeFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

}

func tickInterval(level int) time.Duration {
	interval := max(100*time.Millisecond-time.Duration(level-1)*10*time.Millisecond, 40*time.Millisecond)
	return interval
}

func main() {
	g := NewGame(40, 20)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	g.draw()

	eventCh := make(chan termbox.Event)
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(tickInterval(g.level))

	defer ticker.Stop()
	for {
		select {
		case ev := <-eventCh:
			g.handleInput(ev)
			g.draw()
		case <-ticker.C:
			if !g.gameOver {
				oldLevel := g.level
				g.move()
				if oldLevel != g.level {
					ticker.Reset(tickInterval(g.level))
				}
			}
			g.draw()
		case <-g.quit:
			return
		}
	}

}
