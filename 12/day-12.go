package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	return NewRegions(data).HowManyRegionsFit()
}

func partTwo(data d) (ans any) {
	return "—"
}

//
// BOILERPLATE
//

func init() {
	log.SetOutput(io.Discard)
}

var day = utils.NewDay[d]().WithReader(utils.ReadSpaceSeparatedSections)

func main() {
	day.Read("./input.txt")

	day.Run(partOne)
	day.Run(partTwo)
}
