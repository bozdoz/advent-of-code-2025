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