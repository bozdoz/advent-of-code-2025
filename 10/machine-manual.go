package main

import (
	"math/bits"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Machine struct {
	light_diagram int         // (reversed) binary
	button_wiring map[int]int // index binary -> button binary
	joltage       []int       // part 2 button press counts
}

type Manual struct {
	machines []Machine
}

func NewManual(data []string) *Manual {
	machines := make([]Machine, len(data))

	for i, line := range data {
		// need to update the actual machine
		machine := &machines[i]
		// need to initialize maps
		machine.button_wiring = map[int]int{}
		// inner index (for buttons)
		j := 0

		for part := range strings.FieldsSeq(line) {
			j++
			inside := part[1 : len(part)-1]
			switch {
			case strings.HasPrefix(part, "["):
				// do this in reverse,
				// so the button ints line up properly
				for j := len(inside) - 1; j >= 0; j-- {
					machine.light_diagram <<= 1
					if inside[j] == '#' {
						machine.light_diagram += 1
					}
				}
			case strings.HasPrefix(part, "("):
				button := 0
				for val := range strings.SplitSeq(inside, ",") {
					// 3 is not 3, but instead the 3rd bit
					button |= 1 << utils.ParseInt(val)
				}
				// i will start at 1 here, because buttons come second
				index_binary := 1 << (j - 1)
				machine.button_wiring[index_binary] = button
			case strings.HasPrefix(part, "{"):
				// TODO: part 2
				machine.joltage = append(machine.joltage, -1)
			}
		}
	}

	return &Manual{machines}
}

func (machine *Machine) FewestButtonPresses() (fewest int) {
	// after analyzing the input:
	// I can say *none* of them will have a diagram of 0

	// state is [2]int = current & pressed
	queue := [][2]int{{0, 0}}

	for len(queue) > 0 {
		item := utils.MustShift(&queue)

		current, pressed := item[0], item[1]

		// get next buttons to press
		for k, v := range machine.button_wiring {
			if k&pressed == k {
				// already pressed
				continue
			}
			// only press if it will impact either current or the diagram
			compare := current | machine.light_diagram

			if v&compare != 0 {
				// queue next state
				// toggles lights
				next := current ^ v
				next_pressed := pressed | k

				if next == machine.light_diagram {
					// wow, we're done
					return bits.OnesCount(uint(next_pressed))
				}

				queue = append(queue, [2]int{next, next_pressed})
			}
		}
	}

	return
}
