package main

import (
	"maps"
	"slices"
)

type MovieTheater struct {
	redTiles map[[2]int]struct{} // x,y

}

func NewMovieTheater(data []int) *MovieTheater {
	// first time using Chunk and Collect
	redTiles := make(map[[2]int]struct{}, len(data))

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
	for i, a := range tiles {
		for _, b := range tiles[i+1:] {
			area := rectArea(a, b)

			// first time using max
			largest = max(largest, area)
		}
	}

	return
}

type Perimeter struct {
	horizontal, vertical [][3]int // h = y,x1,x2; v = x,y1,y2
}

const (
	UP = iota
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

func (theater *MovieTheater) formPerimeter() *Perimeter {
	// start from any point, and move in each direction to draw the perimeter lines
	var start [2]int
	for k := range theater.redTiles {
		start = k
		break
	}
	cur := start
	ignore_dir := -1 // first check every direction, then ignore prev
	dist := 0
	horizontal := make([][3]int, 0, len(theater.redTiles))
	vertical := make([][3]int, 0, len(theater.redTiles))

outer:
	for {
		dist++
		for d, dir := range DIRS {
			if d == ignore_dir {
				continue
			}
			// move current by (dist) amount in (dir) direction
			check := [2]int{cur[0] + dir[0]*dist, cur[1] + dir[1]*dist}

			// check if check is a red tile
			if _, found := theater.redTiles[check]; found {
				// draw line
				switch d {
				case UP, DOWN:
					// first int is the shared val
					// the rest are min/max
					rest := [2]int{cur[1], check[1]}
					if cur[1] > check[1] {
						rest = [2]int{check[1], cur[1]}
					}
					vertical = append(vertical, [3]int{cur[0], rest[0], rest[1]})
				case RIGHT, LEFT:
					rest := [2]int{cur[0], check[0]}
					if cur[0] > check[0] {
						rest = [2]int{check[0], cur[0]}
					}
					horizontal = append(horizontal, [3]int{cur[1], rest[0], rest[1]})
				}

				if check == start {
					// we're done
					break outer
				}
				// get ready for next loop
				cur = check
				// ignore the opposite of this direction
				ignore_dir = (d + 2) % 4
				// reset dist
				dist = 0
				break
			}
		}
	}

	return &Perimeter{horizontal, vertical}
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
	tiles := slices.Collect(maps.Keys(theater.redTiles))
	for i, a := range tiles {
		for _, b := range tiles[i+1:] {
			// then (2) check if any intersections with rect (invalid),
			if !perimeter.intersectsRect(a, b) {
				// before (3) getting area
				area := rectArea(a, b)

				// fmt.Println(a, b, area)

				largest = max(largest, area)
			}
		}
	}

	return
}
