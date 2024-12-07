package main

import (
	"bufio"
	"fmt"
	"os"
)

func SearchXMASOccurrences(data []string) (int, error) {
	count := 0

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] == 'A' &&
				(ProbeWord(data, i-1, j-1, 1, 1, 0, "MAS") || ProbeWord(data, i-1, j-1, 1, 1, 0, "SAM")) &&
				(ProbeWord(data, i+1, j-1, -1, 1, 0, "MAS") || ProbeWord(data, i+1, j-1, -1, 1, 0, "SAM")) {

				fmt.Printf("Found at [%v][%v]\n", i+1, j+1)
				count += 1
			}
		}
	}

	return count, nil
}

func ProbeWord(data []string, i, j, iOffset, jOffset, wordOffset int, word string) bool {
	// Validate that we did not exceed MAS or data bounderies when adding
	// steps to current "coordinates"
	if wordOffset >= len(word) || i >= len(data) || i < 0 || j >= len(data[i]) || j < 0 {
		return false
	}

	charToSearch := word[wordOffset]
	isCharFound := data[i][j] == charToSearch

	if isCharFound && charToSearch == word[len(word)-1] {
		return true
	} else if isCharFound {
		return ProbeWord(data, i+iOffset, j+jOffset, iOffset, jOffset, wordOffset+1, word)
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
