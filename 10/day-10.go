package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	manual := NewManual(data)

	sum := 0
	for _, machine := range manual.machines {
		sum += machine.FewestButtonPresses()
	}

	return sum
}

func partTwo(data d) (ans any) {
	return 0
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
