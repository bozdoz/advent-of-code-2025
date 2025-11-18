package utils

import (
	"bufio"
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
