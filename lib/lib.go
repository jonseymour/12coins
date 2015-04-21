package lib

import (
	"fmt"
)

type Weight int

const (
	Light Weight = iota
	Equal
	Heavy
)

func (w Weight) String() string {
	switch w {
	case Light:
		return "light"
	case Heavy:
		return "heavy"
	case Equal:
		return "equal"
	default:
		return "invalid"
	}
}

func (w Weight) Invert() Weight {
	return Heavy - w
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

type Solution func(Scale) (int, Weight)

func NewOracle(coin int, w Weight) *Oracle {
	return &Oracle{
		coin:   coin,
		weight: w,
	}
}

func (o *Oracle) fail(err error) {
	o.err = err
	panic(err)
}

func (o *Oracle) check(a []int, b []int) {
	seen := [12]bool{}
	if o.attempts == 3 {
		o.fail(fmt.Errorf("too many attempts to use the scale!"))
	}
	for _, e := range a {
		if e < 0 || e > 11 {
			o.fail(fmt.Errorf("invalid coin: %d", e))
		}
		if seen[e] {
			o.fail(fmt.Errorf("duplicate detected: %d", e))
		} else {
			seen[e] = true
		}
	}
	for _, e := range b {
		if seen[e] {
			o.fail(fmt.Errorf("duplicate detected: %d", e))
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
			return Heavy - o.weight
		}
	}

	return Equal
}

// test checks whether decide answers the right coin for a given coin and relative weight
func Test(i int, w Weight, p Solution) error {
	oracle := &Oracle{coin: i, weight: w}
	func() {
		defer func() {
			if err := recover(); err != nil {
				oracle.err = err.(error)
			}
		}()
		ri, rw := p(oracle)
		if ri != i {
			panic(fmt.Errorf("decide chose coin %d", ri))
		}
		if rw != w {
			panic(fmt.Errorf("decide chose weight %v", rw))
		}
	}()
	return oracle.err
}
