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