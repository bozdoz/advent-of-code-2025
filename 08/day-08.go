package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

// changed in tests
var ConnectionCount = 1000

func partOne(data d) (ans any) {
	junction := NewJunction(data)

	junction.GetDistances()
	// connect nodes
	junction.ConnectNodes(ConnectionCount)
	// iterate those nodes.... (check if visited)
	// return top 3 largest circuits

	product := 1
	for _, size := range junction.GetTopThreeCircuits(ConnectionCount) {
		product *= size
	}

	return product
}

func partTwo(data d) (ans any) {
	junction := NewJunction(data)

	junction.GetDistances()

	lastA, lastB := junction.ConnectUntilAllJoined()

	return lastA.x * lastB.x
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
