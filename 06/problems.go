package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

// this is useless
type Operation rune

const (
	ADD Operation = '+'
	MUL Operation = '*'
)

type Problem struct {
	numbers   []int
	operation Operation
}

// maybe unnecessary
type Collection struct {
	problems []Problem
}

func NewCollection(data []string) *Collection {
	problems := []Problem{}

	last := len(data) - 1

	for i, problem := range data[:last] {
		numbers := strings.Fields(problem)

		if i == 0 {
			// initialize problems
			problems = make([]Problem, len(numbers))
		}

		for j, num := range numbers {
			problems[j].numbers = append(problems[j].numbers, utils.ParseInt(num))
		}
	}

	operations := strings.Fields(data[last])

	for i, op := range operations {
		problems[i].operation = Operation(op[0])
	}

	return &Collection{
		problems,
	}
}

func (problem *Problem) DoOperation() (cum int) {
	if problem.operation == MUL {
		// multiplication starts at 1
		// add starts at 0
		cum = 1

		for _, num := range problem.numbers {
			cum *= num
		}
	} else {
		for _, num := range problem.numbers {
			cum += num
		}
	}

	return
}

func NewGridCollection(data []string) *Collection {
	height := len(data)
	width := len(data[0])

	problems := []Problem{}
	problem := Problem{}
	// start top-down, right to left
	for c := width - 1; c >= 0; c-- {
		num := 0
		for r := range height {
			cell := data[r][c]

			// last row
			if r == height-1 {
				// number is always finished
				problem.numbers = append(problem.numbers, num)

				// operation row
				if cell == '+' || cell == '*' {
					// problem is finished
					problem.operation = Operation(cell)
					problems = append(problems, problem)
					// next problem
					problem = Problem{}
					// extra (blank) column skip
					c--
				}
				continue
			}

			if cell != ' ' {
				// is number (simplified from utils.ParseInt)
				num = num*10 + int(cell-'0')
			}
		}
	}

	return &Collection{
		problems,
	}
}
