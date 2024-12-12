package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type StonesArrangement []int

func main() {
	stones := ReadStonesInitialArrangement()
	rollsCount := 75
	fmt.Printf("Number of stones after %v rolls: %v\n", rollsCount, RollStones(stones, rollsCount))
}

func ReadStonesInitialArrangement() StonesArrangement {
	content, err := os.ReadFile("./stones.input.txt")
	if err != nil {
		panic(err)
	}

	values := strings.Fields(string(content))
	stones := make(StonesArrangement, len(values))

	for i, strVal := range values {
		num, err := strconv.Atoi(strVal)
		if err != nil {
			panic(err)
		}
		stones[i] = num
	}

	return stones
}

// Stole this from https://github.com/dbut2/advent-of-code/blob/main/2024/11/2.go
// Never would have though of storing values as map (collection of unique numbers on this iteration)
func RollStones(sa StonesArrangement, rollsCount int) int {
	values := make(map[int]int)
	for _, num := range sa {
		values[num]++
	}

	cache := make(map[int][]int)
	cache[0] = []int{1}

	for i := 0; i < rollsCount; i++ {
		next := make(map[int]int)

		for number, amount := range values {
			if _, exists := cache[number]; !exists {
				if HasEvenNumberOfDigits(number) {
					left, right := SplitDigits(number)
					cache[number] = []int{left, right}
				} else {
					cache[number] = []int{number * 2024}
				}
			}

			for _, val := range cache[number] {
				next[val] += amount
			}
		}

		values = next
	}

	total := 0
	for _, count := range values {
		total += count
	}
	return total
}

func HasEvenNumberOfDigits(val int) bool {
	return len(strconv.Itoa(val))%2 == 0
}

func SplitDigits(val int) (int, int) {
	strVal := strconv.Itoa(val)

	leftNumEncoded := strVal[:len(strVal)/2]
	left, err := strconv.Atoi(leftNumEncoded)
	if err != nil {
		panic(err)
	}

	rightNumEncoded := strVal[len(strVal)/2:]
	right, err := strconv.Atoi(rightNumEncoded)
	if err != nil {
		panic(err)
	}

	return left, right
}
