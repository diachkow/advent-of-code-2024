package main

import (
	"bufio"
	"fmt"
	"os"
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
		perimeter := CalculatePerimeter(region, garden)
		total += area * perimeter
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

func CalculatePerimeter(region Region, garden Garden) int {
	perimeter := 0
	for _, plot := range region.Plots {
		for _, do := range directionOffsets {
			otherPlot, err := garden.At(plot.X+do.OffsetX, plot.Y+do.OffsetY)
			if err != nil || otherPlot.Plant != region.Plant {
				perimeter++
			}
		}
	}
	return perimeter
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
