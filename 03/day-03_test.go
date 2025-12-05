package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func TestLargest(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"987654321111111", 98},
		{"811111111111119", 89},
		{"234234234234278", 78},
		{"818181911112111", 92},
	}

	// TODO: should this be some utility to stop writing this?
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			jolt := NewJoltage(test.input)
			got := jolt.Largest()

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 357,
	2: 3121910778619,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func TestLargestChain(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"987654321111111", 987654321111},
		{"811111111111119", 811111111119},
		{"234234234234278", 434234234278},
		{"818181911112111", 888911112111},
	}

	// TODO: should this be some utility to stop writing this?
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			jolt := NewJoltage(test.input)
			got := jolt.largestChain()

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

func TestLargestOptimized(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"987654321111111", 987654321111},
		{"811111111111119", 811111111119},
		{"234234234234278", 434234234278},
		{"818181911112111", 888911112111},
	}

	// TODO: should this be some utility to stop writing this?
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			jolt := NewJoltage(test.input)
			got := jolt.LargestOptimized()

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

var benchTests = []string{
	"987654321111111",
	"811111111111119",
}

func BenchmarkLargestChain(b *testing.B) {
	log.SetOutput(io.Discard)
	for _, s := range benchTests {
		jolt := NewJoltage(s)
		b.Run(fmt.Sprintf("%v", s), func(b *testing.B) {
			for b.Loop() {
				jolt.largestChain()
			}
		})
	}
}

func BenchmarkLargestOptimized(b *testing.B) {
	log.SetOutput(io.Discard)
	for _, s := range benchTests {
		jolt := NewJoltage(s)
		b.Run(fmt.Sprintf("%v", s), func(b *testing.B) {
			for b.Loop() {
				jolt.LargestOptimized()
			}
		})
	}
}
