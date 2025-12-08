package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	grid := NewGrid(data)

	return len(grid.Travel())
}

func partTwo(data d) (ans any) {
	grid := NewGrid(data)

	//
	// no idea why this didn't work on my input, but worked on test data
	//
	// splitters := grid.Travel()

	// sum := 0

	// for _, splitter := range splitters {
	// 	if splitter.beamsReachEnd > 0 {
	// 		sum += splitter.beams * splitter.beamsReachEnd
	// 	}
	// }

	return grid.TryTopDownToArray()
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
