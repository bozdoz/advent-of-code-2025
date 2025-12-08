package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 21,
	2: 40,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func TestExampleTwo(t *testing.T) {
	want := answers[2]

	got := partTwo(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}

}

//
// copilot made some useless tests for me below
//

func TestSplittersProperties(t *testing.T) {
	grid := NewGrid(data)
	splitters := grid.Travel()

	if len(splitters) == 0 {
		t.Fatal("no splitters found in example data")
	}

	for i, s := range splitters {
		if s.beams < 0 {
			t.Fatalf("splitter %d has negative beams: %d", i, s.beams)
		}
		if s.beamsReachEnd < 0 {
			t.Fatalf("splitter %d has negative beamsReachEnd: %d", i, s.beamsReachEnd)
		}
	}
}

func TestPartTwoManualComputation(t *testing.T) {
	grid := NewGrid(data)
	splitters := grid.Travel()

	expected := 0
	for _, s := range splitters {
		if s.beamsReachEnd > 0 {
			expected += s.beams * s.beamsReachEnd
		}
	}

	got := partTwo(data)
	if got != expected {
		t.Fatalf("partTwo mismatch: got %v, manual expected %v", got, expected)
	}
}

func TestPartOneMatchesTravelCount(t *testing.T) {
	data := day.Read("./example.txt")
	grid := NewGrid(data)
	splitters := grid.Travel()

	got := partOne(data)
	if got != len(splitters) {
		t.Fatalf("partOne (%v) does not match number of splitters (%d)", got, len(splitters))
	}
}

func TestTravelIdempotent(t *testing.T) {
	data := day.Read("./example.txt")
	grid := NewGrid(data)

	a := grid.Travel()
	b := grid.Travel()

	if len(a) != len(b) {
		t.Fatalf("Travel returned different lengths: %d vs %d", len(a), len(b))
	}

	for i := range a {
		if a[i].beams != b[i].beams || a[i].beamsReachEnd != b[i].beamsReachEnd {
			t.Fatalf("splitter %d differs between travels: a=(beams=%d,reached=%d) b=(beams=%d,reached=%d)",
				i, a[i].beams, a[i].beamsReachEnd, b[i].beams, b[i].beamsReachEnd)
		}
	}
}
