package utils

import (
	"fmt"
	"iter"
	"time"
)

// not public, so we make sure we use `NewDay` constructor
type day[T any] struct {
	// Check the type of the generic to choose a `reader`
	__type__ T
	// I want to time the whole process
	start time.Time
	// how to read today's puzzle input
	reader func(string) T
	// store the data from the reader (should be pointer?)
	data T
	// increment to print times for each part
	part int
}

// it's a new day
func NewDay[T any]() *day[T] {
	start := time.Now()

	day := &day[T]{
		start: start,
	}

	// figure out what type we just created
	switch any(day.__type__).(type) {
	case []string:
		day.reader = any(ReadAsLines).(func(string) T)
	case []int:
		day.reader = any(ReadCSVInt).(func(string) T)
	case iter.Seq[string]:
		day.reader = any(ReadLinesIter).(func(string) T)
	default:
		panic(fmt.Sprintf("need to implement a reader for type: %T", day.__type__))
	}

	return day
}

// adds some custom reader
func (day *day[T]) WithReader(fun func(string) T) *day[T] {
	day.reader = fun
	return day
}

// read the data according to whatever parser we set
func (day *day[T]) Read(filename string) (out T) {
	day.data = day.reader(filename)

	return day.data
}

// run partone or parttwo
func (day *day[T]) Run(fun func(T) any) {
	begin := time.Now()

	ans := fun(day.data)

	day.part++

	fmt.Printf("%v | %v (%v)\n", day.part, ans, time.Since(begin))

	if day.part == 2 {
		// final time
		fmt.Printf("Total Time: (%v)\n", time.Since(day.start))
	}
}
