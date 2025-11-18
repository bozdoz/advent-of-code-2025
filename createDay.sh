#!/bin/bash

YEAR=2025
NEW_DAY=$1
STRIPPED=${NEW_DAY#0}
NEW_DAY_NAME="day-${NEW_DAY}"

usage() {
    cat >&2 <<END_USAGE

Create a new boilerplate directory from a template

USAGE:
    ./create-day.sh 01
END_USAGE
}

if [ -z $NEW_DAY ]; then
  echo "Provide ## for new day directory"
		usage
  exit 1
fi

kill() {
	echo $1
	exit 1
}

# get session cookie from network tab
# make sure you create this file with:
# COOKIE="Cookie: session=some-cookie-here"
source ./.session.sh || kill

mkdir $NEW_DAY || kill

cd $NEW_DAY

# start touching things
touch example.txt

# get input from API endpoint
curl "https://adventofcode.com/${YEAR}/day/${STRIPPED}/input" --compressed -H "${COOKIE}" > ./input.txt

# create main go file
cat > $NEW_DAY_NAME.go <<EOF
package main

import (
	"io"
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

func init() {
	log.SetOutput(io.Discard)
}

type d = []string

var day = utils.NewDay[d]()

func partOne(data d) (ans any) {
	return 0
}

func partTwo(data d) (ans any) {
	return 0
}

func main() {
	data := day.Read("./input.txt")

	day.Run(partOne, data, 1)
	day.Run(partTwo, data, 2)
}

EOF

# create test file
cat > ${NEW_DAY_NAME}_test.go << EOF
package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Llongfile)
}

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 0,
	2: 0,
}

var data = day.Read("./example.txt")

func TestExampleOne(t *testing.T) {
	want := answers[1]

	got := partOne(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

func TestExampleTwo(t *testing.T) {
	want := answers[2]

	got := partTwo(data)

	if got != want {
		t.Errorf("Answer should be %v, but got %v", want, got)
	}
}

EOF