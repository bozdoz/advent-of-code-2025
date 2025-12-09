package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	ConnectionCount = 10
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 40,
	2: 25272,
}

var data = day.Read("./example.txt")

func TestNearest(t *testing.T) {
	junction := NewJunction(data)

	junction.GetDistances()

	nearest := junction.distances[0]

	want1 := Point{162, 817, 812}
	want2 := Point{425, 690, 689}

	got1 := junction.points[nearest.a]
	got2 := junction.points[nearest.b]

	if !want1.Equals(got1) {
		t.Errorf("Got %v, want %v", got1, want1)
	}

	if !want2.Equals(got2) {
		t.Errorf("Got %v, want %v", got2, want2)
	}
}

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
