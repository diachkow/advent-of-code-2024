package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	arr1, arr2 := readInputData()

	sortArray(arr1)
	sortArray(arr2)

	totalDistance := calculateTotalDistance(arr1, arr2)
	fmt.Printf("Total distance is %v\n", totalDistance)

	totalSimilarityScore := calculateTotalSimilarityScore(arr1, arr2)
	fmt.Printf("Total similarity score is %v\n", totalSimilarityScore)
}

func readInputData() ([]uint32, []uint32) {
	// Read the actual file contents
	f, err := os.Open("./day01.input.txt")
	panicOnError(err)
	defer f.Close()

	// Allocate memory to store location IDs from file
	var arr1, arr2 []uint32

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		left, err1 := strconv.ParseUint(line[:5], 10, 32)
		panicOnError(err1)
		right, err2 := strconv.ParseUint(line[8:13], 10, 32)
		panicOnError(err2)

		arr1 = append(arr1, uint32(left))
		arr2 = append(arr2, uint32(right))
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	return arr1, arr2
}

func calculateTotalDistance(arr1, arr2 []uint32) uint64 {
	res := uint64(0)
	for i := 0; i < len(arr1); i++ {
		res += uint64(subAbs(arr1[i], arr2[i]))
	}
	return res
}

func calculateTotalSimilarityScore(arr1, arr2 []uint32) uint64 {
	occurrences := make(map[uint32]uint32)
	for _, num := range arr2 {
		occurrences[num] += 1
	}

	totalSimilarityScore := uint64(0)
	for i := 0; i < len(arr1); i++ {
		// Similarity score is caulculate as number * number of its occurences in arr1
		num := arr1[i]
		numberOfOccurences := occurrences[num]
		totalSimilarityScore += uint64(num * numberOfOccurences)
	}

	return totalSimilarityScore
}

// Utility functions

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func sortArray(arr []uint32) {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
}

func subAbs(x, y uint32) uint32 {
	if x > y {
		return x - y
	} else {
		return y - x
	}
}
