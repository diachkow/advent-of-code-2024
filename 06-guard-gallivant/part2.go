package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	POINT_BLOCKED  rune = '#'
	POINT_EMPTY         = '.'
	POINT_GUARD         = '^'
	POINT_OBSTACLE      = 'O'
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

type RoutePoint struct {
	Pos Position
	Dir Direction
}
type GuardRoute struct {
	Visited map[RoutePoint]bool
}

func (gr *GuardRoute) Add(rp RoutePoint) {
	gr.Visited[rp] = true
}

func (gr *GuardRoute) AlreadyVisited(rp RoutePoint) bool {
	_, result := gr.Visited[rp]
	return result
}

func NewRoute() *GuardRoute {
	return &GuardRoute{Visited: make(map[RoutePoint]bool)}
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

func CausesGuardLoop(pm PatrolingMap, guard Guard) bool {
	route := NewRoute()

	for {
		rp := RoutePoint{Pos: guard.Pos, Dir: guard.Dir}
		if route.AlreadyVisited(rp) {
			return true
		}
		route.Add(rp)

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
		if nextPoint == POINT_BLOCKED || nextPoint == POINT_OBSTACLE {
			guard.Dir.Turn()
		} else {
			guard.Pos = nextPos
		}
	}

	return false
}

func GetSimulationMap(orig PatrolingMap, obstaclePos Position) PatrolingMap {
	pm := make(PatrolingMap, len(orig))
	for i := range orig {
		pm[i] = make([]rune, len(orig[i]))
		copy(pm[i], orig[i])
	}

	pm[obstaclePos.X][obstaclePos.Y] = POINT_OBSTACLE
	return pm
}

func RunSimulation(pm PatrolingMap) int {
	guard, err := GetGuard(pm)
	if err != nil {
		panic(err)
	}

	var total int
	for i := range pm {
		for j := range pm[i] {
			if pm[i][j] == POINT_EMPTY {
				sm := GetSimulationMap(pm, Position{i, j})
				if CausesGuardLoop(sm, *guard) {
					total++
				}
			}
		}
	}

	return total
}

func main() {
	pm := ReadMap()
	fmt.Printf("Loops with simulation: %v\n", RunSimulation(pm))

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
