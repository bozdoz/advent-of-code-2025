package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	sum := 0
	for _, rng := range data {
		product := NewProductRange(rng)
		product.InvalidCount(hasDouble)
		sum += product.sum
	}
	return sum
}

func partTwo(data d) (ans any) {
	sum := 0
	for _, rng := range data {
		product := NewProductRange(rng)
		product.InvalidCount(hasRepeatingString)
		sum += product.sum
	}
	return sum
}

//
// BOILERPLATE
//

func init() {
	log.SetOutput(io.Discard)
}

var day = utils.NewDay[d]().WithReader(utils.ReadCSV)

func main() {
	day.Read("./input.txt")

	day.Run(partOne)
	day.Run(partTwo)
}
