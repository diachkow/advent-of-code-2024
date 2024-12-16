package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Map []string
type Position struct {
	X, Y int
}

type Direction uint8

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type RaindeerState struct {
	Pos   Position
	Dir   Direction
	Score int
	Path  []Position
}

func (rs RaindeerState) Less(other RaindeerState) bool {
	return rs.Score < other.Score
}

const (
	SYMBOL_START rune = 'S'
	SYMBOL_END        = 'E'
	SYMBOL_WALL       = '#'
	SYMBOL_EMPTY      = '.'
)

func main() {
	m := ReadMap()
	PrintMap(m)

	// Get start coordinates
	var rs RaindeerState
	for i, row := range m {
		for j, el := range row {
			if el == SYMBOL_START {
				rs = RaindeerState{
					Pos:   Position{i, j},
					Dir:   EAST,
					Score: 0,
					Path:  []Position{{i, j}},
				}
				break
			}
		}
	}

	score := SearchPath(m, rs)
	fmt.Printf("Lowest score is %v\n", score)
}

func SearchPath(m Map, start RaindeerState) int {
	bestPaths := make([][]Position, 0)
	bestScore := math.MaxInt64
	visitedScores := map[string]int{}

	pq := &PriorityQueue[RaindeerState]{}
	heap.Init(pq)
	heap.Push(pq, start)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(RaindeerState)

		mapKey := fmt.Sprintf("%d:%d:%d", current.Pos.X, current.Pos.Y, current.Dir)

		if prevScore, exists := visitedScores[mapKey]; exists && prevScore < current.Score {
			continue
		}

		visitedScores[mapKey] = current.Score

		if m[current.Pos.X][current.Pos.Y] == SYMBOL_END {
			score := current.Score
			if score < bestScore {
				bestScore = score
				bestPaths = [][]Position{current.Path}
			} else if score == bestScore {
				bestPaths = append(bestPaths, current.Path)
			}
			continue
		}

		// Check next step
		offsetX, offsetY := getOffset(current.Dir)
		newPos := Position{X: current.Pos.X + offsetX, Y: current.Pos.Y + offsetY}
		symbol, err := m.At(newPos)
		if err == nil && symbol != SYMBOL_WALL {
			newPath := make([]Position, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, newPos)

			heap.Push(pq, RaindeerState{
				Pos:   newPos,
				Dir:   current.Dir,
				Score: current.Score + 1,
				Path:  newPath,
			})
		}

		// Check turning clockwise
		heap.Push(pq, RaindeerState{
			Pos:   current.Pos,
			Dir:   (current.Dir + 1) % 4,
			Score: current.Score + 1000,
			Path:  current.Path,
		})

		// Check turning counter-clockwise
		heap.Push(pq, RaindeerState{
			Pos:   current.Pos,
			Dir:   (current.Dir + 3) % 4,
			Score: current.Score + 1000,
			Path:  current.Path,
		})
	}

	uniqueTiles := map[Position]bool{}
	for _, path := range bestPaths {
		for _, pos := range path {
			uniqueTiles[pos] = true
		}
	}

	return len(uniqueTiles)
}

func getOffset(d Direction) (int, int) {
	switch d {
	case NORTH:
		return -1, 0
	case SOUTH:
		return +1, 0
	case WEST:
		return 0, -1
	case EAST:
		return 0, +1
	}
	panic("Unknown direction")
}

func (m Map) At(p Position) (rune, error) {
	if p.X < 0 || p.Y < 0 || p.X >= len(m) || p.Y >= len(m[p.X]) {
		return 'X', fmt.Errorf("Position is out of bounds")
	}

	return rune(m[p.X][p.Y]), nil
}

func ReadMap() Map {
	file, err := os.Open("./map.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	m := make(Map, 0)

	for scanner.Scan() {
		m = append(m, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return m
}

func PrintMap(m Map) {
	for _, row := range m {
		fmt.Println(row)
	}
}

type Comparable[T any] interface {
	Less(other T) bool
}

type PriorityQueue[T Comparable[T]] struct {
	items []T
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.items[i].Less(pq.items[j])
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PriorityQueue[T]) Push(item interface{}) {
	pq.items = append(pq.items, item.(T))
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}
