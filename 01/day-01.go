package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	dial := NewDial()

	count := 0

	for _, line := range data {
		dial.Turn(line)

		if dial.current == 0 {
			count++
		}
	}

	return count
}

func partTwo(data d) (ans any) {
	dial := NewDial()

	count := 0

	for _, line := range data {
		count += dial.Turn(line)
	}

	return count
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
