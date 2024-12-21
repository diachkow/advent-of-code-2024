package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	parts, patterns := ReadInputData()

	fmt.Printf("Parts: %v\n", parts)
	fmt.Printf("Patterns: %v\n", patterns)

	total := 0
	for _, pattern := range patterns {
		if variants := CanBeMade(pattern, parts); variants > 0 {
			fmt.Printf("Pattern '%v' can be made using %v variants\n", pattern, variants)
			total += variants
		} else {
			fmt.Printf("Pattern '%v' cannot be made\n", pattern)
		}
	}

	fmt.Printf("Possible design count: %v\n", total)
}

func CanBeMade(pattern string, parts []string) int {
	cache := make(map[string]int)
	return solveRecursive(pattern, parts, cache)
}

func solveRecursive(pattern string, parts []string, cache map[string]int) int {
	if len(pattern) == 0 {
		return 1
	}

	if result, exists := cache[pattern]; exists {
		return result
	}

	total := 0
	for _, part := range parts {
		if strings.HasPrefix(pattern, part) {
			total += solveRecursive(pattern[len(part):], parts, cache)
		}
	}

	cache[pattern] = total
	return total
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
