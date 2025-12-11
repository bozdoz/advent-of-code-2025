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
	1: 5,
	2: 2,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

var second = day.Read("./example-2.txt")

func TestExampleTwo(t *testing.T) {
	want := answers[2]

	got := partTwo(second)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}
