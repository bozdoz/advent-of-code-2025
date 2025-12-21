package main

import (
	"io"
	"iter"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

// I don't know; just trying it out
type d = iter.Seq[string]

func partOne(data d) (ans any) {
	manual := NewManual(data)

	sum := 0
	for _, machine := range manual.machines {
		sum += machine.FewestButtonPresses()
	}

	return sum
}

func partTwo(data d) (ans any) {
	manual := NewManual(data)

	return utils.SumFunc(manual.machines, func(machine Machine) int {
		return machine.GaussianElimination()
	})
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
