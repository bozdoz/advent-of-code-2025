package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func TestDialTurn(t *testing.T) {
	tests := []struct {
		current, want int
		turn          string
	}{
		{current: 11, turn: "R8", want: 19},
		{current: 19, turn: "L19", want: 0},
		{current: 0, turn: "L1", want: 99},
		{current: 99, turn: "R1", want: 0},
		{current: 5, turn: "L10", want: 95},
		{current: 84, turn: "L884", want: 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			dial := NewDial()
			dial.current = test.current
			dial.Turn(test.turn)

			if dial.current != test.want {
				t.Errorf("got: %v, want: %v", dial.current, test.want)
			}
		})
	}
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 3,
	2: 6,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func TestDialClicks(t *testing.T) {
	tests := []struct {
		current, want int
		turn          string
	}{
		{current: 50, turn: "L82", want: 1},
		{current: 50, turn: "L182", want: 2},
		{current: 50, turn: "R382", want: 4},
		// start at 0 doesn't pass 0
		{current: 0, turn: "L5", want: 0},
		{current: 0, turn: "R100", want: 1},
		{current: 0, turn: "R101", want: 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			dial := NewDial()
			dial.current = test.current
			got := dial.Turn(test.turn)

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

func TestExampleTwo(t *testing.T) {
	want := answers[2]

	got := partTwo(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}
