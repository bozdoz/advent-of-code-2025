package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	grid := NewGrid(data)
	return grid.CleanPaper()
}

func partTwo(data d) (ans any) {
	grid := NewGrid(data)

	sum := 0

	for {
		removed := grid.CleanPaper()

		if removed == 0 {
			break
		}

		sum += removed
	}

	return sum
}

//
// BOILERPLATE
//

func init() {
	log.SetOutput(io.Discard)
}

var day = utils.NewDay[d]()

func main() {
	day.Read("./input.txt")

	day.Run(partOne)
	day.Run(partTwo)
}
