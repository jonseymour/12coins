package lib

import (
	"fmt"
)

type flag uint

const (
	INVALID       flag = 0
	REVERSED           = 1 << 1
	GROUPED            = 1 << 2
	ANALYSED           = 1 << 3
	RELABLED           = 1 << 4
	NORMALISED         = 1 << 5
	CANONICALISED      = 1 << 6
)

// Describes a test failure. A test failure is an instance of a coin and weight such that the
// results of the weighings for that coin and weight are indistiguishable from some other coin
// and weight.
type Failure struct {
	Coin   int    `json:"coin"`
	Weight Weight `json:"weight"`
}

// Describes a possibly invalid solution to the 12 coins problem.
type Solution struct {
	encoding
	Weighings [3]Weighing  `json:"-"`
	Coins     []int        `json:"coins,omitempty"`     // a mapping between abs(27*a+9*b+c-13)-1 and the coin identity
	Weights   []Weight     `json:"weights,omitempty"`   // a mapping between sgn(27*a+9*b+c-13)-1 and the coin weight
	ZeroCoin  int          `json:"zero-coin,omitempty"` // the zero coin of the weighings. either 0 or 1.
	Unique    CoinSet      `json:"-"`                   // the coins that appear in one weighing
	Pairs     [3]CoinSet   `json:"-"`                   // the pairs that appear in exactly two weighings
	Triples   CoinSet      `json:"-"`                   // the coins that appear in all 3 weighings
	Flip      *int         `json:"flip,omitempty"`      // the weighing which needs to be flipped to guarantee abs(27*a+9*b+c-13)-1 is between 0 and 11
	Valid     *bool        `json:"valid,omitempty"`     // true if the solution is valid
	Failures  []Failure    `json:"failures,omitempty"`  // a list of tests for which the solution is ambiguous
	Structure [3]Structure `json:"-"`                   // the structure of the permutation
	flags     flag         // as
}

// Decide the relative weight of a coin by generating a linear combination of the three weighings and using
// this to index the array.
func (s *Solution) decide(scale Scale) (int, Weight, int) {
	scale.SetZeroCoin(s.ZeroCoin)

	results := [3]Weight{}

	results[0] = scale.Weigh(s.Weighings[0].Left().AsCoins(s.ZeroCoin), s.Weighings[0].Right().AsCoins(s.ZeroCoin))
	results[1] = scale.Weigh(s.Weighings[1].Left().AsCoins(s.ZeroCoin), s.Weighings[1].Right().AsCoins(s.ZeroCoin))
	results[2] = scale.Weigh(s.Weighings[2].Left().AsCoins(s.ZeroCoin), s.Weighings[2].Right().AsCoins(s.ZeroCoin))

	if s.Flip != nil {
		results[*s.Flip] = Heavy - results[*s.Flip]
	}

	a := results[0]
	b := results[1]
	c := results[2]

	i := int(a*9 + b*3 + c - 13) // must be between 0 and 26, inclusive.
	o := abs(i)
	if len(s.Coins) == 12 {
		if o < 1 || o > 12 {
			// this can only happen if flip hasn't be set correctly.
			panic(fmt.Errorf("index out of bounds: %d, %v", o, []Weight{a, b, c}))
		}
		o = o - 1
	} else {
		o = i + 13
	}

	f := s.Coins[o]
	w := s.Weights[o]

	if i > 0 {
		w = Heavy - w
	}

	return f, w, o
}

func (s *Solution) resetAnalysis() {
	s.Unique = nil
	s.Triples = nil
	s.Pairs = [3]CoinSet{nil, nil, nil}
	s.Structure = [3]Structure{nil, nil, nil}
	s.flags = s.flags &^ (GROUPED | ANALYSED | CANONICALISED)
}

// Invoke the internal decide method to decide which coin
// is counterfeit and what it's relative weight is.
func (s *Solution) Decide(scale Scale) (int, Weight) {
	if s.flags&REVERSED == 0 {
		panic(fmt.Errorf("This solution must be reversed first."))
	}
	f, w, _ := s.decide(scale)
	return f, w
}

// Configure the zero coin of the solution.
func (s *Solution) SetZeroCoin(coin int) {
	s.ZeroCoin = coin
}

// Create a deep clone of the receiver.
func (s *Solution) Clone() *Solution {
	tmp := s.Flip
	if tmp != nil {
		tmp = pi(*tmp)
	}
	v := s.Valid
	if v != nil {
		v = pbool(*v)
	}
	clone := Solution{
		Weighings: [3]Weighing{},
		Coins:     make([]int, len(s.Coins)),
		Weights:   make([]Weight, len(s.Weights)),
		ZeroCoin:  s.ZeroCoin,
		Unique:    s.Unique,
		Triples:   s.Triples,
		Failures:  make([]Failure, len(s.Failures)),
		Flip:      tmp,
		Valid:     v,
		flags:     s.flags,
	}

	copy(clone.Pairs[0:], s.Pairs[0:])
	copy(clone.Weighings[0:], s.Weighings[0:])
	copy(clone.Coins[0:], s.Coins[0:])
	copy(clone.Weights[0:], s.Weights[0:])
	copy(clone.Failures[0:], s.Failures[0:])
	copy(clone.Structure[0:], s.Structure[0:])

	return &clone
}

// Sort the coins in each weighing in increasing numerical order.
func (s *Solution) Normalize() *Solution {
	clone := s.Clone()

	for i, _ := range clone.Weighings {
		clone.Weighings[i] = NewWeighing(clone.Weighings[i].Pan(0).Sort(), clone.Weighings[i].Pan(1).Sort())
	}
	clone.flags |= NORMALISED &^ (CANONICALISED)
	return clone
}
