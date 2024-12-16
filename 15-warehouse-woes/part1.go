package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ITEM_WALL  rune = '#'
	ITEM_BOX        = 'O'
	ITEM_ROBOT      = '@'
	ITEM_EMPTY      = '.'
)

type Map [][]rune

const (
	MOVE_UP    rune = '^'
	MOVE_DOWN       = 'v'
	MOVE_LEFT       = '<'
	MOVE_RIGHT      = '>'
)

func main() {
	aMap := ReadMap()
	moves := ReadMoves()

	fmt.Printf("Initial configuration:\n")
	PrintMap(aMap)

	// Get robot position
	var robotPosX, robotPosY int
	for i, row := range aMap {
		for j, el := range row {
			if el == ITEM_ROBOT {
				robotPosX, robotPosY = i, j
			}
		}
	}

	// Execute all move commands
	for _, move := range moves {
		// fmt.Printf("Moving robot '%v' (current position: [%v][%v])\n", string(move), robotPosX, robotPosY)
		robotPosX, robotPosY = MoveRobot(move, robotPosX, robotPosY, aMap)
		// PrintMap(aMap)
	}

	fmt.Printf("Final configuration:\n")
	PrintMap(aMap)

	// Calculate total sum of "GPS" coordinates
	total := 0
	for i, row := range aMap {
		for j, el := range row {
			if el == ITEM_BOX {
				total += 100*i + j
			}
		}
	}

	fmt.Printf("Total sum of GPS coordinates is %v\n", total)
}

func ReadMap() Map {
	file, err := os.Open("./map.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	aMap := make(Map, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]rune, len(line))
		for i, el := range line {
			row[i] = el
		}
		aMap = append(aMap, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return aMap
}

func ReadMoves() []rune {
	bytes, err := os.ReadFile("./moves.input.txt")
	if err != nil {
		panic(err)
	}

	return []rune(strings.TrimSpace(string(bytes)))
}

func PrintMap(aMap Map) {
	fmt.Println()
	for _, row := range aMap {
		for _, el := range row {
			fmt.Printf("%v", string(el))
		}
		fmt.Printf("\n")
	}
	fmt.Println()
}

func MoveRobot(move rune, robotPosX, robotPosY int, aMap Map) (int, int) {
	offsetX, offsetY := 0, 0
	switch move {
	case MOVE_UP:
		offsetX = -1
	case MOVE_DOWN:
		offsetX = +1
	case MOVE_LEFT:
		offsetY = -1
	case MOVE_RIGHT:
		offsetY = +1
	}

	nextX, nextY := robotPosX+offsetX, robotPosY+offsetY
	if !isValidCoord(aMap, nextX, nextY) {
		return robotPosX, robotPosY
	}

	itemInNextCell := aMap[nextX][nextY]
	if itemInNextCell == ITEM_WALL {
		return robotPosX, robotPosY
	}

	if itemInNextCell == ITEM_BOX {
		canMove := false
		nextItemX, nextItemY := nextX, nextY

		for {
			nextItemX += offsetX
			nextItemY += offsetY

			if !isValidCoord(aMap, nextItemX, nextItemY) {
				canMove = false
				break
			}

			nextItem := aMap[nextItemX][nextItemY]
			if nextItem == ITEM_EMPTY {
				canMove = true
				break
			} else if nextItem == ITEM_WALL {
				canMove = false
				break
			}
		}

		if !canMove {
			return robotPosX, robotPosY
		}

		aMap[robotPosX][robotPosY] = ITEM_EMPTY
		aMap[nextX][nextY] = ITEM_ROBOT

		nextItemX, nextItemY = nextX, nextY

		for {
			nextItemX += offsetX
			nextItemY += offsetY

			if !isValidCoord(aMap, nextItemX, nextItemY) {
				break
			}

			nextItem := aMap[nextItemX][nextItemY]

			aMap[nextItemX][nextItemY] = ITEM_BOX

			if nextItem == ITEM_EMPTY {
				break
			}
		}

		return nextX, nextY
	} else {
		aMap[robotPosX][robotPosY] = ITEM_EMPTY
		aMap[nextX][nextY] = ITEM_ROBOT

		return nextX, nextY
	}
}

func isValidCoord(aMap Map, x, y int) bool {
	return x > 0 && y > 0 && x < len(aMap) && y < len(aMap[x])
}
