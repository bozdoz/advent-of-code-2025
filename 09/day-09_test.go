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

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 50,
	2: 24,
}

func TestRectArea(t *testing.T) {
	tests := []struct {
		a, b [2]int
		want int
	}{
		{[2]int{2, 5}, [2]int{9, 7}, 24},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test), func(t *testing.T) {
			got := rectArea(test.a, test.b)
			want := test.want

			if got != want {
				t.Errorf("Answer should be %v, but got %v", want, got)
			}
		})
	}
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
