package main

import (
	"log"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Joltage struct {
	batteries []int
}

func NewJoltage(data string) *Joltage {
	batteries := make([]int, len(data))

	for i, char := range data {
		batteries[i] = utils.ParseInt(string(char))
	}

	return &Joltage{
		batteries,
	}
}

func (jolt *Joltage) Largest() int {
	first := jolt.batteries[0]
	second := jolt.batteries[1]

	// first considers 0 through len-2
	for i := 1; i < len(jolt.batteries)-1; i++ {
		n := jolt.batteries[i]

		if n > first {
			first = n
			second = jolt.batteries[i+1]
		} else if n > second {
			second = n
		}
	}

	// second considers the very last digit also
	if jolt.batteries[len(jolt.batteries)-1] > second {
		second = jolt.batteries[len(jolt.batteries)-1]
	}

	return first*10 + second
}

func (jolt *Joltage) LargestChain() int {
	start := len(jolt.batteries) - 12
	// start with right-most 12 digits as largest
	indices := make([]int, 12)

	for i := range 12 {
		indices[i] = start + i
	}

	// iterate backwards
	for i := start - 1; i >= 0; i-- {
		num := jolt.batteries[i]

		// we look for a new head, and only then re-arrange the tail
		if num >= jolt.batteries[indices[0]] {
			// update the index
			indices[0] = i
			// start with next index
			compare := 1

			// we are checking here
			//    v
			// 234234 batteries
			//   0 12 indices
			//     ^
			j := i + 1
			cur := indices[compare]

			did_update := false
			// re-arrange the rest
			// find the left-most highest value
			for {
				// still look backwards (because we need left-most)
				cur--
				// we need to move to the highest within a range
				// x  y
				// 8181
				// y should move to 8, not 1; but if it were:
				// x  y
				// 8881
				// then y should move to the left-most 8

				if jolt.batteries[cur] >= jolt.batteries[indices[compare]] {
					// update
					did_update = true
					indices[compare] = cur
				}

				if cur == j {
					if did_update && compare < 11 {
						log.Println("update", indices[compare])
						// continue down the tail
						did_update = false
						j = indices[compare] + 1
						compare++

						cur = indices[compare]
						log.Println("j", j, "cur:", cur)
					} else {
						break
					}
				}
			}
		}
	}

	// recombine the indices
	log.Println(indices)

	out := 0

	for i := range 12 {
		out *= 10
		out += jolt.batteries[indices[i]]
	}

	return out
}

func (jolt *Joltage) LargestOptimized(length int) int {
	start := len(jolt.batteries) - length
	// start with right-most 12 digits as largest
	indices := make([]int, length)

	for i := range length {
		indices[i] = start + i
	}

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

	// recombine the indices
	log.Println(indices)

	out := 0

	for i := range length {
		out *= 10
		out += jolt.batteries[indices[i]]
	}

	return out
}
