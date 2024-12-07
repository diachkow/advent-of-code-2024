package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	POINT_BLOCKED rune = '#'
	POINT_EMPTY        = '.'
	POINT_GUARD        = '^'
	POINT_VISITED      = 'X'
)

type Direction uint8

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d *Direction) Turn() {
	*d = (*d + 1) % 4
}

type Position struct{ X, Y int }
type Guard struct {
	Pos Position
	Dir Direction
}

type PatrolingMap [][]rune

func PrintMap(pm PatrolingMap) {
	fmt.Println()
	for _, row := range pm {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func GetGuard(pm PatrolingMap) (*Guard, error) {
	for i := range pm {
		for j := range pm[i] {
			if pm[i][j] == POINT_GUARD {
				return &Guard{
					Pos: Position{i, j},
					Dir: UP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("Failed to find guard at start")
}

func RunMap(pm PatrolingMap) {
	guard, err := GetGuard(pm)
	if err != nil {
		panic(err)
	}

	for {
		// Paint current guard position as visited
		pm[guard.Pos.X][guard.Pos.Y] = POINT_VISITED

		var nextPos Position
		switch guard.Dir {
		case UP:
			nextPos = Position{guard.Pos.X - 1, guard.Pos.Y}
		case RIGHT:
			nextPos = Position{guard.Pos.X, guard.Pos.Y + 1}
		case DOWN:
			nextPos = Position{guard.Pos.X + 1, guard.Pos.Y}
		case LEFT:
			nextPos = Position{guard.Pos.X, guard.Pos.Y - 1}
		}

		if nextPos.X < 0 || nextPos.Y < 0 || nextPos.X >= len(pm) || nextPos.Y >= len(pm[nextPos.X]) {
			// This means guard has reached his destination at the end
			break
		}

		nextPoint := pm[nextPos.X][nextPos.Y]
		if nextPoint == POINT_BLOCKED {
			guard.Dir.Turn()
		} else {
			guard.Pos = nextPos
		}
	}
}

func CalculateVisited(pm PatrolingMap) int {
	var total int
	for i := range pm {
		for j := range pm[i] {
			if pm[i][j] == POINT_VISITED {
				total++
			}
		}
	}
	return total
}

func main() {
	pm := ReadMap()
	RunMap(pm)
	PrintMap(pm)
	fmt.Printf("Total visited points: %v\n", CalculateVisited(pm))
}

func ReadMap() PatrolingMap {
	file, err := os.Open("./map.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res := make(PatrolingMap, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		res = append(res, []rune(line))
	}

	return res
}
