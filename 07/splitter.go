package main

import "fmt"

type Grid struct {
	start         int // index
	cells         []byte
	height, width int
}

func NewGrid(data []string) *Grid {
	start := 0
	height := len(data)
	width := len(data[0])
	cells := make([]byte, 0, height*width)

	for r := range height {
		for c := range width {
			cell := data[r][c]
			if cell == 'S' {
				start = r*width + c
			}
			cells = append(cells, cell)
		}
	}

	return &Grid{
		start,
		cells,
		height,
		width,
	}
}

// gets last (LIFO)
func pop[T any](s *[]T) *T {
	if len(*s) == 0 {
		return nil
	}

	end := len(*s) - 1
	val := (*s)[end]
	*s = (*s)[:end]

	return &val
}

// gets first (FIFO)
func shift[T any](s *[]T) *T {
	if len(*s) == 0 {
		return nil
	}

	val := (*s)[0]
	*s = (*s)[1:]

	return &val
}

type State struct {
	position   int // index of a splitter
	prev       int // index of previous
	reachesEnd bool
}

type Splitter struct {
	beams         int // usually splits 1
	beamsReachEnd int // out of 2, how many beams
}

// this works for test data, and for part 1,
// but fails for part 2 real input........ :(
func (grid *Grid) Travel() map[int]*Splitter {
	// moves down from start
	// splits on '^'
	states := []State{
		{grid.start, -1, false},
	}
	// visited splitters,
	// how many reach end,
	// how many beams
	splitters := map[int]*Splitter{}

	for len(states) > 0 {
		// breadth-first; still not sure if good idea
		state := shift(&states)
		reachesEnd := state.reachesEnd
		// skip seen
		if splitter, ok := splitters[state.position]; ok {
			if reachesEnd {
				splitter.beamsReachEnd++
			} else {
				// this splitter also sees beams from another splitter above
				splitter.beams += splitters[state.prev].beams
			}
			continue
		}

		beams := 1

		if splitter, ok := splitters[state.prev]; ok {
			beams = splitter.beams
			// fmt.Println("prev splitter", grid.toRC(state.prev), splitter.beams, "cur splitter", grid.toRC(state.position))
		}

		splitters[state.position] = &Splitter{beams: beams}

		if reachesEnd {
			splitters[state.position].beamsReachEnd++
			continue
		}

		// get next states
		states = append(states, grid.GetNextStates(state)...)
	}

	return splitters
}

func (grid *Grid) GetNextStates(state *State) []State {
	// splitter creates 2 new states
	states := []State{}

	// get position
	pos := state.position
	// create 2 new beams to the left and right
	col := pos % grid.width
	can_left := col != 0
	can_right := col != grid.width-1

	// fmt.Println(grid.toRC(pos), col, can_left, can_right)

	// search downward until next splitter for each
	check_side := func(cur int) {
		cur += grid.width

		// while not off the grid
		for cur < grid.height*grid.width {
			if grid.cells[cur] == '^' {
				// found next state
				states = append(states, State{
					position: cur,
					prev:     pos,
				})
				return
			}
			cur += grid.width
		}
		// at end
		states = append(states, State{
			position:   pos, // original position?
			prev:       pos,
			reachesEnd: true,
		})
	}

	if can_left {
		check_side(pos - 1)
	} else {
		// nothing off grid...
		fmt.Println("off grid left")
	}

	if can_right {
		check_side(pos + 1)
	} else {
		// nothing off grid...
		fmt.Println("off grid right")
	}

	return states
}

func (grid *Grid) DebugToRC(i int) [2]int {
	r := i / grid.width
	c := i % grid.width

	return [2]int{r, c}
}

func (grid *Grid) DebugFromRC(rc [2]int) int {
	return rc[0]*grid.width + rc[1]
}

// From: https://www.reddit.com/r/adventofcode/comments/1pgnmou/2025_day_7_lets_visualize/
func (grid *Grid) TryTopDownToArray() int {
	// we can start iterating slowly; the beams can only ever increase the width by 3
	// .S.
	// .^. start is 1, width is 1
	// ^.^ start is 0, width is 3
	start := grid.start + grid.width // one row lower
	width := 1                       // start is 1 beam, 1 column
	// out is slice of columns that *could* have beams exiting
	out := make([]int, grid.width)

	// start with beam going into start column
	out[start%grid.width] = 1

	// while we're not off the grid
	for start < grid.height*grid.width {
		foundSplitter := false
		// iterate start through width looking for splitters
		for i := start; i < start+width; i++ {
			if grid.cells[i] == '^' {
				foundSplitter = true
				// split the beam left and right
				beams := out[i%grid.width]
				// left
				out[(i-1)%grid.width] += beams
				// right
				out[(i+1)%grid.width] += beams
				// empty column
				out[i%grid.width] = 0
			}
		}
		// we can only ever increase by one cell each direction
		if foundSplitter {
			start += grid.width - 1
			width += 2
		} else {
			// next row, same width
			start += grid.width
		}
	}

	sum := 0
	for _, v := range out {
		sum += v
	}

	return sum
}
