package utils

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"testing"
)

//
// I asked copilot to make this file, and then I tweaked some parts
//

func TestNewDayUnknownTypePanics(t *testing.T) {
	// I really like defers
	defer func() {
		// check if we had to recover or not
		if r := recover(); r == nil {
			t.Fatalf("expected NewDay[int]() to panic for unsupported type, but it did not")
		}
	}()

	// int is not handled and should panic
	NewDay[int]()
	t.Fatalf("expected panic from NewDay[int](), reached after call")
}

func TestWithReaderAndRead(t *testing.T) {
	want := []string{"line1", "line2"}
	day := NewDay[[]string]().WithReader(func(filename string) []string {
		// return deterministic data for testing
		return want
	})

	got := day.Read("ignored-filename")

	// first time using slices.Equal
	if !slices.Equal(got, want) {
		t.Fatalf("Read returned unexpected data: got=%v", got)
	}
}

func TestRunOutputs(t *testing.T) {
	day := NewDay[[]string]().WithReader(func(filename string) []string {
		return []string{}
	})

	day.Read("ignored")

	// We capture stdout and check for the presence of the expected substrings.
	old := os.Stdout
	// restore stdout at the end?
	defer func() { os.Stdout = old }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run first part
	day.Run(func(data []string) any {
		return "answer1"
	})

	// Run second part
	day.Run(func(data []string) any {
		return 123
	})

	// finish capturing
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)

	out := buf.String()

	if !strings.Contains(out, "1 | answer1") {
		t.Fatalf("output missing first line, got:\n%s", out)
	}
	if !strings.Contains(out, "2 | 123") {
		t.Fatalf("output missing second line, got:\n%s", out)
	}
	if !strings.Contains(out, "Total Time:") {
		t.Fatalf("output missing Total Time summary, got:\n%s", out)
	}
}
