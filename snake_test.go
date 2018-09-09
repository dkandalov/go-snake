package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSnakeMovesRight(t *testing.T) {
	var snake = &Snake{
		cells:     []Cell{{2, 0}, {1, 0}, {0, 0}},
		direction: right,
	}
	actualSnake := snake.Move()
	expectedSnake := &Snake{
		cells:     []Cell{{3, 0}, {2, 0}, {1, 0}},
		direction: right,
	}
	if !reflect.DeepEqual(actualSnake, expectedSnake) {
		fmt.Print(actualSnake)
		t.Fail()
	}
}

func TestSnakeChangesDirection(t *testing.T) {
	var snake = &Snake{
		cells:     []Cell{{2, 0}, {1, 0}, {0, 0}},
		direction: right,
	}
	actualSnake := snake.Turn(down).Move()
	expectedSnake := &Snake{
		cells:     []Cell{{2, 1}, {2, 0}, {1, 0}},
		direction: down,
	}
	if !reflect.DeepEqual(actualSnake, expectedSnake) {
		fmt.Print(actualSnake)
		t.Fail()
	}
}
