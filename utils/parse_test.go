package utils

import (
	"strconv"
	"testing"
)

func TestParseInt(t *testing.T) {
	expect := map[string]int{
		"-1":      -1,
		"1":       1,
		"+1":      0,
		"1,000":   1,
		"1000":    1000,
		"3,4,5":   3,
		"  5":     0,
		"	5":      0,
		"garbage": 0,
	}

	for k, v := range expect {
		t.Run(k, func(t *testing.T) {
			got := ParseInt(k)

			if got != v {
				t.Errorf("Wanted: %d; Got: %d", v, got)
			}
		})
	}

	t.Run("panics for empty string", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected a panic, but got %v", r)
			}
		}()

		ParseInt("")
		t.Errorf("Expected to panic, but didn't!")
	})
}

var benchInputs = []string{
	"123",
	"-123",
	"1234567890123",
	"-1234567890123",
}

func BenchmarkParseIntCustom(b *testing.B) {
	for _, s := range benchInputs {
		b.Run(s, func(b *testing.B) {
			for b.Loop() {
				_ = ParseInt(s)
			}
		})
	}
}

func BenchmarkAtoi(b *testing.B) {
	for _, s := range benchInputs {
		b.Run(s, func(b *testing.B) {
			for b.Loop() {
				_, _ = strconv.Atoi(s)
			}
		})
	}
}
