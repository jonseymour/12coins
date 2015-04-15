package main

import (
	"fmt"
	"os"
)

type Weight int

const (
	light Weight = iota
	equal
	heavy
)

func (w Weight) String() string {
	switch w {
	case light:
		return "light"
	case heavy:
		return "heavy"
	case equal:
		return "equal"
	default:
		return "invalid"
	}
}

// A Scale can weigh two collections of coins and answer whether the
// a collection weighs less, the same as or more than the right collection.
//
// Each slice can contain the numbers 0, 11. A number that appears in one slice
// may not appear in another slice.
type Scale interface {
	Weigh(a []int, b []int) Weight
}

type Oracle struct {
	coin     int
	weight   Weight
	attempts int
	err      error
}

func (o *Oracle) fail(err error) {
	o.err = err
	panic(err)
}

func (o *Oracle) check(a []int, b []int) {
	seen := [12]bool{}
	if o.attempts == 3 {
		panic(fmt.Errorf("too many attempts to use the scale!"))
	}
	for _, e := range a {
		if e < 0 || e > 11 {
			panic(fmt.Errorf("invalid coin: %d", e))
		}
		if seen[e] {
			panic(fmt.Errorf("duplicate detected: %d", e))
		} else {
			seen[e] = true
		}
	}
	for _, e := range b {
		if seen[e] {
			panic(fmt.Errorf("duplicate detected: %d", e))
		} else {
			seen[e] = true
		}
	}
}

// The oracle implements the Scale interface and happens to know which
// coin is the different coin.
func (o *Oracle) Weigh(a []int, b []int) Weight {
	o.check(a, b)
	o.attempts += 1

	for _, e := range a {
		if e == o.coin {
			return o.weight
		}
	}

	for _, e := range b {
		if e == o.coin {
			switch o.weight {
			case light:
				return heavy
			case heavy:
				return light
			}
		}
	}

	return equal
}

// test checks whether decide answers the right coin for a given coin and relative weight
func test(i int, w Weight) bool {
	ri, rw := decide(&Oracle{coin: i, weight: w})
	return ri == i && rw == w
}

// exhaustively test the decision procedure against all possibilities and return those that fail
func main() {

	fail := false
	for i := 0; i < 12; i++ {
		if !test(i, heavy) {
			fmt.Fprintf(os.Stderr, "failed for %v, heavy\n", i)
			fail = true
		}
		if !test(i, light) {
			fmt.Fprintf(os.Stderr, "failed for %v, light\n", i)
			fail = true
		}
	}
	if fail {
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "ok\n")
	}
}
