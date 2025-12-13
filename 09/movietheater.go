package main

import (
	"iter"
	"maps"
	"slices"
)

type MovieTheater struct {
	redTiles map[[2]int]struct{} // x,y

	// consider saving "start" tile so we can predictably assume
	// either CW or CCW in perimeter building below
}

func NewMovieTheater(data []int) *MovieTheater {
	redTiles := make(map[[2]int]struct{}, len(data))

	// first time using Chunk
	for chunk := range slices.Chunk(data, 2) {
		redTiles[[2]int{chunk[0], chunk[1]}] = struct{}{}
	}

	return &MovieTheater{
		redTiles,
	}
}

func rectArea(a, b [2]int) int {
	height := a[1] - b[1]

	if height < 0 {
		height = -height
	}

	// even with 0 height, we include the col,row that
	// we're in, meaning we need to add 1 to both sides
	height++

	width := a[0] - b[0]

	if width < 0 {
		width = -width
	}

	width++

	return height * width
}

func (theater *MovieTheater) LargestArea() (largest int) {
	// iterate all points and return largest area
	// now it's the first time using Collect
	tiles := slices.Collect(maps.Keys(theater.redTiles))
	var area int
	for i, a := range tiles {
		for _, b := range tiles[i+1:] {
			area = rectArea(a, b)

			// first time using max
			largest = max(largest, area)
		}
	}

	return
}

//
// PART TWO HOCUS POCUS
//

// acceptable diagonals
// each corner has an inside and an outside
// I don't want to consider comparing tiles that
// are aiming outwards, so I track acceptable diagonals
type DIAGONALS uint8

// using bits because a corner could have 1 or 3
const (
	NW DIAGONALS = 1 << iota
	NE
	SE
	SW
)

type Perimeter struct {
	// these include whatever their constant variable is as their 1st value
	//
	// horizontal = y,x1,x2;
	//
	// vertical = x,y1,y2
	horizontal, vertical [][3]int
	// give each tile acceptable diagonals
	tileAndDiagonals map[[2]int]DIAGONALS
}

// indexes of DIRS below
const (
	UP int = iota
	RIGHT
	DOWN
	LEFT
)

var DIRS = [4][2]int{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
}

// moved this to a new function just for cleaner code
func (theater *MovieTheater) eachAdjacentTile() iter.Seq2[[2]int, int] {
	// tried this, and it helped debug issues,
	// but it didn't improve performance, so I dropped it
	// we start from the start for some predictability
	// start := theater.start

	// start from any point,
	// and move in each direction to draw the perimeter lines
	var start [2]int
	// unordered map
	for k := range theater.redTiles {
		start = k
		break
	}

	cur := start

	ignore_dir := -1 // first check every direction, then ignore prev
	dist := 0
	// not sure if declaring variables outside the loop improves perf
	var candidate [2]int

	return func(yield func([2]int, int) bool) {
		// start with start?
		// kind of annoying, because I'd like to do something like iter.Seq3
		// with values: prev, current, direction
		// but only Seq and Seq2 exist, and I don't want to create a struct
		if !yield(start, -1) {
			// gopls is mad at me if I don't if-statement all my yields
			return
		}

		// we need to iterate start twice
		// in order to get the direction it came from
		// because I didn't want to make a struct,
		// explained in previous comment
		final_loop := false

		for {
			// moving outwards from a given tile
			dist++
			// iterate in each direction
			for d, dir := range DIRS {
				// don't check the direction you came
				if d == ignore_dir {
					continue
				}
				// move current by (dist) amount in (dir) direction
				candidate = [2]int{cur[0] + dir[0]*dist, cur[1] + dir[1]*dist}

				// check if candidate is a red tile
				if _, found := theater.redTiles[candidate]; found {
					if !yield(candidate, d) {
						// exit if false (never)
						return
					}
					if final_loop {
						// finally done, after 2nd round with start
						return
					}
					if candidate == start {
						// we're *almost* done (exit)
						// we need to iterate the start one more time because
						// we finally know what direction it came from
						final_loop = true
					}
					// get ready for next loop
					cur = candidate
					// ignore the opposite of this direction
					ignore_dir = (d + 2) % 4
					// reset dist
					dist = 0
					// don't keep iterating directions
					break
				}
			}
		}
	}
}

func (theater *MovieTheater) formPerimeter() *Perimeter {
	horizontal := make([][3]int, 0, len(theater.redTiles))
	vertical := make([][3]int, 0, len(theater.redTiles))

	// prev tile, with which we draw lines
	var prev [2]int
	// d and prev_d *helps* determine which diagonals are inside
	var prev_d int
	// the line we will draw
	var line [3]int

	// we'll assume it's CCW and left-side is inside,
	// otherwise, we'll need to inverse
	// TIL making maps can take an optional *hint* at its size:
	// https://go.dev/ref/spec#Making_slices_maps_and_channels
	diagonals := make(map[[2]int]DIAGONALS, len(theater.redTiles))

	// shoelace theorem gives negative for CCW and positive for CW
	// we don't actually care about the area value
	shoelace_area := 0

	for cur, d := range theater.eachAdjacentTile() {
		// `d` starts at -1 (unfortunately)
		if d != -1 {
			// draw line, checking direction with switch
			switch d {
			case UP, DOWN:
				// first int is the shared val
				// the rest are min/max
				line = [3]int{cur[0], cur[1], prev[1]}
				// is there a better way to sort this?
				if cur[1] > prev[1] {
					line[1] = prev[1]
					line[2] = cur[1]
				}
				vertical = append(vertical, line)
			case RIGHT, LEFT:
				line := [3]int{cur[1], cur[0], prev[0]}
				if cur[0] > prev[0] {
					line[1] = prev[0]
					line[2] = cur[0]
				}
				horizontal = append(horizontal, line)
			}

			// shoelace theorem for area
			// negative area is CCW; positive area is CW
			// x1*y2 - x2*y1
			shoelace_area += prev[0]*cur[1] - cur[0]*prev[1]

			// assume CCW, meaning left-side is inside
			// ! is there a better way to do all of this?
			switch [2]int{prev_d, d} {
			// left, then down, and so on
			case [2]int{LEFT, DOWN}:
				diagonals[prev] = SE
			case [2]int{LEFT, UP}:
				diagonals[prev] = SE | SW | NW
			case [2]int{DOWN, LEFT}:
				diagonals[prev] = SW | SE | NE
			case [2]int{DOWN, RIGHT}:
				diagonals[prev] = NE
			case [2]int{RIGHT, DOWN}:
				diagonals[prev] = NW | NE | SE
			case [2]int{RIGHT, UP}:
				diagonals[prev] = NW
			case [2]int{UP, RIGHT}:
				diagonals[prev] = SW | NW | NE
			case [2]int{UP, LEFT}:
				diagonals[prev] = SW
				// pretty sure I don't care about straight lines
			}
			// discovered everything was 0 when I used `&` accidentally
			// also, discovered I was setting [2]int{0,0} and had 0 for
			// the start tile (needed to iterate start one extra time)
			// fmt.Println("diagonal", prev, diagonals[prev], prev_d, d)
		}

		prev = cur
		prev_d = d
	}

	// if I set an explicit start,
	// this either always happens or never happens
	// as it stands, it *could* happen...
	// just thinking now, if I always started at top-left most point
	// I could *always* start downward, and avoid checking for CW
	is_clockwise := shoelace_area > 0

	// in case of CW, we need to inverse the acceptable diagonals
	inverse_mask := NW + NE + SW + SE

	// iterate diagonals (again?)
	// then they're convex(?) and we should extend the lines
	// to make an extra border
	for k, v := range diagonals {
		if is_clockwise {
			// inverse all the diagonals (all bits on)
			v ^= inverse_mask
			// update
			diagonals[k] = v
		}
		// TOO many switch statements
		// we're making an extra (small) 2 lines around every convex(?) corner
		// they're convex if they have only 1 acceptable diagonal;
		// this switch statement checks for EXACTlY one diagonal
		switch v {
		case NE:
			// draw a horizontal 2-unit length line bottom
			// constant y+1, x-1 -> x+1
			line = [3]int{k[1] + 1, k[0] - 1, k[0] + 1}
			horizontal = append(horizontal, line)
			// draw a vertical 2-unit length line on the left
			// constant x-1, y-1 -> y+1
			line = [3]int{k[0] - 1, k[1] - 1, k[1] + 1}
			vertical = append(vertical, line)
		case SE:
			// draw top&left
			line = [3]int{k[1] - 1, k[0] - 1, k[0] + 1}
			horizontal = append(horizontal, line)

			line = [3]int{k[0] - 1, k[1] - 1, k[1] + 1}
			vertical = append(vertical, line)
		case SW:
			// draw top&right
			line = [3]int{k[1] - 1, k[0] - 1, k[0] + 1}
			horizontal = append(horizontal, line)

			line = [3]int{k[0] + 1, k[1] - 1, k[1] + 1}
			vertical = append(vertical, line)
		case NW:
			// draw bottom&right
			line = [3]int{k[1] + 1, k[0] - 1, k[0] + 1}
			horizontal = append(horizontal, line)

			line = [3]int{k[0] + 1, k[1] - 1, k[1] + 1}
			vertical = append(vertical, line)
		}
	}

	// this was failing because I accidentally did SW & NW & NE instead of `|`
	// fmt.Println("9,5", diagonals[[2]int{9, 5}])

	// so now this perimeter contains all drawn lines, and
	// extra lines around each convex(?) corner to avoid
	// edge cases where we perfectly avoid corner lines
	return &Perimeter{horizontal, vertical, diagonals}
}

func (perimeter *Perimeter) intersectsRect(p1, p2 [2]int) bool {
	// create 4 lines from a and b corners
	// then check if those lines intersect any *perpendicular* perimeter lines

	// top-left and bottom-right
	tl := [2]int{min(p1[0], p2[0]), min(p1[1], p2[1])}
	br := [2]int{max(p1[0], p2[0]), max(p1[1], p2[1])}

	// algorithm is line1.x (p1 and p2) is between p3.x and p4.x (line2 points)
	// and line2.y (p3 and p4) is between p1.y and p2.y

	// check horizontal lines against vertical
	for _, vert := range perimeter.vertical {
		// same check for top and bottom, since their x values are the same
		crosses_x := vert[0] > tl[0] && vert[0] < br[0]
		// top
		if crosses_x && tl[1] > vert[1] && tl[1] < vert[2] {
			return true
		}
		// bottom
		if crosses_x && br[1] > vert[1] && br[1] < vert[2] {
			return true
		}
	}

	// check vertical lines against horizontal
	for _, hor := range perimeter.horizontal {
		crosses_y := hor[0] > tl[1] && hor[0] < br[1]
		// left
		if crosses_y && tl[0] > hor[1] && tl[0] < hor[2] {
			return true
		}
		// right
		if crosses_y && br[0] > hor[1] && br[0] < hor[2] {
			return true
		}
	}

	return false
}

// points form a perimeter
//
// largest area must be within it
func (theater *MovieTheater) LargestAreaWithinPerimeter() (largest int) {
	// need to draw perimeter: (1) track lines,
	perimeter := theater.formPerimeter()
	// ! these map keys are unordered, which *might* cause issues
	tiles := slices.Collect(maps.Keys(theater.redTiles))

	var area int
	var dir [2]bool
	var compare_diagonal DIAGONALS

	for i, a := range tiles {
		for _, b := range tiles[i+1:] {
			if a[0] == b[0] || a[1] == b[1] {
				// just outright ignore straight lines
				continue
			}

			// check if a->b is in acceptable direction (pointed inwards)
			dir = [2]bool{
				b[0] > a[0], // rightward
				b[1] > a[1], // downward
			}

			switch dir {
			case [2]bool{true, true}:
				// SE
				compare_diagonal = SE
			case [2]bool{true, false}:
				// NE
				compare_diagonal = NE
			case [2]bool{false, false}:
				// NW
				compare_diagonal = NW
			case [2]bool{false, true}:
				// SW
				compare_diagonal = SW
			}

			// TODO: could check CW or CCW here instead of inversing diagonals above
			// OR: make sure it is never CW

			if perimeter.tileAndDiagonals[a]&compare_diagonal != compare_diagonal {
				// this area is outside the polygon
				// fmt.Println("Outside:", a, b, perimeter.tileAndDiagonals[a], compare_diagonal)
				continue
			}

			// then (2) check if any intersections with rect (invalid),
			if !perimeter.intersectsRect(a, b) {
				// before (3) getting area
				area = rectArea(a, b)

				largest = max(largest, area)
			}
		}
	}

	return
}
