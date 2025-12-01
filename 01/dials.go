package main

import (
	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Dial struct {
	current int
}

func NewDial() *Dial {
	return &Dial{
		current: 50,
	}
}

func (dial *Dial) Turn(turn string) (clicks int) {
	dir := turn[0]
	val := utils.ParseInt(turn[1:])

	// count number of times around the dial
	clicks = val / 100

	// modulus here to prevent issues with the less than 0 check
	val %= 100

	// check if we're passing 0 or not
	started_zero := dial.current == 0

	if dir == 'R' {
		dial.current += val
	} else {
		dial.current -= val
	}
	// count clicks
	if !started_zero && (dial.current <= 0 || dial.current > 99) {
		clicks++
	}

	// correct negative or >99
	if dial.current < 0 {
		// will not be less than -99
		dial.current = 100 + dial.current
	} else if dial.current > 99 {
		dial.current %= 100
	}

	return
}
