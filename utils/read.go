package utils

import (
	"bufio"
	"bytes"
	"log"
	"os"
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
