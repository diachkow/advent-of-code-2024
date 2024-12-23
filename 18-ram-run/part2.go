package main

import (
	"bufio"
	"fmt"
	"os"
)

// size stand for map size
const SIZE = 71
const N_BYTES = 1024

type Coordinate struct {
	X, Y  int
	Steps int
}

func main() {
	coords := ReadCoordinates()
	aMap := GenerateMap(coords)

	PrintMap(aMap)

	for _, nextByte := range coords[N_BYTES:] {
		fmt.Printf("Trying with (X:%v, Y:%v)\n", nextByte.X, nextByte.Y)
		aMap[nextByte.Y][nextByte.X] = '#'
		// PrintMap(aMap)
		_, err := SearchPath(aMap)
		if err != nil {
			fmt.Printf("Cannot find path with fallen byte (X: %v, Y: %v)\n", nextByte.X, nextByte.Y)
			break
		}
	}
}

func ReadCoordinates() []Coordinate {
	file, err := os.Open("./coords.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	coords := []Coordinate{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		coord := Coordinate{}
		_, parseErr := fmt.Sscanf(
			scanner.Text(),
			"%d,%d",
			&coord.X, &coord.Y,
		)
		if parseErr != nil {
			panic(parseErr)
		}
		coords = append(coords, coord)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return coords
}

func GenerateMap(coords []Coordinate) [][]rune {
	m := make([][]rune, SIZE)
	for i := 0; i < SIZE; i++ {
		m[i] = make([]rune, SIZE)
		for j := 0; j < SIZE; j++ {
			m[i][j] = '.'
		}
	}

	for _, coord := range coords[:N_BYTES] {
		m[coord.Y][coord.X] = '#'
	}

	return m
}

func PrintMap(m [][]rune) {
	fmt.Println()
	for _, row := range m {
		for _, el := range row {
			fmt.Printf("%v", string(el))
		}
		fmt.Println()
	}
	fmt.Println()
}

var directions = [][2]int{
	{+1, 0}, // Right
	{-1, 0}, // Left
	{0, -1}, // Up
	{0, +1}, // Down
}

func SearchPath(m [][]rune) (int, error) {
	start := Coordinate{0, 0, 0}
	endX, endY := SIZE-1, SIZE-1

	queue := []Coordinate{start}
	visited := map[string]bool{}
	visited[fmt.Sprintf("%v:%v", start.X, start.Y)] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.X == endX && curr.Y == endY {
			return curr.Steps, nil
		}

		for _, dir := range directions {
			nextX, nextY := curr.X+dir[0], curr.Y+dir[1]

			if nextX < 0 || nextY < 0 || nextX >= SIZE || nextY >= SIZE {
				continue
			}

			if m[nextY][nextX] == '#' {
				continue
			}

			key := fmt.Sprintf("%v:%v", nextX, nextY)
			if visited[key] {
				continue
			}

			visited[key] = true
			queue = append(queue, Coordinate{nextX, nextY, curr.Steps + 1})
		}
	}

	return 0, fmt.Errorf("No path found")
}
