package main

import (
	"bufio"
	"fmt"
	"os"
)

type Map struct {
	Grid [][]int
}

const (
	MAP_HEIGHT int = 103
	MAP_WIDTH      = 101
	ITERATIONS     = 100
)

type Quadrant uint8

const (
	TOP_LEFT_QUADRANT Quadrant = iota
	TOP_RIGHT_QUADRANT
	BOTTOM_RIGHT_QUADRANT
	BOTTOM_LEFT_QUADRANT
)

type Position struct {
	X, Y int
}
type Velocity struct {
	X, Y int
}

type Robot struct {
	P   struct{ X, Y int }
	V   struct{ X, Y int }
	Map *Map
}

func main() {
	robots, aMap := ReadRobotsData()
	for i := 1; true; i++ {
		for _, robot := range robots {
			robot.Move()
		}

		// Christmas tree ASCII art has a 3x3 square as a tree foot
		for x := range MAP_HEIGHT - 3 {
			for y := range MAP_WIDTH - 3 {
				squareFound := true
				for dx := range 3 {
					for dy := range 3 {
						if aMap.Grid[x+dx][y+dy] == 0 {
							squareFound = false
							break
						}
					}
				}

				if squareFound {
					aMap.Print(true)
					fmt.Printf("\n\nIt took %v seconds to found it!\n", i)
					return
				}
			}
		}
	}
}

func ReadRobotsData() ([]*Robot, *Map) {
	file, err := os.Open("./robots.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mapPtr := NewMap()
	robots := make([]*Robot, 0)

	for scanner.Scan() {
		robot := Robot{Map: mapPtr}
		_, parseErr := fmt.Sscanf(
			scanner.Text(),
			"p=%d,%d v=%d,%d",
			&robot.P.Y, &robot.P.X,
			&robot.V.Y, &robot.V.X,
		)
		if parseErr != nil {
			panic(parseErr)
		}
		mapPtr.Grid[robot.P.X][robot.P.Y]++
		robots = append(robots, &robot)
	}

	return robots, mapPtr
}

func NewMap() *Map {
	grid := make([][]int, 0, MAP_HEIGHT)
	for i := 0; i < MAP_HEIGHT; i++ {
		grid = append(grid, make([]int, MAP_WIDTH))
	}
	return &Map{Grid: grid}
}

func (m *Map) Print(includeMiddle bool) {
	for i, row := range m.Grid {
		for j, el := range row {
			if !includeMiddle && (len(m.Grid)/2 == i || len(m.Grid[i])/2 == j) {
				fmt.Printf(" ")
			} else if el == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%v", el)
			}
		}
		fmt.Printf("\n")
	}
}

func (r *Robot) Move() {
	nextX, nextY := r.P.X+r.V.X, r.P.Y+r.V.Y

	if nextX < 0 {
		offset := 0 - nextX
		nextX = MAP_HEIGHT - offset
	} else if nextX >= MAP_HEIGHT {
		offset := nextX - MAP_HEIGHT
		nextX = offset
	}

	if nextY < 0 {
		offset := 0 - nextY
		nextY = MAP_WIDTH - offset
	} else if nextY >= MAP_WIDTH {
		offset := nextY - MAP_WIDTH
		nextY = offset
	}

	r.Map.Grid[r.P.X][r.P.Y]--
	r.P = Position{nextX, nextY}
	r.Map.Grid[r.P.X][r.P.Y]++
}

func (r Robot) Quadrant() (Quadrant, error) {
	midHeight := MAP_HEIGHT / 2
	midWidth := MAP_WIDTH / 2

	if r.P.X == midHeight && r.P.Y == midWidth {
		return 0, fmt.Errorf("Robot is not within any quadrant")
	}

	if r.P.X < midHeight && r.P.Y < midWidth {
		return TOP_LEFT_QUADRANT, nil
	} else if r.P.X > midHeight && r.P.Y < midWidth {
		return BOTTOM_LEFT_QUADRANT, nil
	} else if r.P.X < midHeight && r.P.Y > midWidth {
		return TOP_RIGHT_QUADRANT, nil
	} else if r.P.X > midHeight && r.P.Y > midWidth {
		return BOTTOM_RIGHT_QUADRANT, nil
	}

	return 0, fmt.Errorf("Failed to detect Quadrant")
}
