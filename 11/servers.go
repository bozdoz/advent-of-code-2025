package main

import (
	"maps"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Cache struct {
	// from -> to -> count
	data map[string]map[string]int
	// I needlessly want to run goroutines
	// mu sync.Mutex
}

func NewCache() Cache {
	return Cache{
		data: map[string]map[string]int{},
	}
}

func (cache *Cache) Get(start, goal string) (int, bool) {
	// for goroutines (takes 200µs longer)
	// cache.mu.Lock()
	// defer cache.mu.Unlock()
	if a, found := cache.data[start]; found {
		if b, found := a[goal]; found {
			return b, true
		}
	}
	return 0, false
}

func (cache *Cache) Set(start, goal string, val int) {
	// cache.mu.Lock()
	// defer cache.mu.Unlock()
	if a, found := cache.data[start]; found {
		a[goal] = val
	} else {
		cache.data[start] = map[string]int{
			goal: val,
		}
	}
}

type ServerRack struct {
	servers map[string][]string
	cache   Cache
}

func NewServerRack(data []string) *ServerRack {
	servers := make(map[string][]string, len(data))

	for _, line := range data {
		var key string
		for server := range strings.FieldsSeq(line) {
			if strings.HasSuffix(server, ":") {
				// this is the key
				key = strings.TrimSuffix(server, ":")
				// initialize slice
				servers[key] = []string{}
				continue
			}

			// add to key
			servers[key] = append(servers[key], server)
		}
	}

	return &ServerRack{servers, NewCache()}
}

type State struct {
	current string
	visited map[string]struct{}
}

// not used: slower than recursive below
func (rack *ServerRack) NumPathsConnecting(start, goal string) (count int) {
	// start at "you", end at "out"?
	stack := []State{
		{
			current: start,
			visited: map[string]struct{}{
				start: {},
			},
		},
	}

	pruned := map[string]struct{}{}

	for len(stack) > 0 {
		// DFS
		state := utils.MustPop(&stack)

		// check if done
		if state.current == goal {
			count++
			continue
		}

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

			// need to copy visited (expensive?)
			// first time using maps.Clone
			next_visited := maps.Clone(state.visited)

			// add next
			next_visited[next] = struct{}{}

			stack = append(stack, State{
				next,
				next_visited,
			})
		}

		if !added {
			// prune!
			// fmt.Println("Pruned", state.current)
			pruned[state.current] = struct{}{}
		}
	}

	return
}

// much faster than above
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
