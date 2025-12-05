package main

import "iter"

const (
	EMPTY = '.'
	PAPER = '@'
	ERROR = '!' // my own definition
)

type Grid struct {
	height, width int
	cells         []rune
}

func NewGrid(data []string) *Grid {
	height := len(data)
	width := len(data[0])
	cells := make([]rune, height*width)

	for r, row := range data {
		for c, cell := range row {
			cells[r*width+c] = cell
		}
	}

	return &Grid{
		height,
		width,
		cells,
	}
}

func (grid *Grid) getByIndex(index int) rune {
	if index < 0 || index >= len(grid.cells) {
		return ERROR
	}
	return grid.cells[index]
}

// try a sequence
func (grid *Grid) eachNeighbour(i int) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		// top
		top := i - grid.width
		if !yield(grid.getByIndex(top)) {
			return
		}

		// bottom
		bottom := i + grid.width
		if !yield(grid.getByIndex(bottom)) {
			return
		}

		// left-most top-down
		col := i % grid.width
		if col != 0 {
			if !yield(grid.getByIndex(top - 1)) {
				return
			}
			if !yield(grid.getByIndex(i - 1)) {
				return
			}
			if !yield(grid.getByIndex(bottom - 1)) {
				return
			}
		}

		// right-most top-down
		if col != grid.width-1 {
			if !yield(grid.getByIndex(top + 1)) {
				return
			}
			if !yield(grid.getByIndex(i + 1)) {
				return
			}
			if !yield(grid.getByIndex(bottom + 1)) {
				return
			}
		}
	}
}

func (grid *Grid) forEachNeighbour(i int, fun func(char rune)) {
	//     0 1 2
	//
	// 0   0 1 2
	// 1   3 4 5
	// 2   6 7 8
	//
	// r := i / grid.width
	// c := i % grid.width

	// we can get neighbours easier with index

	// top
	top := i - grid.width
	fun(grid.getByIndex(top))

	// bottom
	bottom := i + grid.width
	fun(grid.getByIndex(bottom))

	// left-most top-down
	col := i % grid.width
	if col != 0 {
		fun(grid.getByIndex(top - 1))
		fun(grid.getByIndex(i - 1))
		fun(grid.getByIndex(bottom - 1))
	}

	// right-most top-down
	if col != grid.width-1 {
		fun(grid.getByIndex(top + 1))
		fun(grid.getByIndex(i + 1))
		fun(grid.getByIndex(bottom + 1))
	}
}

func (grid *Grid) CleanPaper() int {
	// this seems slightly faster than a
	// map[int]struct{}
	removable := []int{}

	for cell := range grid.height * grid.width {
		// we only care about the paper ones
		if grid.cells[cell] == PAPER {
			count := 0

			for char := range grid.eachNeighbour(cell) {
				if char == PAPER {
					count++
					if count == 4 {
						// don't keep counting the rest of the 8 neighbours
						break
					}
				}
			}

			// fewer than 4
			if count < 4 {
				// remove it (later)
				removable = append(removable, cell)
			}
		}
	}

	// actually remove
	for _, cell := range removable {
		grid.cells[cell] = EMPTY
	}

	return len(removable)
}
