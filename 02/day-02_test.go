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

func TestHasDouble(t *testing.T) {
	tests := []struct {
		num  int
		want bool
	}{
		{11, true},
		{12, false},
		{222, false},
		{2222, true},
		{1188511885, true},
		{1010, true},
		{1011, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			got := hasDouble(test.num)

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

func TestInvalidIds(t *testing.T) {
	tests := []struct {
		rng  string
		want int
	}{
		{"11-22", 2},
		{"95-115", 1},
		{"1188511880-1188511890", 1},
		{"1698522-1698528", 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			product := NewProductRange(test.rng)
			got := product.InvalidCount(hasDouble)

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 1227775554,
	2: 4174379265,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

var repeatingTests = []struct {
	num  int
	want bool
}{
	{1, false},
	{11, true},
	{12, false},
	{20, false},
	{101, false},
	{112, false},
	{222, true},
	{1010, true},
	{1011, false},
	{1212, true},
	{101010, true},
	{101011, false},
	{111011, false},
	{110111, false},
	{101111, false},
	{123123, true},
	{121212, true},
	{121112, false},
	// this broke in my clever method
	{1020102, false},
	{1188511885, true},
	{10000000001, false},
	{101010101010, true},
	{100111001110011, true},
}

func TestHasRepeating(t *testing.T) {
	for _, test := range repeatingTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			got := hasRepeating(test.num)

			if got != test.want {
				t.Errorf("got: %v, want: %v", got, test.want)
			}
		})
	}
}

func TestHasRepeatingString(t *testing.T) {
	for _, test := range repeatingTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			got := hasRepeatingString(test.num)

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

var benchTests = []int{
	111,
	// this broke in my clever method
	1020102,
	123123123123,
}

func BenchmarkMath(b *testing.B) {
	log.SetOutput(io.Discard)
	for _, s := range benchTests {
		b.Run(fmt.Sprintf("%v", s), func(b *testing.B) {
			for b.Loop() {
				_ = hasRepeating(s)
			}
		})
	}
}

// this is slower in benchmarks...
func BenchmarkString(b *testing.B) {
	log.SetOutput(io.Discard)
	for _, s := range benchTests {
		b.Run(fmt.Sprintf("%v", s), func(b *testing.B) {
			for b.Loop() {
				_ = hasRepeatingString(s)
			}
		})
	}
}
