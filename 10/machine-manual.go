package main

import (
	"fmt"
	"iter"
	"log"
	"math"
	"math/bits"
	"slices"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type Machine struct {
	light_diagram int           // (reversed) binary
	button_wiring map[int]int   // index binary -> button binary
	plain_buttons [][]int       // joltage_indexes
	joltage       []int         // part 2 button press counts
	cache         map[int][]int // another one
}

type Manual struct {
	machines []Machine
}

func NewManual(data d) *Manual {
	machines := []Machine{}

	for line := range data {
		// need to update the actual machine
		machine := Machine{}
		// need to initialize maps
		machine.button_wiring = map[int]int{}
		// inner index (for buttons)
		index := 0

		for part := range strings.FieldsSeq(line) {
			index++
			inside := part[1 : len(part)-1]
			switch {
			case strings.HasPrefix(part, "["):
				// do this in reverse,
				// so the button ints line up properly
				for j := len(inside) - 1; j >= 0; j-- {
					machine.light_diagram <<= 1
					if inside[j] == '#' {
						machine.light_diagram += 1
					}
				}
			case strings.HasPrefix(part, "("):
				button := 0
				plain := []int{}
				for val := range strings.SplitSeq(inside, ",") {
					// 3 is not 3, but instead the 3rd bit
					num := utils.ParseInt(val)
					button |= 1 << num
					// this one 3 is just plain 3
					plain = append(plain, num)
				}
				machine.plain_buttons = append(machine.plain_buttons, plain)
				// i will start at 2! here, because buttons come second
				index_binary := 1 << (index - 2)
				machine.button_wiring[index_binary] = button
			case strings.HasPrefix(part, "{"):
				for val := range strings.SplitSeq(inside, ",") {
					machine.joltage = append(machine.joltage, utils.ParseInt(val))
				}
			}
		}

		machine.cache = make(map[int][]int)

		machines = append(machines, machine)
	}

	return &Manual{machines}
}

func (machine *Machine) FewestButtonPresses() (fewest int) {
	// after analyzing the input:
	// I can say *none* of them will have a diagram of 0

	// state is [2]int = current & pressed
	queue := [][2]int{{0, 0}}

	for len(queue) > 0 {
		item := utils.MustShift(&queue)

		current, pressed := item[0], item[1]

		// get next buttons to press
		for k, v := range machine.button_wiring {
			if k&pressed == k {
				// already pressed
				continue
			}
			// only press if it will impact either current or the diagram
			compare := current | machine.light_diagram

			if v&compare != 0 {
				// queue next state
				// toggles lights
				next := current ^ v
				next_pressed := pressed | k

				if next == machine.light_diagram {
					// wow, we're done
					return bits.OnesCount(uint(next_pressed))
				}

				queue = append(queue, [2]int{next, next_pressed})
			}
		}
	}

	return
}

func (machine *Machine) DebugButtonsInt(prefix string, buttons int) {
	// initial_buttons := buttons
	out := [][]int{}
	for i := 0; buttons > 0; i++ {
		if buttons&1 == 1 {
			// action button i
			out = append(out, machine.plain_buttons[i])
		}
		buttons >>= 1
	}
	// log.Println(prefix, initial_buttons, out)
}

// TODO: terrible name
// this gets all possible button presses that produces a given diagram:
// [##.#]
func (machine *Machine) AllFewest(diagram int) iter.Seq[int] {
	if cache, found := machine.cache[diagram]; found {
		// return cached sequence
		return slices.Values(cache)
	}

	machine.cache[diagram] = []int{}

	// I don't know what I'm doing
	set_cache := func(buttons int) {
		machine.cache[diagram] = append(machine.cache[diagram], buttons)
	}

	return func(yield func(int) bool) {
		// after analyzing the input:
		// I can say *none* of them will have a diagram of 0

		// state is [2]int = current & pressed
		queue := [][2]int{{0, 0}}
		sent := map[int]struct{}{}

		for len(queue) > 0 {
			item := utils.MustPop(&queue)

			current, pressed := item[0], item[1]

			// get next buttons to press
			for k, v := range machine.button_wiring {
				if k&pressed == k {
					// already pressed
					continue
				}
				// only press if it will impact either current or the diagram
				compare := current | diagram

				if v&compare != 0 {
					// queue next state
					// toggles lights
					next := current ^ v
					next_pressed := pressed | k

					if _, found := sent[next_pressed]; found {
						continue
					}

					if next == diagram {
						// wow, we're done
						sent[next_pressed] = struct{}{}
						set_cache(next_pressed)
						if !yield(next_pressed) {
							return
						}
					}

					queue = append(queue, [2]int{next, next_pressed})
				}
			}
		}
	}
}

type DebugPrev struct {
	last    []int
	buttons int
}

type State struct {
	cur []int
	// when we divide by 2, we increase the magnitude of each press
	magnitude int
	presses   int
	prev      []DebugPrev
}

// create new diagram from joltage
// figure out which buttons create that diagram
// divide joltage by at least half (they'll all be even)
// determine which buttons to press next
func (machine *Machine) FewestButtonPressJoltage() (fewest int) {

	fewest = math.MaxInt

	queue := []State{
		{
			cur:       machine.joltage,
			magnitude: 0,
			presses:   0,
			prev: []DebugPrev{
				{last: machine.joltage, buttons: -1},
			}},
	}

	for len(queue) > 0 {
		state := utils.MustShift(&queue)

		// if presses already greater than fewest, we can ignore
		if state.presses > fewest {
			continue
		}

		cur := state.cur

		// check if all same number (and/or zeroes)
		same_number := -1
		// only if all same number
		// [ 0 0 2 2 ]
		// 1 1 0 0
		nonzero_binary := 0
		// example: [3 5 4 7]
		// converts (reversed) to:
		//
		// [ # . # # ]
		//
		//	7 4 5 3
		odd_binary := 0

		for i, num := range cur {
			if num%2 == 1 {
				odd_binary += 1 << i
			}
			if num != 0 {
				nonzero_binary += 1 << i
				if same_number == -1 {
					same_number = num
				} else if same_number != num {
					same_number = -2 // ignore
				}
			}
		}

		// fmt.Println("cur", state)
		// fmt.Printf("same: %v, nonzero: %b, odd: %b\n", same_number, nonzero_binary, odd_binary)

		// is done
		if nonzero_binary == 0 {
			// we've filled all the joltage
			if state.presses < fewest {
				fewest = state.presses
				// fmt.Println("Fewest State", state.presses)
				// fmt.Println("Prev:", state.prev)
				for _, prev := range state.prev[1:] {
					machine.DebugButtonsInt("--", prev.buttons)
				}
			}
			continue
		}

		// special multiplier for when all nonzeroes are the same number
		multiplier := 1

		// what to get buttons for
		diagram := 0

		// we're pretty much done here
		if same_number > 0 {
			diagram = nonzero_binary

			// need to split for all possible
			multiplier = same_number
		} else if odd_binary == 0 {
			// is even (divide by two, increase magnitude)
			next := make([]int, len(cur))
			copy(next, cur)

			for i := range next {
				// shift right is equivalent to div 2
				// now I'm just intentionally being obtuse
				next[i] >>= 1
			}
			state.magnitude++

			prev_states := make([]DebugPrev, len(state.prev))

			copy(prev_states, state.prev)

			state.prev = append(prev_states, DebugPrev{
				next,
				-2,
			})

			state.cur = next

			queue = append(queue, state)

			continue
		} else {
			// is odd
			// get single buttons that match this shape
			diagram = odd_binary
		}

		// for each of all possible, press the buttons
		// and save the state to some BFS
	outer:
		for buttons := range machine.AllFewest(diagram) {

			// press the buttons, and adjust joltage
			next := make([]int, len(cur))
			copy(next, cur)
			presses := 0

			// machine.DebugButtonsInt("buttons", buttons)

			prev_buttons := buttons

			for i := 0; buttons > 0; i++ {
				if buttons&1 == 1 {
					// action button i
					for _, j := range machine.plain_buttons[i] {
						// check if we can decrease joltage
						if next[j] == 0 {
							continue outer
						}
						// decrease joltage (multiplier is for final operation)
						next[j] -= multiplier
					}
					// increase presses using magnitude (and multiplier)
					presses += int(math.Pow(2, float64(state.magnitude))) * multiplier
				}
				buttons >>= 1
			}

			// log.Println("presses", presses, "stamped", stamped)

			prev_states := make([]DebugPrev, len(state.prev))

			copy(prev_states, state.prev)

			prev_states = append(prev_states, DebugPrev{
				next,
				prev_buttons,
			})

			queue = append(queue, State{
				next,
				state.magnitude,
				state.presses + presses,
				prev_states,
			})
		}
	}

	return
}

// should be a good enough hash?
func Hashed(in []int) (out any) {
	return fmt.Sprintf("%v", in)
}

func (machine *Machine) GetFewestRecursive() (fewest int) {
	return machine.RecursiveJoltagePresses(NextState{
		joltage:   machine.joltage,
		presses:   0,
		magnitude: 0,
	})
}

type NextState struct {
	presses int
	joltage []int
	// magnitude indicates that the path has been divided by 2,
	// and all presses should count extra
	magnitude int
}

func (machine *Machine) GetNextButtonPresses(state NextState) (states []NextState) {
	// next possible:
	// 1. same number (other than 0's)
	// 2. all even numbers (increase magnitude)
	// 3. get buttons that press odd

	// check if all same number (and/or zeroes)
	same_number := -1
	// only if all same number
	// [ 0 0 2 2 ]
	// 1 1 0 0
	nonzero_binary := 0
	// example: [3 5 4 7]
	// converts (reversed) to:
	//
	// [ # . # # ]
	//
	//	7 4 5 3
	odd_binary := 0

	for i, num := range state.joltage {
		if num%2 == 1 {
			odd_binary += 1 << i
		}
		if num != 0 {
			nonzero_binary += 1 << i
			if same_number == -1 {
				same_number = num
			} else if same_number != num {
				same_number = -2 // ignore
			}
		}
	}

	// special multiplier for when all nonzeroes are the same number
	multiplier := 1
	// what to get buttons for
	diagram := 0

	if same_number > 0 {
		multiplier = same_number
		diagram = nonzero_binary

	} else if odd_binary == 0 {
		// is even
		next := make([]int, len(state.joltage))
		copy(next, state.joltage)

		for i := range next {
			// shift right is equivalent to div 2
			// now I'm just intentionally being obtuse
			next[i] >>= 1
		}
		return []NextState{
			{presses: 0, joltage: next, magnitude: state.magnitude + 1},
		}
	} else {
		// odd
		diagram = odd_binary
	}

outer:
	for buttons := range machine.AllFewest(diagram) {
		next := make([]int, len(state.joltage))
		copy(next, state.joltage)
		presses := 0

		for i := 0; buttons > 0; i++ {
			if buttons&1 == 1 {
				// action button i
				for _, j := range machine.plain_buttons[i] {
					// check if we can decrease joltage
					if next[j] == 0 {
						continue outer
					}
					// decrease joltage (multiplier is for final operation)
					next[j] -= multiplier
				}
				// increase presses using magnitude (and multiplier)
				presses += int(math.Pow(2, float64(state.magnitude))) * multiplier
			}
			buttons >>= 1
		}

		states = append(states, NextState{
			joltage:   next,
			presses:   presses,
			magnitude: state.magnitude,
		})
	}

	return
}

const INVALID = -1

func (machine *Machine) RecursiveJoltagePresses(state NextState) (fewest int) {
	all_zero := true
	for num := range slices.Values(state.joltage) {
		if num != 0 {
			all_zero = false
		}
	}
	// if end
	if all_zero {
		return state.presses
	}

	// get next presses and positions from possible buttons
	fewest = math.MaxInt
	for _, next := range machine.GetNextButtonPresses(state) {
		// then recurse to find fewest presses
		next_presses := machine.RecursiveJoltagePresses(next)

		if next_presses == INVALID {
			continue
		}

		if state.presses+next_presses < fewest {
			fewest = state.presses + next_presses
		}
	}

	if fewest < 0 || fewest == math.MaxInt {
		return INVALID
	}

	return
}

//
// GAUSSIAN ELIMINATION
//

// by: icub3d
// https://www.youtube.com/watch?v=xibCHVRF6oI

type Matrix struct {
	RREF   [][]float64
	Pivots []int
	Free   []int
}

// from CHATGPT
func RREF(A [][]float64) Matrix {
	m := len(A)    // height
	n := len(A[0]) // width

	// Copy matrix
	M := make([][]float64, m)
	for i := range A {
		M[i] = append([]float64(nil), A[i]...)
	}

	row := 0
	pivots := []int{}

	for col := 0; col < n-1 && row < m; col++ {
		// Find pivot row
		pivot := row
		maxVal := math.Abs(M[row][col])
		for i := row + 1; i < m; i++ {
			if math.Abs(M[i][col]) > maxVal {
				maxVal = math.Abs(M[i][col])
				pivot = i
			}
		}

		// if the best value is 0 it's a free variable
		if math.Abs(M[pivot][col]) < 1e-12 {
			continue // No pivot in this column
		}

		// Swap pivot row into place
		M[row], M[pivot] = M[pivot], M[row]

		// Normalize pivot row
		p := M[row][col]
		for j := col; j < n; j++ {
			M[row][j] /= p
		}

		// Eliminate above and below
		for i := range m {
			if i == row {
				continue
			}
			f := M[i][col]
			for j := col; j < n; j++ {
				M[i][j] -= f * M[row][j]
			}
		}

		pivots = append(pivots, col)
		row++
	}

	// Determine free columns
	free := []int{}
	varMap := make(map[int]bool)
	for _, p := range pivots {
		varMap[p] = true
	}
	for c := 0; c < n-1; c++ {
		if !varMap[c] {
			free = append(free, c)
		}
	}

	return Matrix{
		RREF:   M,
		Pivots: pivots,
		Free:   free,
	}
}

func (machine *Machine) ToMatrix() [][]float64 {
	out := make([][]float64, len(machine.joltage))

	row_len := len(machine.plain_buttons) + 1

	for i := range out {
		out[i] = make([]float64, row_len)
		out[i][row_len-1] = float64(machine.joltage[i])
	}

	for col, buttons := range machine.plain_buttons {
		for _, row := range buttons {
			out[row][col] = 1
		}
	}

	return out
}

// do we need this?
const EPSILON float64 = 1e-9

func (matrix *Matrix) is_valid(values *[]int) (total int, valid bool) {
	cols := len(matrix.RREF[0])
	for row := range len(matrix.Pivots) {
		joltage := matrix.RREF[row][cols-1]
		// fmt.Println("Finding joltage", joltage)
		for i, col := range matrix.Free {
			joltage -= matrix.RREF[row][col] * float64((*values)[i])
		}

		if joltage < -EPSILON {
			// checking if it's negative
			return 0, false
		}
		rounded := math.Round(joltage)

		if math.Abs(joltage-rounded) > EPSILON {
			// checking if it's a float
			return 0, false
		}

		total += int(rounded)
	}

	free_total := utils.Sum(*values)

	return total + free_total, true
}

func dfs(matrix Matrix, index int, values *[]int, min *int, max int) {
	if index == len(matrix.Free) {
		if total, valid := matrix.is_valid(values); valid {
			if total < *min {
				*min = total
			}
		}
		return
	}

	sum := utils.Sum((*values)[:index])

	for val := range max {
		if sum+val >= *min {
			break
		}
		(*values)[index] = val
		dfs(matrix, index+1, values, min, max)
	}
}

// TODO:
// try guassian elimination and linear equations
func (machine *Machine) GaussianElimination() (fewest int) {
	A := machine.ToMatrix()

	res := RREF(A)

	log.Println("")
	for _, r := range res.RREF {
		log.Println(r)
	}

	log.Println("\nPivot columns:", res.Pivots)
	log.Print("Free columns:", res.Free, "\n\n")

	// GET MAX FROM JOLTAGES
	// ? first time using slices.Max
	max := slices.Max(machine.joltage)
	min := math.MaxInt
	values := make([]int, len(res.Free))

	// DFS
	dfs(res, 0, &values, &min, max)

	return min
}
