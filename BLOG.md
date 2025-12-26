### Day 12

**Difficulty: ?/10 ☆☆☆☆☆☆☆☆☆☆**

**Time: 1 hrs**

**Run Time: 500µs**

I had no idea how to do this, and then I looked on Reddit and apparently the test data is just hypothetical, but the real data was incredibly basic:

Each region either can contain all the presents no matter the shape, OR:

Each region cannot contain all the presents no matter the shape.  

All that matters is the region area, the present area, and the cell count for the present shape.

Then check worst case:

```go
// check if more space than worst case (always works)
	totalarea := region.height * region.width
	worstcase := 0
	for i, count := range region.presentCount {
		if count > 0 {
			// all cells (3x3)
			worstcase += presents[i].area * count
		}
	}
	if worstcase <= totalarea {
		// we can always fit
		return true
	}
```

And that's pretty much it.  The rest can all return `false`.

Tried to make an iterator for a day reader, but I couldn't figure out how to do it.
Maybe I need to open the file **within** the iterator.

### Day 11

**Difficulty: 6/10 ★★★★★★☆☆☆☆**

**Time: 2 hrs**

**Run Time: ~1.5s~ ~1.4s~ ~750µs~ 570µs**

DFS and needless goroutines!  Today I used `maps.Clone` against my wishes.

Shaved another 0.1s by removing the `avoid` map idea.

I used the same function for both parts, and later added parameters for `start` and `goal`.

The huge savings was when I started pruning:

```go
// get next states
added := false
for _, next := range rack.servers[state.current] {
	// check if visited
	if _, found := state.visited[next]; found {
		continue
	}

	// check if pruned
	if _, found := pruned[next]; found {
		// fmt.Println("already pruned 2", next)
		continue
	}

	added = true
	// ...
}

if !added {
	// prune!
	pruned[state.current] = struct{}{}
}
```

Without pruning, part 2 never finished.  I had also tried passing paths to ignore as parameters, as the second part is done in segments:

```go
// get all segment path possibilities

// segments
// srv -> [fft -> dac] -> out

// each segment avoids the other parts
segments := []struct {
	start, goal string
}{
	{"svr", "fft"},
	{"fft", "dac"},
	// {"dac", "fft"}, // ! this returns 0
	{"dac", "out"},
}
product := 1

for _, v := range segments {
	product *= rack.NumPathsConnecting(v.start, v.goal)
}

return product
```

Shaved a lot of time off, by just using a simple recursive function, and a shared cache, instead:

```go
func (rack *ServerRack) RecursivePathsConnecting(start, goal string) (count int) {
	if start == goal {
		return 1
	}

	// if cached
	if val, found := rack.cache.Get(start, goal); found {
		return val
	}

	sum := 0
	for _, next := range rack.servers[start] {
		sum += rack.RecursivePathsConnecting(next, goal)
	}

	// cache
	rack.cache.Set(start, goal, sum)

	return sum
}
```

And the cache was just a struct with a `map`, with some custom getters and setters:

```go
type Cache struct {
	// from -> to -> count
	data map[string]map[string]int
}

func NewCache() Cache {
	return Cache{
		// need to initialize this
		data: map[string]map[string]int{},
	}
}

func (cache *Cache) Get(start, goal string) (int, bool) {
	if a, found := cache.data[start]; found {
		if b, found := a[goal]; found {
			return b, true
		}
	}
	return 0, false
}

func (cache *Cache) Set(start, goal string, val int) {
	if a, found := cache.data[start]; found {
		a[goal] = val
	} else {
		cache.data[start] = map[string]int{
			goal: val,
		}
	}
}
```

I tried a `sync.Mutex` on the cache (so that I could continue to abuse goroutines), but it *cost* me 100µs on average.  Also, TIL that `Option+M` creates `µ` on a Mac.

### Day 10

**Difficulty: 10/10 ★★★★★★★★★★**

**Time: 20+ hrs**

**Run Time: ~235ms**

Almost 700 LOC, and multiple attempts.

Part one was basically summarized in the `allPossibleButtonPresses` function, which I eventually turned into an iterator, and tried to use for Part 2.  The function is meant only to return all button combinations that would produce some diagram: for part 1 that's the light diagram.

My machine (for part 1 anyway) uses a lot of binary:

```go

type Machine struct {
	light_diagram int         // (reversed) binary
	button_wiring map[int]int // index binary -> button binary
	// ...
}

// for light diagram:
// do this in reverse,
// so the button ints line up properly
for j := len(inside) - 1; j >= 0; j-- {
	machine.light_diagram <<= 1
	if inside[j] == '#' {
		machine.light_diagram += 1
	}
}

// for buttons:
button := 0
for val := range strings.SplitSeq(inside, ",") {
	// 3 is not 3, but instead the 3rd bit
	num := utils.ParseInt(val)
	button |= 1 << num
}
// i will start at 2 here, because buttons come second
index_binary := 1 << (index - 2)
machine.button_wiring[index_binary] = button
```

So now I've got button wiring that translates:

`(3) (1,3) (2)`

to (binary):

`1: 1000, 10: 1010, 100: 0100`

OR (decimal):

`1: 8, 2: 10, 4: 4`

And the diagrams are translated from:

`[...#.], [.###.#]` 

to (reversed):

`8, 46`

Anyway, I thought this was smart, given in Part 1 I leveraged the binary to toggle buttons (via `^` XOR), check if keys were already pressed (via `&` AND), and decide whether a button is worth pressing (via `|` OR).

```go
item := utils.MustShift(&queue)

current, pressed := item[0], item[1]

// get next buttons to press
for k, v := range machine.button_wiring {
	if k&pressed == k {
		// already pressed
		continue
	}
	// only press if it will impact either current or the diagram
	// seems to save about ~10ms
	compare := current | diagram

	if v&compare != 0 {
		// queue next state
		// toggles lights
		next := current ^ v
		next_pressed := pressed | k
```

This is all fine.  Part 2 I had no idea how to do, but I saw that a bunch of people on Reddit used third party libraries.  I saw someone suggest generating a binary for odd digits in Joltage, pressing those buttons, then dividing joltage by 2, and so on, until joltage is 0's; I could get that to work on the test data, but I couldn't get it to work on my real data.  I also saw people talking about Gaussian Elimination.  

I read [the wiki page on it](https://en.wikipedia.org/wiki/Gaussian_elimination), and [watched a video solution using it](https://www.youtube.com/watch?v=xibCHVRF6oI), and eventually thought I could try it.  I [understood most of the rust code, so was able to use some of that too](https://gist.github.com/icub3d/16eea2a8b4a94d193a148fef908779a9).

First I needed to figure out how to get a reduced row echelon form matrix, which I gleaned from wikipedia, and then asked ChatGPT to convert from pseudocode. I understand some of it, but didn't think too hard about it.  I understand that there are properties of a matrix that we can exploit to help reduce the scope of this problem.  I would have to practice a bit more on matrices to fully understand how to simplify this code, or why it works.

Then I created a small script to convert the joltages to matrices. The rest is mostly lifted from icub3d's implementation:

```go
func (machine *Machine) GaussianElimination() (fewest int) {
	matrix := machine.ToMatrix()

	reduced := RREF(matrix)

	// never going to be higher than the highest joltage
	max := slices.Max(machine.joltage)
	min := math.MaxInt
	values := make([]int, len(reduced.Free))

	// DFS
	dfs(reduced, 0, values, &min, max)

	return min
}
```

This is nice and brief.  First time using `slices.Max` to return some max value from a slice.  I think this is the first time I've passed a reference to an `int` so that some recursive function can just update it, but it's very clever: easier than keeping track of return values. I think in previous years I passed a `channel` and waited for values, though I'm not sure how that would help me here.

So yeah, learned a lot about matrices: maybe some better practices for dfs recursive functions.



### Day 9

**Difficulty: 8/10 ★★★★★★★★☆☆**

**Time: 8 hrs**

**Run Time: 30ms**

Nearly a full day on this problem. And ~400 LOC. And ~5 pages in my notebook.

First, I added a new type for my Day Reader:

```go
day.reader = any(ReadCSVInt).(func(string) T)
```

Then I was able to just use `[]int` as my data type, and parse with `slices.Chunk`:

```go
// first time using Chunk
for chunk := range slices.Chunk(data, 2) {
	redTiles[[2]int{chunk[0], chunk[1]}] = struct{}{}
}
```

In order to do this I updated the `ScanCommas` scanner function to check for commas first, then check for newlines:

```go
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
```

First time using `min`, and `max` today, which were apparently added in 1.21 (after my last AOC in go).

So, for **part one**, I created a slice from the map keys, using `slices.Collect`, and `maps.Keys`, I think both for the first time:

```go
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
```

This returned the largest area within ~100µs.

**Part TWO**, I thought about doing a flood fill on all lines, and all inner cells, but worried it would be the wrong solution.  I thought I could check intersections of lines to invalidate areas, though that became challenging.  Thus ~400 LOC.

First, perpendicular axis-aligned line intersections are pretty basic:

Given two lines:

```console
	A
	|
	|  C---D
	|
	B
```

It intersects if `A.x` is between `C.x` and `D.x` AND if `C.y` is between `A.y` and `B.y`.  I just had to do that twice with a rect's top, bottom, left, and right sides (against *all* stored vertical and horizontal lines).

The collection of those lines were rather simple: I started anywhere (because my tiles are in an unordered `map`), I iterate in 4 directions, and increase distance each loop, in order to check for nearest adjacent.

Something like:

```go
var DIRS = [4][2]int{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
}
// ...

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
```

The rest is even more painful.  After doing that, I realized it was still passing for points like `A` and `B` here:

```console
         A...o
outside  .
         .
 B.......o
 .
 .    inside
 o
```

Because I was only checking intersections **within** the polygon.  Now I had to figure out how to prevent checks that are **outside** the polygon.  I also realized that the way I was creating the polygon, I had no idea at any given line, which side was inside, and which was outside.  I could *guess* that I was going CCW and the left was inside, but I'd be wrong half the time, due to my (probably foolish) implementation of grabbing the first key from an unordered map as starting point.

In any case, I was able to use shoelace theorem (first used in [2023 Day 10](https://github.com/bozdoz/advent-of-code-2023/blob/75966a2785074241aad99c8ab42b77f670d77de2/day-10/src/main.rs#L173)) to determine if I was going CCW or CW (negative or positive areas).

I also compared each line's direction with the previous line's direction, in order to determine if the corner was convex (<=180) or concave (>=180) (still not sure which is which).  So I tracked diagonals which point inward for every corner: some have 3 diagonals, and some have 1.

```console
         A...o
outside  .↘️
         .↗️
 B.......o
 .↘️    ↙️ ↘️
 .    inside
 o
```

This was painful:

```go
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
```

And then it didn't work.  Actually, it passed the tests half of the time.  Meaning there was some problem with the directions.  First, I initially used `&` instead of `|` to create the diagonals.  Second, I ran into this comparison (A->B):

```console
         o...o
         .   .
         .   .
 B.......o   .
 .           .
 .           .
 o.......C   .
         .   .
         D...A
```

Here a rectangle between A->B includes the empty space bottom-left. It doesn't cross any of the perimeter lines that I set up.  For example, `D.y` is **equal to** `A.y` (not between `A.y` and `B.y`).

So I had to come up with something else: extra lines on all convex(?) corners:

```console
        ┌-   -┐
        |o...o|
         .   .
┌-       .   .
|o.......o   .
 .           .
 .           .
|o.......o   .
∟-       .   .
				|o...o|
        ∟-   -┘
```

Now I had lines between all tiles, and lines around each outer corner, which should at least catch the edge case above.

```go
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
```

Now I did all of this, because *I didn't know how to grow the polygon outwards* to begin with: not that that would definitely be easier.

Anyway, this worked.  It's under a second.  I'm happy with it.  The actual checks before the area looks like this (still long):

```go
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
```

### Day 8

**Difficulty: 7/10 ★★★★★★★☆☆☆**

**Time: 3 hrs**

**Run Time: 100ms**

TIL that `go test` != `go test .`, and I don't know why.  It looks like `go test` preserves the stdout, and `go test` squelches it so long as tests pass.  This seems more useful to me for debugging.

I had a lot going on today with node id's, which didn't turn out terribly.  These are all the same node id's as `int`:

```go
a, b     int // index

adjacent []int

nodes    map[int]*Node
```

Parsing was pretty easy.  I copied some of my 3d distance function from 2021:

```go
func (p1 *Point) Distance(p2 Point) float64 {
	x := p2.x - p1.x
	y := p2.y - p1.y
	z := p2.z - p1.z

	return math.Sqrt(float64(x*x + y*y + z*z))
}
```

I'm surprised that the scripts run as fast as they do.  I expected all the looping and sorting I was doing would be costly.

Creating all the nodes and edges was straight forward:

```go
for i, a := range junction.points {
	for j := i + 1; j < len(junction.points); j++ {
		dist := a.Distance(junction.points[j])
```

`j` starts at `i+1` to make sure each node sees every other node once.

Then it's sorted, ascending.

Then I connect all the adjacent nodes in order of distance, which is very short, and very ugly:

```go
func (junction *Junction) ConnectNodes(count int) {
	for _, pair := range junction.distances[:count] {
		// ugly
		junction.nodes[pair.a].adjacent = append(junction.nodes[pair.a].adjacent, pair.b)
		junction.nodes[pair.b].adjacent = append(junction.nodes[pair.b].adjacent, pair.a)
	}
}
```

Then I did a silly thing, and I made another `Sequence` to iterate.

So I iterated the distances again, that I sorted, which gave me edges, then I did a `range` to do a BFS over connected nodes.

```go
// maybe not the best time for me to be throwing myself
// into another sequence
func (junction *Junction) eachConnected(start int, visited *map[int]struct{}) iter.Seq[int] {
	queue := append([]int{}, junction.nodes[start].adjacent...)

	return func(yield func(int) bool) {
		for len(queue) > 0 {
			connectedId := shift(&queue)

			if _, found := (*visited)[connectedId]; !found {
				// add to visited
				(*visited)[connectedId] = struct{}{}
				// add next adjacent
				queue = append(queue, junction.nodes[connectedId].adjacent...)
				// check if statement, though I'm not sure we'll ever break
				if !yield(connectedId) {
					return
				}
			}
		}
	}
}
```

This shared a visited map with the caller, since we'd never run into a visited cell within the loop that was part of an alternate circuit.  

Then it counts and sorts the counts to get the top 3.

For Part 2, I just keep connecting until all are visited.  Realizing now I just have to iterate them until they're all visited.


### Day 7

**Difficulty: 5/10 ★★★★★☆☆☆☆☆**

**Time: 2 hrs**

**Run Time: 340µs**

This was difficult to keep in my head, due to trying out DFS and BFS and also trying to do some kind of path traversing, while ignoring visited nodes.

I put to pen and paper an algorithm for part 2, which worked just fine for the test data, but not for the actual data.  So I have no idea where I went wrong.  What I tried to do worked also for Part 1:

1. I did a BFS, tracking all the seen Splitters
2. Each state was either the start or a splitter, and tracked how many beams hit it, and whether any reach the end

```go
type Splitter struct {
	beams         int // usually splits 1
	beamsReachEnd int // out of 2, how many beams
}
// ... (in loop)
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
// ...
```

3. So it carried its parents beams, without having to continue to produce next states for the loop
4. I iterated the returned seen Splitter map, and multiplied number of beams by number of beams that reached the end

```go
//
// no idea why this didn't work on my input, but worked on test data
//	
splitters := grid.Travel()

sum := 0

for _, splitter := range splitters {
	// we only count the beams that exit the grid
	if splitter.beamsReachEnd > 0 {
		sum += splitter.beams * splitter.beamsReachEnd
	}
}
```

Again, this worked just fine for test data.

I asked copilot to figure out some edge case I'm missing, and it produced a bunch of garbage.

Then I saw [a visualization on Reddit](https://www.reddit.com/r/adventofcode/comments/1pgnmou/2025_day_7_lets_visualize/) which I knew would work, so I made some code for that.  This code worked in under 30µs, which is incredible.  Though I have no idea how to use it to compare and figure out where I went wrong with my algorithm.  For my real input, the two methods differed by hundreds of thousands: the test data, identical 🫠.

I'm really liking 2d arrays for grids now.  Up and down is just subtracting or adding `grid.width`, and conversion to columns is just `index%grid.width`; also, checking out-of-grid is a combination of checking the column and whether index is < 0 or > height*width.  I did make some debugging functions while I worked though:

```go
func (grid *Grid) DebugToRC(i int) [2]int {
	r := i / grid.width
	c := i % grid.width

	return [2]int{r, c}
}

func (grid *Grid) DebugFromRC(rc [2]int) int {
	return rc[0]*grid.width + rc[1]
}
```

### Day 6

**Difficulty: 1/10 ★☆☆☆☆☆☆☆☆☆**

**Time: 30 min**

**Run Time: 260µs**

Another easy one.  This was just a data-parsing puzzle.  Nothing too advanced.

I made two structs, one I regretted immediately.  I made some constants with a custom rune-type, which I later regretted.

I heavily relied on `strings.Fields` which breaks a line up by space-separated values (I wasn't exactly sure if multiple spaces would work, but it did).

So Part 1 was a breeze.

Part 2 was just thinking of the input as a grid (thankfully there were empty spaces in the data):

```console
123 328  51 64.
 45 64  387 23. 
  6 98  215 314
*   +   *   +  
```

Otherwise, I would have had a harder time analyzing that right-most column.

But yeah, this grid-cell traversing was simple:

```go
func NewGridCollection(data []string) *Collection {
	height := len(data)
	width := len(data[0])

	// start top-down, right to left
	for c := width - 1; c >= 0; c-- {
		num := 0
		for r := range height {
			cell := data[r][c]
```

I think I tweaked the logic once or twice and got the right answer.

The annoying part of making a custom rune-type was checks like this where I couldn't even use it: 

```go
// operation row
if cell == '+' || cell == '*' {
	// problem is finished
	problem.operation = Operation(cell)
```

Instead I had to check the rune, and then convert to `Operation`, which seems completely needless.

I copied over a part of the `ParseInt` function, because I knew the cell was going to contain a single digit:

```go
if cell != ' ' {
	// is number (simplified from utils.ParseInt)
	num = num*10 + int(cell-'0')
}
```

Which removes negative checks, and iterating a full string.  Maybe I should save that to a utility function, but I guess it's less characters already than any function name I could give it: `int(cell-'0')`

### Day 5

**Difficulty: 1/10 ★☆☆☆☆☆☆☆☆☆**

**Time: 30 min**

**Run Time: 550µs**

Using `fmt.Sscanf` to parse ints in this puzzle; I had first accidentally tried `fmt.Scanf`, which requires a reader.

Used `sort.Slice` to sort the ranges, before determining if they overlap (or are enveloped).

Super easy day.  I did add one extra test to verify that sort.Slice worked, but that was all.

#### Update

I changed the sorter from `sort.Slice` to the *new* Sequence-driven `slices.SortFunc`, which also uses `cmp.Compare`, and it is 4x quicker:

```console
cpu: Apple M4
BenchmarkSortRange-10    26835813     43.86 ns/op
BenchmarkSortSlices-10   149276900     7.891 ns/op
```


### Day 4

**Difficulty: 1/10 ★☆☆☆☆☆☆☆☆☆**

**Time: 30 min**

~~**Run Time: 6.7ms**~~
**Run Time: 3.3ms**

#### Update

I revisited the algorithm and realized I could use one of Go's *new* iterator sequences: https://pkg.go.dev/iter (since go 1.23)

So I took the old code (dependency injection): 

```go
grid.forEachNeighbour(cell, func(char rune) {
	if char == PAPER {
		count++
	}
})
```

And did this instead:

```go
for char := range grid.eachNeighbour(cell) {
	if char == PAPER {
		count++
		if count == 4 {
			// don't keep counting the rest of the 8 neighbours
			break
		}
	}
}
```

Though the function is a bit more verbose with all the handling of `yield` results (checking for `break` I'm assuming):

```go
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
		// ...
```

But this halved my overall time, which is great!

#### Original

Didn't feel like I did anything today: not many tests, nothing new with utility scripts. Just the first grid puzzle of the year.

I anticipated many pitfalls from previous years and created an ERROR cell:

```go
const (
	EMPTY = '.'
	PAPER = '@'
	ERROR = '!' // my own definition
)
```

And a getter (I ended up just needing the index, instead of querying by row and column):

```go
func (grid *Grid) getByIndex(index int) rune {
	if index < 0 || index >= len(grid.cells) {
		return ERROR
	}
	return grid.cells[index]
}
```

That way I can ignore the edges of the grid.  The other issues was anticipating dealing with a 2d array for the grid:

```go
type Grid struct {
	height, width int
	cells         []rune
}

func (grid *Grid) forEachNeighbour(i int, fun func(char rune)) {
	// ...
	col := i % grid.width
	// if column is 0 then we can't go left!
	if col != 0 {
		fun(grid.getByIndex(top - 1))
		fun(grid.getByIndex(i - 1))
		fun(grid.getByIndex(bottom - 1))
	}
	// ...
}
```

I'm actually not sure if 2d array is faster than `map[int]map[int]rune`, but from memory of previous years it felt like it would be faster.

The only other gotcha was remembering to remove the removable `PAPER` *after* the counting was done (so we didn't alter the counts).  But this was an obvious gotcha; one that `rust` would have never let me get away with in the first place 🦀

### Day 3

**Difficulty: 1/10 ★☆☆☆☆☆☆☆☆☆**

**Time: 1 hrs**

**Run Time: 345µs**

This seemed very straight forward with a couple interesting challenges.  The first part was very simple to just iterate once and compare first and second max values.  The second part initially made me start backwards, with all initial values on the right-most 12 digits.  Then I iterated backwards, moving the head, then the tail.  I realized quickly after that I should just always get the left-most highest number, and it made for a much simpler function:

```go
func (jolt *Joltage) LargestOptimized() int {
	start := len(jolt.batteries) - 12
	// start with right-most 12 digits as largest
	indices := make([]int, 12)

	for i := range 12 {
		indices[i] = start + i
	}

  // don't need the *copied* value, because we always want the current value
	for i := range indices {
		// don't iterate to index -1
		until := -1
		if i > 0 {
			// don't go beyond previous node
			until = indices[i-1]
		}
		for j := indices[i] - 1; j > until; j-- {
			if jolt.batteries[j] >= jolt.batteries[indices[i]] {
				indices[i] = j
			}
		}
	}

	out := 0

	for i := range 12 {
		out *= 10
		out += jolt.batteries[indices[i]]
	}

	return out
}
```

The only awkward part was keeping track of `jolt.batteries[indices[i]]`.  I benchmarked the two approaches, and found the "optimized" one ran 5x as fast:

```console
BenchmarkLargestChain/987654321111111-10          4313073     261.5 ns/op
BenchmarkLargestChain/811111111111119-10          4479896     268.5 ns/op
BenchmarkLargestOptimized/987654321111111-10     22342346      52.80 ns/op
BenchmarkLargestOptimized/811111111111119-10     21657750      54.55 ns/op
```

OK, I updated it again to reduce the number of iterations to 1 instead of 3, and it halved the time:

```console
BenchmarkLargestChain/987654321111111-10          4415940     259.6 ns/op
BenchmarkLargestChain/811111111111119-10          4405172     272.1 ns/op
BenchmarkLargestOptimized/987654321111111-10     50055966      24.58 ns/op
BenchmarkLargestOptimized/811111111111119-10     49564675      24.12 ns/op
```

### Day 2

**Difficulty: 2/10 ★★☆☆☆☆☆☆☆☆**

**Time: 2 hrs**

~~**Run Time: 100ms**~~
**Run Time: 53ms**

Added a `WithReader` method for the Day utility:

```go
var day = utils.NewDay[d]().WithReader(utils.ReadCSV)

// adds some custom reader
func (day *day[T]) WithReader(fun func(string) T) *day[T] {
	day.reader = fun
	return day
}
```

I copied over `ScanLines` and converted it into a `ScanCommas` function for the bufio scanner.

Today I thought I could try to figure out how to compare the segments of a number using math, instead of string comparisons.

The idea was to get the count of digits in the number, get segments of that number using math.Pow10:

```go
val := 123456
pow := int(math.Pow10(3))
compare := val / pow 
// 123 which can then be compared with the next segments
```

After getting that chunk, I would get the right-most digits for comparing, via modulus of a similar math.Pow10 value.

All the tests passed today, without the input data working.  Had to try regex, until I found out that Go doesn't support back references: `\1`.

Then I just converted the value to a string and chunked it to compare strings.  This gave me the correct answer, and was way easier to reason about.  I used that to compare with the previous function, in order to find which inputs I was failing with.

Surprisingly, the string comparison function ran faster with the real input, but benchmarked roughly similar, if not sometimes slower:

```console
goos: darwin
goarch: arm64
pkg: github.com/bozdoz/advent-of-code-2025/02
cpu: Apple M4
BenchmarkMath/111-10             39575978     30.04 ns/op
BenchmarkMath/1020102-10         18349874     64.88 ns/op
BenchmarkMath/123123123123-10    21972421     53.57 ns/op
BenchmarkString/111-10           37649468     31.80 ns/op
BenchmarkString/1020102-10       32080914     36.64 ns/op
BenchmarkString/123123123123-10  24248792     48.75 ns/op
```

### Day 1

**Difficulty: 1/10 ★☆☆☆☆☆☆☆☆☆**

**Time: 30 min**

**Run Time: 260µs**

First time using this new boilerplate.  

```go
type d = []string

var day = utils.NewDay[d]()

func main() {
	day.Read("./input.txt")

	day.Run(partOne)
	day.Run(partTwo)
}
```

It's supposed to figure out how to read today's puzzle based on the type:

```go
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
	default:
		panic(fmt.Sprintf("need to implement a reader for type: %T", day.__type__))
	}

	return day
}
```

Surprisingly needed to cast as `any` a couple times there, but it looks like it will work quite well.

This new setup also makes the testing just as easy:

```go
var data = day.Read("./example.txt")
```