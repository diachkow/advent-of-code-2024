package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LevelsOrder string

const (
	DescendingOrder LevelsOrder = "Descending"
	AscendingOrder  LevelsOrder = "Ascending"
)

type Report struct {
	Levels []uint8
}

func (r Report) String() string {
	return fmt.Sprintf("<Report: %v>", r.Levels)
}

func (r Report) Order() LevelsOrder {
	for i := 0; i < len(r.Levels)-1; i++ {
		if r.Levels[i] < r.Levels[i+1] {
			return AscendingOrder
		} else if r.Levels[i] > r.Levels[i+1] {
			return DescendingOrder
		}
	}

	// If order is not determenied after iterating through all levels, we
	// just default it to asc
	return AscendingOrder
}

// Determine whether report is Safe or not
func (r Report) IsSafe() bool {
	order := r.Order()

	for i := 0; i < len(r.Levels)-1; i++ {
		if order == AscendingOrder {
			if r.Levels[i] > r.Levels[i+1] {
				return false
			} else {
				diff := r.Levels[i+1] - r.Levels[i]
				if diff < 1 || diff > 3 {
					return false
				}
			}
		} else if order == DescendingOrder {
			if r.Levels[i] < r.Levels[i+1] {
				return false
			} else {
				diff := r.Levels[i] - r.Levels[i+1]
				if diff < 1 || diff > 3 {
					return false
				}
			}
		}
	}

	return true
}

type ProblemDampener struct {
	Report Report
}

func (p ProblemDampener) IsSafe() bool {
	if p.Report.IsSafe() {
		return true
	}

	// Pre-allocate memory for experiments with removal of 1 level
	newLevels := make([]uint8, len(p.Report.Levels)-1, len(p.Report.Levels)-1)

	for i := 0; i < len(p.Report.Levels); i++ {
		// Create new report without level at index i
		copy(newLevels, p.Report.Levels[:i])
		copy(newLevels[i:], p.Report.Levels[i+1:])
		newReport := Report{newLevels}

		if newReport.IsSafe() {
			return true
		}
	}

	return false
}

func main() {
	var reports []Report = readInputData()

	safeReportsCount := 0
	for _, report := range reports {
		p := ProblemDampener{report}
		if p.IsSafe() {
			safeReportsCount += 1
		}
	}

	fmt.Printf("Out of %v reports %v are safe\n", len(reports), safeReportsCount)
}

func readInputData() []Report {
	file, err := os.Open("./reports.input.txt")
	panicOnError(err)
	defer file.Close()

	var res []Report

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		reportLine := scanner.Text()
		levelsAsText := strings.Fields(reportLine)

		levels := make([]uint8, len(levelsAsText), len(levelsAsText))

		for i, levelAsText := range levelsAsText {
			level, err := strconv.ParseUint(levelAsText, 10, 8)
			panicOnError(err)
			levels[i] = uint8(level)
		}

		res = append(res, Report{levels})
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
