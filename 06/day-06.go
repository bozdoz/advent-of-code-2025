package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	collection := NewCollection(data)

	sum := 0

	for _, problem := range collection.problems {
		sum += problem.DoOperation()
	}

	return sum
}

func partTwo(data d) (ans any) {
	collection := NewGridCollection(data)

	sum := 0

	for _, problem := range collection.problems {
		sum += problem.DoOperation()
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
