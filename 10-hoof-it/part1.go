package main

import (
	"bufio"
	"fmt"
	"os"
)

type TopographicMap [][]uint8
type Position struct {
	X, Y int
}

const (
	TRAILHEAD uint8 = 0
	ROUTE_END       = 9
)

func main() {
	tm := ReadMap()
	fmt.Printf("Total is %v\n", RunAllRoutes(tm))
}

func ReadMap() TopographicMap {
	file, err := os.Open("./map.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res := make(TopographicMap, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]uint8, len(line))
		for i, char := range line {
			row[i] = CharToNum(char)
		}
		res = append(res, row)
	}

	return res
}

func CharToNum(char rune) uint8 {
	num := int(char) - '0'
	return uint8(num)
}

func RunAllRoutes(tm TopographicMap) int {
	visited := make([][]bool, len(tm))
	for i, row := range tm {
		visited[i] = make([]bool, len(row))
	}

	totalRoutes := 0

	for i, row := range tm {
		for j, el := range row {
			if el == TRAILHEAD {
				uniqueRouteEnds := make(map[Position]bool)
				RunRoutes(tm, visited, i, j, 0, uniqueRouteEnds)
				fmt.Printf("[%v][%v] has %v routes\n", i, j, len(uniqueRouteEnds))
				totalRoutes += len(uniqueRouteEnds)
			}
		}
	}

	return totalRoutes
}

var directions = []struct{ OffsetX, OffsetY int }{
	{-1, 0}, // Up
	{1, 0},  // Down
	{0, -1}, // Left
	{0, 1},  // Right
}

func RunRoutes(tm TopographicMap, visited [][]bool, posX, posY int, currentVal uint8, uniqueRouteEnds map[Position]bool) {
	if currentVal == ROUTE_END {
		uniqueRouteEnds[Position{posX, posY}] = true
		return
	}

	visited[posX][posY] = true

	for _, dir := range directions {
		nextX, nextY := posX+dir.OffsetX, posY+dir.OffsetY
		nextVal := currentVal + 1
		if IsValidStep(tm, visited, nextX, nextY, nextVal) {
			RunRoutes(tm, visited, nextX, nextY, nextVal, uniqueRouteEnds)
		}
	}

	visited[posX][posY] = false
}

func IsValidStep(tm TopographicMap, visited [][]bool, nextX, nextY int, nextVal uint8) bool {
	// Next cell is out of map bonds
	if nextX < 0 || nextY < 0 || nextX >= len(tm) || nextY >= len(tm[nextX]) {
		return false
	}

	// Next cell is already visited in this route
	if visited[nextX][nextY] {
		return false
	}

	// Next cell is not incremental
	if tm[nextX][nextY] != nextVal {
		return false
	}

	return true
}
