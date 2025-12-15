package main

import (
	"fmt"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

const (
	FILLED rune = '#'
	EMPTY  rune = '.'
)

type Present struct {
	cells     []rune
	area      int
	cellCount int
}

type Region struct {
	width, height int
	presentCount  []int
}

type Regions struct {
	presents []Present
	regions  []Region
}

func NewRegions(data []string) *Regions {
	presents := []Present{}
	regions := []Region{}

	for _, item := range data[:len(data)-1] {
		present := Present{
			cells: []rune{},
		}
		height := -1
		for line := range strings.SplitSeq(item, "\n") {
			height++
			if height == 0 {
				continue
			}

			for _, val := range line {
				present.cells = append(present.cells, val)
			}
		}
		present.cellCount = strings.Count(item, "#")
		// height * width
		present.area = len(present.cells)
		presents = append(presents, present)
	}
	// each region
	for line := range strings.SplitSeq(data[len(data)-1], "\n") {
		i := strings.Index(line, ": ")

		if i == -1 {
			// extra new line at the end
			break
		}

		region := Region{}

		fmt.Sscanf(line, "%dx%d:", &region.width, &region.height)

		for val := range strings.FieldsSeq(line[i+1:]) {
			region.presentCount = append(region.presentCount, utils.ParseInt(val))
		}
		regions = append(regions, region)
	}

	return &Regions{presents, regions}
}

func (region *Region) DoPresentsFit(presents []Present) bool {
	totalarea := region.height * region.width
	// check if more space than worst case (always works)
	worstcase := 0
	// check if less space than best case (couldn't possibly fit)
	bestcase := 0
	for i, count := range region.presentCount {
		if count > 0 {
			// all cells (3x3)
			worstcase += presents[i].area * count
			bestcase += presents[i].cellCount * count
		}
	}
	if worstcase <= totalarea {
		// we can always fit
		return true
	}
	if bestcase > totalarea {
		// we can never fit
		return false
	}

	panic("Reddit told me this was impossible")
}

func (regions *Regions) HowManyRegionsFit() (count int) {
	for _, region := range regions.regions {
		if region.DoPresentsFit(regions.presents) {
			count++
		}
	}
	return
}
