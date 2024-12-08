package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Map [][]rune
type Position struct {
	X, Y int
}
type Antinode struct {
	Pos Position
}

func main() {
	aMap := ReadMap()
	antinodes := PlaceAntinodes(aMap)
	fmt.Printf("Unique antinodes number: %v\n", len(antinodes))
	// PaintAntinodes(aMap, antinodes, false)
	// PrintMap(aMap)
}

func ReadMap() Map {
	file, err := os.Open("./map.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res := make(Map, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		res = append(res, []rune(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return res
}

func PlaceAntinodes(aMap Map) map[Antinode]bool {
	uniqueAntinodes := make(map[Antinode]bool)

	for i := range aMap {
		for j := range aMap[i] {
			node := aMap[i][j]
			if IsFrequencyNode(node) {
				for k := range aMap {
					for l := range aMap[k] {
						anotherNode := aMap[k][l]
						if anotherNode == node && (i != k || j != l) {
							antinodes := CreateAntinodes(Position{i, j}, Position{k, l}, len(aMap), len(aMap[i]))

							for _, antinode := range antinodes {
								if _, exists := uniqueAntinodes[antinode]; !exists {
									uniqueAntinodes[antinode] = true
								}
							}

						}
					}
				}
			}
		}
	}

	return uniqueAntinodes
}

func IsFrequencyNode(node rune) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]$`)
	return re.MatchString(string(node))
}

func CreateAntinodes(pos1, pos2 Position, boundX, boundY int) []Antinode {
	antinodes := make([]Antinode, 0)
	mul := 0
	for {
		xOffset := (pos2.X - pos1.X) * mul
		yOffset := (pos2.Y - pos1.Y) * mul
		p := Position{X: pos1.X - xOffset, Y: pos1.Y - yOffset}
		if p.X < 0 || p.Y < 0 || p.X >= boundX || p.Y >= boundY {
			break
		}
		antinodes = append(antinodes, Antinode{Pos: p})
		mul++
	}
	return antinodes
}

func PrintMap(aMap Map) {
	fmt.Println()
	for _, row := range aMap {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func PaintAntinodes(aMap Map, antinodes map[Antinode]bool, overwrite bool) {
	for antinode, _ := range antinodes {
		if overwrite {
			aMap[antinode.Pos.X][antinode.Pos.Y] = '#'
		} else {
			node := aMap[antinode.Pos.X][antinode.Pos.Y]
			if !IsFrequencyNode(node) {
				aMap[antinode.Pos.X][antinode.Pos.Y] = '#'
			}
		}
	}
}
