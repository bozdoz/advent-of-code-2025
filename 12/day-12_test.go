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

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	defer func() {
		// check if we had to recover or not
		if r := recover(); r == nil {
			t.Fatalf("expected 12.1 to panic because I don't know packing algorithms")
		}
	}()

	// reddit told me actual data is basic, but test data is actually difficult
	partOne(data)
	t.Fatalf("expected panic")
}
