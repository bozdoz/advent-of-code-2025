package main

import (
	"fmt"
	"sort"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Product struct {
	freshRanges  [][2]int
	availableIDs []int
}

func NewProduct(data []string) *Product {
	freshRanges := [][2]int{}
	availableIDs := []int{}

	i := 0

	for _, item := range data {
		i++
		if item == "" {
			break
		}

		var first, second int
		// ignoring doing split and parse here
		fmt.Sscanf(item, "%d-%d", &first, &second)

		freshRanges = append(freshRanges, [2]int{first, second})
	}

	// ids are the rest
	for _, item := range data[i:] {
		availableIDs = append(availableIDs, utils.ParseInt(item))
	}

	return &Product{
		freshRanges,
		availableIDs,
	}
}

func (prod *Product) CountFresh() (count int) {
	for _, id := range prod.availableIDs {
		for _, rng := range prod.freshRanges {
			min, max := rng[0], rng[1]

			if id >= min && id <= max {
				count++
				break
			}
		}
	}

	return
}

func (prod *Product) sortRanges() {
	sort.Slice(prod.freshRanges, func(i, j int) bool {
		return prod.freshRanges[i][0] < prod.freshRanges[j][0]
	})
}

func (prod *Product) AllPossibleFresh() (count int) {
	// merge ranges, then get diffs
	// sort ranges first
	prod.sortRanges()

	// check in order
	prev := prod.freshRanges[0]

	for i, rng := range prod.freshRanges[1:] {
		if rng[0] <= prev[1] {
			// merge, but check which range max is highest
			// in case one envelopes the other
			if rng[1] > prev[1] {
				prev[1] = rng[1]
			}
		} else {
			// already merged, so we can count the diff
			count += prev[1] - prev[0] + 1
			prev = rng
		}
		// -2 because we omit the first in this loop
		if i == len(prod.freshRanges)-2 {
			count += prev[1] - prev[0] + 1
		}
	}

	return
}
