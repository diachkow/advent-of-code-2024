package main

import (
	"bufio"
	"fmt"
	"os"
)

const XMAS = "XMAS"

func SearchXMASOccurrences(data []string) (int, error) {
	count := 0

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			// Search forward (move j+1 on each step)
			if ProbeXMAS(data, i, j, 0, 1, 0) {
				count += 1
			}

			// Search backward (move j-1 on each step)
			if ProbeXMAS(data, i, j, 0, -1, 0) {
				count += 1
			}

			// Search down (move i+1 on each step)
			if ProbeXMAS(data, i, j, 1, 0, 0) {
				count += 1
			}

			// Search up (move i-1 on each step)
			if ProbeXMAS(data, i, j, -1, 0, 0) {
				count += 1
			}

			// Search diagonal top left (move i-1, j-1 each step)
			if ProbeXMAS(data, i, j, -1, -1, 0) {
				count += 1
			}

			// Search diagonal top right (move i-1, j+1 on each step)
			if ProbeXMAS(data, i, j, -1, 1, 0) {
				count += 1
			}

			// Search diagonal bottom left (move i+1, j-1 on each step)
			if ProbeXMAS(data, i, j, 1, -1, 0) {
				count += 1
			}

			// Search diagonal bottom right (move i+1, j+1 on each step)
			if ProbeXMAS(data, i, j, 1, 1, 0) {
				count += 1
			}
		}
	}

	return count, nil
}

func ProbeXMAS(data []string, i, j, iOffset, jOffset, wordOffset int) bool {
	// Validate that we did not exceed XMAS or data bounderies when adding
	// steps to current "coordinates"
	if wordOffset >= len(XMAS) || i >= len(data) || i < 0 || j >= len(data[i]) || j < 0 {
		return false
	}

	charToSearch := XMAS[wordOffset]
	isCharFound := data[i][j] == charToSearch

	if isCharFound && charToSearch == XMAS[len(XMAS)-1] {
		return true
	} else if isCharFound {
		return ProbeXMAS(data, i+iOffset, j+jOffset, iOffset, jOffset, wordOffset+1)
	} else {
		return false
	}
}

func main() {
	data := readInputData()

	wordCount, err := SearchXMASOccurrences(data)
	panicOnError(err)
	fmt.Printf("Number of XMAS occurrences: %v\n", wordCount)

}

func readInputData() []string {
	file, err := os.Open("./day04.input.txt")
	panicOnError(err)
	defer file.Close()

	var data []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
