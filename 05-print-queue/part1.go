package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func IsValidUpdate(rules map[int][]int, update []int) bool {
	numToIndexMap := make(map[int]int)
	for i, num := range update {
		numToIndexMap[num] = i
	}

	for i := 0; i < len(update); i++ {
		num := update[i]
		numRules, ok := rules[num]

		// number has no associated rules, so it can be considered as valid
		if !ok {
			continue
		}

		for _, numRule := range numRules {
			numRuleIndex, ok := numToIndexMap[numRule]

			// Number from rule is not a part of update, we can skip
			// the check then
			if !ok {
				continue
			}

			if i >= numRuleIndex {
				return false
			}
		}
	}

	return true
}

func GetMiddlePageNumber(update []int) int {
	mid := len(update) / 2.0
	midIndex := math.Ceil(float64(mid))
	return update[int(midIndex)]
}

func main() {
	rules := readRules()
	updates := readUpdates()

	total := 0
	for _, update := range updates {
		if IsValidUpdate(rules, update) {
			num := GetMiddlePageNumber(update)
			total += num
		}
	}

	fmt.Printf("Total is %v\n", total)

}

func readRules() map[int][]int {
	file, err := os.Open("./rules.input.txt")
	panicOnError(err)
	defer file.Close()

	res := make(map[int][]int)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Split(line, "|")

		if len(numbers) != 2 {
			panic("Rule shall be of format X|Y")
		}

		left, err := strconv.Atoi(numbers[0])
		panicOnError(err)

		right, err := strconv.Atoi(numbers[1])
		panicOnError(err)

		res[left] = append(res[left], right)
	}

	if err := scanner.Err(); err != nil {
		panicOnError(err)
	}

	return res
}

func readUpdates() [][]int {
	file, err := os.Open("./updates.input.txt")
	panicOnError(err)
	defer file.Close()

	res := make([][]int, 0)

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		numbers := strings.Split(line, ",")

		lineNumbers := make([]int, len(numbers), len(numbers))

		for j, numberAsString := range numbers {
			number, err := strconv.Atoi(numberAsString)
			panicOnError(err)
			lineNumbers[j] = number
		}

		res = append(res, lineNumbers)
	}

	if err := scanner.Err(); err != nil {
		panicOnError(err)
	}

	return res

}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
