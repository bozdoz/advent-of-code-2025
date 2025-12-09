package utils

import (
	"bytes"
	"testing"
)

//
// made by copilot
//

func TestScanCommasBasic(t *testing.T) {
	data := []byte("hello,world,test")

	advance, token, err := ScanCommas(data, false)

	if err != nil {
		t.Fatalf("ScanCommas returned error: %v", err)
	}
	if advance != 6 {
		t.Fatalf("ScanCommas advance wrong: want=6, got=%d", advance)
	}
	if !bytes.Equal(token, []byte("hello")) {
		t.Fatalf("ScanCommas token wrong: want=hello, got=%s", token)
	}
}

func TestScanCommasNoComma(t *testing.T) {
	data := []byte("hello")

	advance, token, err := ScanCommas(data, true)

	if err != nil {
		t.Fatalf("ScanCommas returned error: %v", err)
	}
	if advance != 5 {
		t.Fatalf("ScanCommas advance wrong: want=5, got=%d", advance)
	}
	if !bytes.Equal(token, []byte("hello")) {
		t.Fatalf("ScanCommas token wrong: want=hello, got=%s", token)
	}
}

func TestScanCommasEOFEmpty(t *testing.T) {
	data := []byte("")

	advance, token, err := ScanCommas(data, true)

	if err != nil {
		t.Fatalf("ScanCommas returned error: %v", err)
	}
	if advance != 0 || token != nil {
		t.Fatalf("ScanCommas should return (0, nil) for empty EOF data")
	}
}
