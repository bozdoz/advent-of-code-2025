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

// copied from ScanLines; just swapped "\n" with "," naively (ignoring quoted commas)
func ScanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, ','); i >= 0 {
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
