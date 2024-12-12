package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type StonesArrangement []int

const BLINK_COUNT = 25

func main() {
	stones := ReadStonesInitialArrangement()

	for i := 0; i < BLINK_COUNT; i++ {
		stones = RollStones(stones)
		fmt.Printf("Number of stones after %v roll: %v\n", i+1, len(stones))
	}
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

func RollStones(sa StonesArrangement) StonesArrangement {
	res := make(StonesArrangement, len(sa))
	copy(res, sa)

	for origIndex, copyIndex := 0, 0; origIndex < len(sa); origIndex++ {
		stone := sa[origIndex]

		if stone == 0 {
			res[copyIndex] = 1
		} else if HasEvenNumberOfDigits(stone) {
			left, right := SplitDigits(stone)
			res[copyIndex] = left
			res = slices.Insert(res, copyIndex+1, right)
			copyIndex++
		} else {
			res[copyIndex] = stone * 2024
		}

		copyIndex++
	}

	return res
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
