package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type d = []string

func partOne(data d) (ans any) {
	rack := NewServerRack(data)
	// avoid nothing
	avoid := map[string]struct{}{}

	return rack.NumPathsConnecting("you", "out", avoid)
}

func partTwo(data d) (ans any) {
	rack := NewServerRack(data)

	// get all segment path possibilities

	// segments
	// srv -> [fft -> dac] -> out

	// each segment avoids the other parts
	segments := []struct {
		start, goal string
		avoid       map[string]struct{}
	}{
		{"svr", "fft", map[string]struct{}{
			// avoid the other two
			"dac": {},
			"out": {},
		}},
		{"fft", "dac", map[string]struct{}{
			"svr": {},
			"out": {},
		}},
		// returns 0
		// {"dac", "fft", map[string]struct{}{
		// 	"svr": {},
		// 	"out": {},
		// }},
		// this is the fastest one (15131)
		{"dac", "out", map[string]struct{}{
			"svr": {},
			"fft": {},
		}},
	}

	// is it crazy to do goroutines? YES
	var first, second, third int
	// do goroutines to check segments in parallel:
	response := make(chan int, 3)
	// 240% CPU running all four segments
	// 11 MB memory
	// first two resolve quickly though
	// does this blow up my machine?...looks at butterfly
	// 1.5 seconds to run as go routines
	for _, v := range segments {
		go func() {
			num := rack.NumPathsConnecting(v.start, v.goal, v.avoid)
			if num > 0 {
				// fmt.Println("found", v, num)
				response <- num
			}
		}()
	}

	first = <-response
	second = <-response
	third = <-response

	return first * second * third

	// // NON-GOROUTINES here; it's JUST as fast as goroutines
	// // product
	// num := 1

	// for _, v := range segments {
	// 	num *= rack.NumPathsConnecting(v.start, v.goal, v.avoid)
	// }

	// return num
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
