package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Memory line is a obscured string that can contain multiple instructions
type MemoryLine string

// Process a single memory line, returning the result of evaluation instructions
// that were parsed from this memory line
func processMemoryLine(ml MemoryLine) (int, error) {
	var instructionRegexp = regexp.MustCompile(`(mul\(\d+,\d+\)|don't\(\)|do\(\))`)
	tokens := instructionRegexp.FindAllString(string(ml), -1)
	if tokens == nil {
		return 0, errors.New(fmt.Sprintf("No instructions found in memory string %v", ml))
	}

	var numberRegexp = regexp.MustCompile(`\d+`)

	result := 0
	isEnabled := true

	for _, token := range tokens {
		switch {
		case token == "don't()":
			isEnabled = false
		case token == "do()":
			isEnabled = true
		case isEnabled && strings.Contains(token, "mul"):
			numbers := numberRegexp.FindAllString(token, -1)

			left, err := strconv.Atoi(numbers[0])
			if err != nil {
				return 0, err
			}

			right, err := strconv.Atoi(numbers[1])
			if err != nil {
				return 0, err
			}

			result += (left * right)
		}
	}

	return result, nil
}

func main() {
	ml := readInputData()
	result, err := processMemoryLine(ml)
	panicOnError(err)
	fmt.Printf("Result = %v\n", result)
}

func readInputData() MemoryLine {
	b, err := os.ReadFile("./instructions.input.txt")
	panicOnError(err)
	return MemoryLine(string(b))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
