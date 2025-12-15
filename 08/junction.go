package main

import (
	"fmt"
	"iter"
	"math"
	"slices"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Point struct {
	x, y, z int
}

type Pair struct {
	a, b int // index
	dist float64
}

type Node struct {
	adjacent []int
}

type Junction struct {
	points    []Point
	distances []Pair
	nodes     map[int]*Node
}

func NewJunction(data []string) *Junction {
	points := make([]Point, len(data))
	nodes := make(map[int]*Node)

	for i, line := range data {
		var x, y, z int

		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)

		points[i] = Point{x, y, z}
		nodes[i] = &Node{}
	}

	return &Junction{points, []Pair{}, nodes}
}

// from https://github.com/bozdoz/advent-of-code-2021/blob/main/types/vector3d.go
func (p1 *Point) Distance(p2 Point) float64 {
	x := p2.x - p1.x
	y := p2.y - p1.y
	z := p2.z - p1.z

	return math.Sqrt(float64(x*x + y*y + z*z))
}

func (p1 *Point) Equals(p2 Point) bool {
	return p1.x == p2.x && p1.y == p2.y && p1.z == p2.z
}

// get distances and sort
func (junction *Junction) GetDistances() {
	for i, a := range junction.points {
		for j := i + 1; j < len(junction.points); j++ {
			dist := a.Distance(junction.points[j])
			junction.distances = append(junction.distances, Pair{
				a:    i,
				b:    j,
				dist: dist,
			})
		}
	}

	// sort
	slices.SortFunc(junction.distances, func(x, y Pair) int {
		if x.dist < y.dist {
			return -1
		}
		if x.dist > y.dist {
			return 1
		}
		return 0
	})
}

// join nearest nodes
func (junction *Junction) ConnectNodes(count int) {
	for _, pair := range junction.distances[:count] {
		// ugly
		junction.nodes[pair.a].adjacent = append(junction.nodes[pair.a].adjacent, pair.b)
		junction.nodes[pair.b].adjacent = append(junction.nodes[pair.b].adjacent, pair.a)
	}
}

// maybe not the best time for me to be throwing myself
// into another sequence
func (junction *Junction) eachConnected(start int, visited *map[int]struct{}) iter.Seq[int] {
	queue := append([]int{}, junction.nodes[start].adjacent...)

	return func(yield func(int) bool) {
		for len(queue) > 0 {
			connectedId := utils.MustShift(&queue)

			if _, found := (*visited)[connectedId]; !found {
				// only unseen ids
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

func (junction *Junction) GetTopThreeCircuits(count int) []int {
	// iterate same distances
	// check if visited nodes
	visited := map[int]struct{}{}
	counts := []int{}
	// return sizes of circuits
	for _, pair := range junction.distances[:count] {
		// check if visited before iterating
		if _, found := visited[pair.a]; found {
			continue
		}
		visited[pair.a] = struct{}{}
		count := 1

		for range junction.eachConnected(pair.a, &visited) {
			// ! don't check visited here
			count++
		}

		counts = append(counts, count)
	}

	slices.SortFunc(counts, func(x, y int) int {
		// return 1 to make this descending
		if x < y {
			return 1
		}
		if x > y {
			return -1
		}
		return 0
	})

	return counts[:3]
}

func (junction *Junction) ConnectUntilAllJoined() (Point, Point) {
	// once every node is visited, we're done
	visited := map[int]struct{}{}
	// now we start connecting and checking
	want := len(junction.points)

	var lastPair Pair

	for i := 0; len(visited) != want; i++ {
		pair := junction.distances[i]

		visited[pair.a] = struct{}{}
		visited[pair.b] = struct{}{}

		lastPair = pair
	}

	return junction.points[lastPair.a], junction.points[lastPair.b]
}
