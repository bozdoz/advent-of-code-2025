package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	sum := 0

	for _, input := range data {
		jolt := NewJoltage(input)
		sum += jolt.Largest()
	}

	return sum
}

func partTwo(data d) (ans any) {
	sum := 0

	for _, input := range data {
		jolt := NewJoltage(input)
		sum += jolt.LargestOptimized()
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
