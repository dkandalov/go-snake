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

//noinspection GoVarAndConstTypeMayBeOmitted
func main() {
	C.initscr()
	defer C.endwin()
	C.noecho()
	C.curs_set(C.int(0))
	C.halfdelay(C.int(2))

	width := 20
	height := 10

	var snake = &Snake{
		cells:     []Cell{{4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}},
		direction: right,
	}

	window := C.newwin(C.int(height+2), C.int(width+2), 0, 0)
	defer C.delwin(window)

	var c C.int = 0
	for c != 'q' {
		C.wclear(window)
		C.box(window, 0, 0)

		for _, cell := range snake.Tail() {
			C.wprint(window, C.int(cell.y+1), C.int(cell.x+1), C.CString("o"))
		}
		C.wprint(window, C.int(snake.Head().y+1), C.int(snake.Head().x+1), C.CString("Q"))

		C.wrefresh(window)

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
		snake = snake.Turn(direction).Move()
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
