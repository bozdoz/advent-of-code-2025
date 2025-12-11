package main

import (
	"maps"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type ServerRack struct {
	servers map[string][]string
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

	return &ServerRack{servers}
}

type State struct {
	current string
	visited map[string]struct{}
}

func (rack *ServerRack) NumPathsConnecting(start, goal string, avoid map[string]struct{}) (count int) {
	// start at "you", end at "out"?
	queue := []State{
		{
			current: start,
			visited: map[string]struct{}{
				start: {},
			},
		},
	}

	pruned := map[string]struct{}{}

	for len(queue) > 0 {
		// DFS
		state := utils.MustPop(&queue)

		// check if done
		if state.current == goal {
			count++
			continue
		}

		// get next states
		added := false
		for _, next := range rack.servers[state.current] {
			// check if avoided
			if _, found := avoid[next]; found {
				continue
			}

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

			queue = append(queue, State{
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
