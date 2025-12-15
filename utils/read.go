package utils

import (
	"bufio"
	"bytes"
	"iter"
	"log"
	"os"
	"strings"
)

func ReadAsLines(filename string) (out []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	return
}

// copied from ScanLines;
//
// just swapped "\n" with "," naively (ignoring quoted commas)
//
// then checked for newline within the comma check
//
// to handle 7,1\n11,1
func ScanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, ','); i >= 0 {
		// check newlines too
		if j := bytes.IndexByte(data, '\n'); j >= 0 {
			// if we have both comma and a newline, do the lesser
			// first time using `min`?
			k := min(i, j)
			// We have a full comma-separated segment.
			return k + 1, data[0:k], nil
		}
		// We have a full comma-separated segment.
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

func ReadCSV(filename string) (out []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(ScanCommas)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	return
}

func ReadCSVInt(filename string) (out []int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(ScanCommas)

	for scanner.Scan() {
		out = append(out, ParseInt(scanner.Text()))
	}

	return
}

func ReadSpaceSeparatedSections(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(data), "\n\n")
}

// no idea how to do this
func ReadLinesSeq(filename string) func() iter.Seq[string] {
	return func() iter.Seq[string] {
		return _readLinesSeq(filename)
	}
}

// hard to return an iterator here
func _readLinesSeq(filename string) iter.Seq[string] {
	return func(yield func(string) bool) {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		// ! am I closing too soon?
		// can I read a file inside an iterator?
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
