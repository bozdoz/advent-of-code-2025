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