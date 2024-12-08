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
							antinode, err := CreateAntinode(Position{i, j}, Position{k, l}, len(aMap), len(aMap[i]))
							if err != nil {
								continue
							}

							if _, exists := uniqueAntinodes[*antinode]; !exists {
								uniqueAntinodes[*antinode] = true
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

func CreateAntinode(pos1, pos2 Position, boundX, boundY int) (*Antinode, error) {
	p := Position{X: pos1.X - (pos2.X - pos1.X), Y: pos1.Y - (pos2.Y - pos1.Y)}
	if p.X < 0 || p.Y < 0 || p.X >= boundX || p.Y >= boundY {
		return nil, fmt.Errorf("Cannot create Antinode: Position is out of map bounds")
	}
	return &Antinode{Pos: p}, nil
}

// func PrintMap(aMap Map) {
// 	fmt.Println()
// 	for _, row := range aMap {
// 		fmt.Println(string(row))
// 	}
// 	fmt.Println()
// }
//
// func PaintAntinodes(aMap Map, antinodes map[Antinode]bool) {
// 	for antinode, _ := range antinodes {
// 		node := aMap[antinode.Pos.X][antinode.Pos.Y]
// 		if !IsFrequencyNode(node) {
// 			aMap[antinode.Pos.X][antinode.Pos.Y] = '#'
// 		}
// 	}
// }
