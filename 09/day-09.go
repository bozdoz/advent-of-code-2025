package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

// try new type today
type d = []int

func partOne(data d) (ans any) {
	theater := NewMovieTheater(data)

	return theater.LargestArea()
}

func partTwo(data d) (ans any) {
	theater := NewMovieTheater(data)

	return theater.LargestAreaWithinPerimeter()
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
