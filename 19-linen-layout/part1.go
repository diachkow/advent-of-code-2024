package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	parts, patterns := ReadInputData()

	total := 0
	for _, pattern := range patterns {
		if CanBeMade(pattern, parts) {
			total++
		}
	}

	fmt.Printf("Possible design count: %v\n", total)
}

func CanBeMade(pattern string, parts []string) bool {
	cache := make(map[string]bool)
	return solveRecursive(pattern, parts, cache)
}

func solveRecursive(pattern string, parts []string, cache map[string]bool) bool {
	if len(pattern) == 0 {
		return true
	}

	if result, exists := cache[pattern]; exists {
		return result
	}

	for _, part := range parts {
		if strings.HasPrefix(pattern, part) {
			if solveRecursive(pattern[len(part):], parts, cache) {
				cache[pattern] = true
				return true
			}
		}
	}

	cache[pattern] = false
	return false
}

func ReadInputData() ([]string, []string) {
	file, err := os.Open("./towels.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var parts []string
	patterns := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if len(parts) == 0 {
			parts = strings.Split(line, ", ")
			// slices.SortFunc(parts, func(a, b string) int {
			// 	lenA, lenB := len(a), len(b)
			// 	if lenA > lenB {
			// 		return -1
			// 	} else if lenA == lenB {
			// 		return 0
			// 	} else {
			// 		return 1
			// 	}
			// })
		} else if line == "" {
			continue
		} else {
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return parts, patterns
}
