package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Plant rune
type Plot struct {
	X, Y  int
	Plant Plant
}

type Garden struct {
	Plots         []Plot
	Height, Width int
}
type Region struct {
	Plant Plant
	Plots []Plot
}

func main() {
	garden := ReadGardenMap()
	regions := FindRegions(garden)

	total := 0
	for _, region := range regions {
		area := CalculateArea(region)
		sides := CalculateSidesCount(region, garden)
		total += area * sides
	}
	fmt.Printf("Total: %v\n", total)
}

func FindRegions(garden Garden) []Region {
	regs := make([]Region, 0)
	arranged := make(map[Plot]bool)

	for _, plot := range garden.Plots {
		if _, exists := arranged[plot]; !exists {
			regPlots := getRegionPlotsRecursively(garden, plot, 0, 0, arranged)
			regs = append(regs, Region{Plots: regPlots, Plant: plot.Plant})
		}
	}

	return regs
}

var directionOffsets = []struct{ OffsetX, OffsetY int }{
	{-1, 0}, // Up
	{1, 0},  // Down
	{0, -1}, // Left
	{0, 1},  // Right
}

func getRegionPlotsRecursively(garden Garden, plot Plot, offsetX, offsetY int, arranged map[Plot]bool) []Plot {
	nextX, nextY := plot.X+offsetX, plot.Y+offsetY
	newPlot, err := garden.At(nextX, nextY)
	if err != nil {
		return []Plot{}
	}

	if newPlot.Plant != plot.Plant {
		return []Plot{}
	}

	if _, alreadyArranged := arranged[newPlot]; alreadyArranged {
		return []Plot{}
	}

	arranged[newPlot] = true
	plots := []Plot{newPlot}

	for _, dirOffset := range directionOffsets {
		plots = append(plots, getRegionPlotsRecursively(garden, newPlot, dirOffset.OffsetX, dirOffset.OffsetY, arranged)...)
	}

	return plots
}

func (g *Garden) At(x, y int) (Plot, error) {
	for _, plot := range g.Plots {
		if plot.X == x && plot.Y == y {
			return plot, nil
		}
	}
	return Plot{}, fmt.Errorf("Plot was not found")
}

func CalculateArea(region Region) int {
	return len(region.Plots)
}

func CalculateSidesCount(region Region, garden Garden) int {
	if len(region.Plots) == 1 && len(region.Plots) == 2 {
		return 4
	}

	const (
		TOP_LEFT int = iota
		TOP_RIGHT
		BOTTOM_RIGHT
		BOTTOM_LEFT
	)

	const (
		NOT_COUNTS int = iota
		COUNTS
	)

	isSameRegion := func(p1 Plot, offsetX, offsetY int) bool {
		p2, err := garden.At(p1.X+offsetX, p1.Y+offsetY)
		return err == nil && slices.Contains(region.Plots, p2)
	}

	total := 0
	for _, plot := range region.Plots {
		edges := []int{
			COUNTS, // TOP_LEFT
			COUNTS, // TOP_RIGHT
			COUNTS, // BOTTOM_RIGHT
			COUNTS, // BOTTOM_LEFT
		}

		// Check if plot on top is the same figure
		if isSameRegion(plot, -1, 0) {
			edges[TOP_LEFT] = NOT_COUNTS
			edges[TOP_RIGHT] = NOT_COUNTS
		}

		// Check if plot on right is the same figure
		if isSameRegion(plot, 0, 1) {
			edges[TOP_RIGHT] = NOT_COUNTS
			edges[BOTTOM_RIGHT] = NOT_COUNTS
		}

		// Check if plot on bottom is the same figure
		if isSameRegion(plot, 1, 0) {
			edges[BOTTOM_LEFT] = NOT_COUNTS
			edges[BOTTOM_RIGHT] = NOT_COUNTS
		}

		// Check if plot on left is the same figure
		if isSameRegion(plot, 0, -1) {
			edges[TOP_LEFT] = NOT_COUNTS
			edges[BOTTOM_LEFT] = NOT_COUNTS
		}

		// Make some diagonal checks for top edges
		// We don't check for bottom edges not to double calculate them

		// If plot diagonally on top left is within the same region and either left or top plot is missing
		if isSameRegion(plot, -1, -1) && edges[TOP_LEFT] == NOT_COUNTS && (!isSameRegion(plot, -1, 0) || !isSameRegion(plot, 0, -1)) {
			edges[TOP_LEFT] = COUNTS
		}

		// If plot diagonally on top right is within the same region and either right or top plot is missing
		if isSameRegion(plot, -1, +1) && edges[TOP_RIGHT] == NOT_COUNTS && (!isSameRegion(plot, -1, 0) || !isSameRegion(plot, 0, 1)) {
			edges[TOP_RIGHT] = COUNTS
		}

		for _, edge := range edges {
			if edge == COUNTS {
				total++
			}
		}
	}

	return total
}

func ReadGardenMap() Garden {
	file, err := os.Open("./garden.input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	plots := make([]Plot, 0)
	scanner := bufio.NewScanner(file)
	width, height := 0, 0

	for scanner.Scan() {
		row := scanner.Text()
		for i, pl := range row {
			plots = append(plots, Plot{height, i, Plant(pl)})
		}
		width = len(row)
		height++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	garden := Garden{plots, height, width}
	return garden
}
