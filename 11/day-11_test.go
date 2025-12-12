package main

import (
	"io"
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

func TestRecursive(t *testing.T) {
	want := answers[1]

	rack := NewServerRack(data)

	got := rack.NumPathsConnecting("you", "out")

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func BenchmarkLoops(b *testing.B) {
	log.SetOutput(io.Discard)

	b.Run("Stack", func(b *testing.B) {
		rack := NewServerRack(data)
		for b.Loop() {
			rack.NumPathsConnecting("you", "out")
		}
	})

	b.Run("Recursive Cache", func(b *testing.B) {
		rack := NewServerRack(data)
		for b.Loop() {
			rack.RecursivePathsConnecting("you", "out")
		}
	})

	b.Run("Part 2 Stack", func(b *testing.B) {
		rack := NewServerRack(data)
		segments := []struct {
			start, goal string
		}{
			{"svr", "fft"},
			{"fft", "dac"},
			{"dac", "out"},
		}
		for b.Loop() {
			for _, v := range segments {
				rack.NumPathsConnecting(v.start, v.goal)
			}
		}
	})

	b.Run("Part 2 Recursive Cache", func(b *testing.B) {
		rack := NewServerRack(data)
		segments := []struct {
			start, goal string
		}{
			{"svr", "fft"},
			{"fft", "dac"},
			{"dac", "out"},
		}
		for b.Loop() {
			for _, v := range segments {
				rack.RecursivePathsConnecting(v.start, v.goal)
			}
		}
	})
}
