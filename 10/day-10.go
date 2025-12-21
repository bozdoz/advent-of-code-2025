package main

import (
	"fmt"
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

	sum := 0
	for i, machine := range manual.machines {
		fewest := machine.GetFewestRecursive()
		fmt.Println("fewest", i+1, fewest)
		sum += fewest
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
