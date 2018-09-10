package main

// #cgo LDFLAGS: -lncurses
// #include <ncurses.h>
//
// // Need this wrapper function because Go can't use variadic C functions
// // (https://stackoverflow.com/questions/26852407/unexpected-type-with-cgo-in-go)
// static void wprint(WINDOW* window, int line, int column, const char* s) {
//   mvwprintw(window, line, column, s);
// }
import "C"
import (
	"math/rand"
	"strconv"
	"time"
)

//noinspection GoVarAndConstTypeMayBeOmitted
func main() {
	C.initscr()
	defer C.endwin()
	C.noecho()
	C.curs_set(C.int(0))
	C.halfdelay(C.int(2))

	var game = &Game{
		width:  20,
		height: 10,
		snake: Snake{
			cells:     []Cell{{4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}},
			direction: right,
		},
		apples: Apples{
			width:       20,
			height:      10,
			growthSpeed: 3,
			random:      rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}

	window := C.newwin(C.int(game.height+2), C.int(game.width+2), 0, 0)
	defer C.delwin(window)

	var c C.int = 0
	for c != 'q' {
		draw(window, game)

		c = C.wgetch(window)
		var direction = none
		switch c {
		case 'i':
			direction = up
		case 'j':
			direction = left
		case 'k':
			direction = down
		case 'l':
			direction = right
		}

		game = game.update(direction)
	}
}

func draw(window *C.WINDOW, game *Game) {
	C.wclear(window)
	C.box(window, 0, 0)

	for _, cell := range game.apples.cells {
		C.wprint(window, C.int(cell.y+1), C.int(cell.x+1), C.CString("."))
	}
	for _, cell := range game.snake.Tail() {
		C.wprint(window, C.int(cell.y+1), C.int(cell.x+1), C.CString("o"))
	}
	head := game.snake.Head()
	C.wprint(window, C.int(head.y+1), C.int(head.x+1), C.CString("Q"))

	if game.IsOver() {
		C.wprint(window, C.int(0), C.int(4), C.CString("Game is Over"))
		C.wprint(window, C.int(1), C.int(3), C.CString("Your score is "+strconv.Itoa(game.Score())))
	}

	C.wrefresh(window)
}

type Game struct {
	width  int
	height int
	snake  Snake
	apples Apples
}

func (game Game) Score() int {
	return len(game.snake.cells)
}
func (game Game) IsOver() bool {
	if contains(game.snake.Tail(), *game.snake.Head()) {
		return true
	}
	for _, cell := range game.snake.cells {
		if cell.x < 0 || cell.x >= game.width || cell.y < 0 || cell.y >= game.height {
			return true
		}
	}
	return false
}
func (game Game) update(direction Direction) *Game {
	if game.IsOver() {
		return &game
	}

	var newSnake, newApples = game.snake.Turn(direction).Move().eat(game.apples.Grow())

	return &Game{
		width:  game.width,
		height: game.height,
		snake:  *newSnake,
		apples: newApples,
	}
}

type Snake struct {
	cells     []Cell
	direction Direction
}

func (snake *Snake) Move() *Snake {
	newHead := snake.Head().Move(snake.direction)
	newTail := snake.cells[:len(snake.cells)-1]
	return &Snake{
		cells:     append([]Cell{newHead}, newTail...),
		direction: snake.direction,
	}
}
func (snake *Snake) Head() *Cell {
	return &snake.cells[0]
}
func (snake *Snake) Tail() []Cell {
	return snake.cells[1:]
}
func (snake *Snake) Turn(newDirection Direction) *Snake {
	if newDirection == none || newDirection.IsOpposite(snake.direction) {
		return snake
	}
	return &Snake{cells: snake.cells, direction: newDirection}
}
func (snake *Snake) eat(apples Apples) (*Snake, Apples) {
	if !contains(apples.cells, *snake.Head()) {
		return snake, apples
	}
	return snake, apples // TODO
}

type Apples struct {
	width       int
	height      int
	cells       []Cell
	growthSpeed int
	random      *rand.Rand
}

func (apples Apples) Grow() Apples {
	if apples.random.Intn(apples.growthSpeed) != 0 {
		return apples
	}
	randomCell := Cell{
		x: apples.random.Intn(apples.width),
		y: apples.random.Intn(apples.height),
	}
	var newCells []Cell
	if !contains(apples.cells, randomCell) {
		newCells = append(apples.cells, randomCell)
	} else {
		newCells = apples.cells
	}

	return Apples{
		width:       apples.width,
		height:      apples.height,
		cells:       newCells,
		growthSpeed: apples.growthSpeed,
		random:      apples.random,
	}
}

type Cell struct {
	x int
	y int
}

func (cell *Cell) Move(direction Direction) Cell {
	switch direction {
	case up:
		return Cell{cell.x, cell.y - 1}
	case down:
		return Cell{cell.x, cell.y + 1}
	case left:
		return Cell{cell.x - 1, cell.y}
	case right:
		return Cell{cell.x + 1, cell.y}
	}
	return *cell
}
func contains(cells []Cell, cell Cell) bool {
	for _, it := range cells {
		if it == cell {
			return true
		}
	}
	return false
}

type Direction int

func (d1 Direction) IsOpposite(d2 Direction) bool {
	if d1 == up && d2 == down {
		return true
	}
	if d2 == up && d1 == down {
		return true
	}
	if d1 == left && d2 == right {
		return true
	}
	if d2 == left && d1 == right {
		return true
	}
	return false
}

const (
	none Direction = iota
	up
	down
	left
	right
)
