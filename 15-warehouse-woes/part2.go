package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ItemType uint8

const (
	ITEM_WALL ItemType = iota
	ITEM_BOX
	ITEM_ROBOT
	ITEM_EMPTY
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
			if el == '@' {
				robotPosX, robotPosY = i, j
			}
		}
	}

	// Execute all move commands
	for _, move := range moves {
		// fmt.Printf("Moving robot '%v' (current position: [%v][%v]) (move #%v)\n", string(move), robotPosX, robotPosY, i+1)
		robotPosX, robotPosY = MoveRobot(move, robotPosX, robotPosY, aMap)
	}

	fmt.Printf("Final configuration:\n")
	PrintMap(aMap)

	// Calculate total sum of "GPS" coordinates
	total := 0
	for i, row := range aMap {
		for j, el := range row {
			if el == '[' {
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
		row := make([]rune, 0)
		for _, el := range line {
			switch el {
			case '#':
				row = append(row, '#', '#')
			case 'O':
				row = append(row, '[', ']')
			case '.':
				row = append(row, '.', '.')
			case '@':
				row = append(row, '@', '.')
			}
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
	if !aMap.IsValidCoord(nextX, nextY) {
		return robotPosX, robotPosY
	}

	nextCellItemType, err := aMap.ItemAt(nextX, nextY)
	if err != nil {
		panic(err)
	}

	if nextCellItemType == ITEM_WALL {
		return robotPosX, robotPosY
	}

	if nextCellItemType == ITEM_BOX && (move == MOVE_LEFT || move == MOVE_RIGHT) {
		canMove := false
		nextItemX, nextItemY := nextX, nextY

		for {
			nextItemX += offsetX
			nextItemY += offsetY

			nextItemType, err := aMap.ItemAt(nextItemX, nextItemY)
			if err != nil {
				canMove = false
				break
			}

			if nextItemType == ITEM_EMPTY {
				canMove = true
				break
			} else if nextItemType == ITEM_WALL {
				canMove = false
				break
			}
		}

		if !canMove {
			return robotPosX, robotPosY
		}

		nextItemX, nextItemY = nextX, nextY
		for {
			nextItemX += offsetX
			nextItemY += offsetY

			if !aMap.IsValidCoord(nextItemX, nextItemY) {
				break
			}

			shouldStop := aMap[nextItemX][nextItemY] == '.'

			if move == MOVE_LEFT {
				aMap[nextItemX][nextItemY] = ']'
			} else if move == MOVE_RIGHT {
				aMap[nextItemX][nextItemY] = '['
			} else {
				panic(fmt.Sprintf("Unexpected move %v", string(move)))
			}

			nextItemX += offsetX
			nextItemY += offsetY

			if shouldStop || aMap[nextItemX][nextItemY] == '.' {
				shouldStop = true
			}

			if move == MOVE_LEFT {
				aMap[nextItemX][nextItemY] = '['
			} else {
				aMap[nextItemX][nextItemY] = ']'
			}

			if shouldStop {
				break
			}
		}

		aMap[robotPosX][robotPosY] = '.'
		aMap[nextX][nextY] = '@'

		return nextX, nextY
	} else if nextCellItemType == ITEM_BOX && (move == MOVE_UP || move == MOVE_DOWN) {

		lastRow := robotPosX + offsetX

		type itemToUpdate struct {
			X, Y int
			Item rune
		}

		rowsToUpdate := make(map[int][]itemToUpdate)

		if aMap[nextX][nextY] == '[' {
			rowsToUpdate[lastRow] = []itemToUpdate{
				{X: nextX, Y: nextY, Item: '['},
				{X: nextX, Y: nextY + 1, Item: ']'},
			}
		} else if aMap[nextX][nextY] == ']' {
			rowsToUpdate[lastRow] = []itemToUpdate{
				{X: nextX, Y: nextY, Item: ']'},
				{X: nextX, Y: nextY - 1, Item: '['},
			}
		}

		for {
			isLastRow := false
			nextRow := lastRow + offsetX

			for _, item := range rowsToUpdate[lastRow] {
				nextItemX, nextItemY := item.X+offsetX, item.Y+offsetY
				nextItem := aMap[nextItemX][nextItemY]

				if nextItem == '#' {
					// If wall is detected that means robot cannot move further
					return robotPosX, robotPosY
				} else if nextItem == '[' {
					rowsToUpdate[nextRow] = append(rowsToUpdate[nextRow], []itemToUpdate{
						{X: nextItemX, Y: nextItemY, Item: '['},
						{X: nextItemX, Y: nextItemY + 1, Item: ']'},
					}...)
				} else if nextItem == ']' {
					rowsToUpdate[nextRow] = append(rowsToUpdate[nextRow], []itemToUpdate{
						{X: nextItemX, Y: nextItemY, Item: ']'},
						{X: nextItemX, Y: nextItemY - 1, Item: '['},
					}...)
				}
			}

			if isLastRow {
				break
			}

			if len(rowsToUpdate[nextRow]) == 0 {
				break
			}

			lastRow = nextRow
		}

		for i := lastRow; i != robotPosX; i -= offsetX {
			for _, item := range rowsToUpdate[i] {
				aMap[item.X][item.Y] = '.'
				aMap[item.X+offsetX][item.Y] = item.Item
			}
		}

		aMap[robotPosX][robotPosY] = '.'
		aMap[nextX][nextY] = '@'

		return nextX, nextY
	} else {
		aMap[robotPosX][robotPosY] = '.'
		aMap[nextX][nextY] = '@'

		return nextX, nextY
	}
}

func (m Map) ItemAt(x, y int) (ItemType, error) {
	if !m.IsValidCoord(x, y) {
		return 0, fmt.Errorf("Item type cannot be obtained")
	}

	cell := m[x][y]
	if cell == '#' {
		return ITEM_WALL, nil
	} else if cell == '.' {
		return ITEM_EMPTY, nil
	} else if cell == '@' {
		return ITEM_ROBOT, nil
	} else if cell == '[' || cell == ']' {
		return ITEM_BOX, nil
	}

	return 0, fmt.Errorf("Unexpected cell type %v", string(cell))
}

func (m Map) IsValidCoord(x, y int) bool {
	return x > 0 && y > 0 && x < len(m) && y < len(m[x])
}
