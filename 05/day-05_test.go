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
	1: 3,
	2: 14,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func TestSortRanges(t *testing.T) {
	prod := NewProduct(data)
	prod.sortRanges()

	got := prod.freshRanges[2][0]
	want := 12

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}

	// also want 16 to be last
	got = prod.freshRanges[3][0]
	want = 16

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
