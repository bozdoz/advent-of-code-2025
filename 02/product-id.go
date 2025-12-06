package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bozdoz/advent-of-code-2025/utils"
)

type ProductRange struct {
	min, max, sum int
}

func NewProductRange(rng string) *ProductRange {
	items := strings.SplitN(rng, "-", 2)

	min, max := utils.ParseInt(items[0]), utils.ParseInt(items[1])

	sum := 0

	return &ProductRange{
		min,
		max,
		sum,
	}
}

// on modern architectures, multiplying a number is much faster than dividing a number
// https://stackoverflow.com/a/68124773/488784
func digitCount(i int) int {
	if i >= 1e18 {
		return 19
	}
	x, count := 10, 1
	for x <= i {
		x *= 10
		count++
	}
	return count
}

func hasDouble(val int) bool {
	digits := digitCount(val)

	// needs to be exact copy
	if digits%2 != 0 {
		return false
	}

	pow := int(math.Pow10(digits / 2))
	first := val / pow
	second := val - first*pow

	return first == second
}

// this was meant to be clever, but went haywire
// and then worse: it passed all the tests, but it didn't work
// so I made hasRepeatingString, and compared which values didn't match
func hasRepeating(val int) bool {
	digits := digitCount(val)

	is_even := digits%2 == 0

	// loop over this, decreasing in segment size
	half := digits / 2
	// what am I doing
	true_half := half

	if !is_even && digits > 2 {
		// round up
		half++
	}

outer:
	// ? first time doing range int
	for i := range half {
		num := half + i
		pow := int(math.Pow10(num))
		compare := val / pow

		size := true_half - i

		if size == 0 || digits%size != 0 {
			// check if it can evenly appear in val
			continue
		}

		// now look for the same value in each chunk
		rest := val - compare*pow

		// new pow for getting each right-most chunk
		pow = int(math.Pow10(size))

		for range digits/size - 1 {
			// get right-most chunk
			next := rest % pow

			if next != compare {
				continue outer
			}

			// update rest to omit right-most chunk
			rest /= pow

		}
		return true
	}

	return false
}

// go regex doesn't support back references (\1) 😭
func hasRepeatingRegex(val int) bool {
	str := fmt.Sprintf("%d", val)
	half := len(str) / 2

	for i := range half {
		num := half - i
		// pattern := fmt.Sprintf("^(?:(\\d{%d})\\\\1*)$", num)
		pattern := fmt.Sprintf(`^(\d{%d})\\1*$`, num)

		re := regexp.MustCompile(pattern)

		if re.MatchString(str) {
			fmt.Println("Pattern matched!", pattern, str)
			return true
		} else {
			fmt.Println("Pattern did not match.", pattern, str)
		}
	}

	return false
}

func hasRepeatingString(val int) bool {
	// Sprintf slower
	// str := fmt.Sprintf("%d", val)

	str := strconv.FormatInt(int64(val), 10)
	digits := len(str)

outer:
	for chunk := 1; chunk <= digits/2; chunk++ {
		// check if it can evenly appear in val
		if digits%chunk != 0 {
			continue
		}
		segment := str[:chunk]

		// start at next segment, increase by chunk size
		for i := chunk; i < digits; i += chunk {
			if str[i:i+chunk] != segment {
				continue outer
			}
		}
		return true
	}
	return false
}

func (prod *ProductRange) InvalidCount(fun func(i int) bool) (count int) {
	for i := prod.min; i <= prod.max; i++ {
		if fun(i) {
			count++
			prod.sum += i
		}
	}
	return
}

// hasRepeating wasn't working, so I compared with the basic string compare
func (prod *ProductRange) CheckMethods() (count int) {
	for i := prod.min; i <= prod.max; i++ {
		if hasRepeating(i) != hasRepeatingString(i) {
			fmt.Println(i)
		}
	}
	return
}
