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
	RECURSE            = 1 << 7
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
	Coins     []int        `json:"coins,omitempty"`    // a mapping between abs(9*a+3*b+c-13)-1 and the coin identity
	Weights   []Weight     `json:"weights,omitempty"`  // a mapping between sgn(9*a+3*b+c-13)-1 and the coin weight
	Unique    CoinSet      `json:"-"`                  // the coins that appear in one weighing
	Pairs     [3]CoinSet   `json:"-"`                  // the pairs that appear in exactly two weighings
	Triples   CoinSet      `json:"-"`                  // the coins that appear in all 3 weighings
	Failures  []Failure    `json:"failures,omitempty"` // a list of tests for which the solution is ambiguous
	Structure [3]Structure `json:"-"`                  // the structure of the permutation
	flags     flag         // as
}

// Decide the relative weight of a coin by generating a linear combination of the three weighings and using
// this to index the array.
func (s *Solution) decide(scale Scale) (int, Weight, int) {
	z := s.GetZeroCoin()
	scale.SetZeroCoin(z)

	results := [3]Weight{}

	results[0] = scale.Weigh(s.Weighings[0].Left().AsCoins(z), s.Weighings[0].Right().AsCoins(z))
	results[1] = scale.Weigh(s.Weighings[1].Left().AsCoins(z), s.Weighings[1].Right().AsCoins(z))
	results[2] = scale.Weigh(s.Weighings[2].Left().AsCoins(z), s.Weighings[2].Right().AsCoins(z))

	if s.encoding.Flip != nil {
		results[*s.encoding.Flip] = Heavy - results[*s.encoding.Flip]
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

// The internal reset is used to reset the analysis of the receiver but
// does not undo the reversed state.
func (s *Solution) reset() {
	s.Unique = nil
	s.Triples = nil
	s.Pairs = [3]CoinSet{nil, nil, nil}
	s.Structure = [3]Structure{nil, nil, nil}
	s.encoding = encoding{
		ZeroCoin: s.encoding.ZeroCoin,
		Flip:     s.encoding.Flip,
	}
	s.flags = s.flags &^ (GROUPED | ANALYSED | CANONICALISED)
}

// The external reset creates a new clone in which only the weighings
// are preserved.
func (s *Solution) Reset() *Solution {
	r := s.Clone()
	r.reset()
	r.Coins = []int{}
	r.Weights = []Weight{}
	r.flags = INVALID
	r.encoding.Flip = nil
	return r
}

func (s *Solution) markInvalid() {
	s.flags = INVALID | (s.flags & RECURSE)
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
	if coin == ONE_BASED {
		s.encoding.ZeroCoin = nil
	} else {
		s.encoding.ZeroCoin = pi(coin)
	}
}

func (s *Solution) GetZeroCoin() int {
	if s.encoding.ZeroCoin == nil {
		return 1
	} else {
		return *s.encoding.ZeroCoin
	}
}

// Create a deep clone of the receiver.
func (s *Solution) Clone() *Solution {
	tmp := s.encoding.Flip
	if tmp != nil {
		tmp = pi(*tmp)
	}
	clone := Solution{
		encoding: encoding{
			ZeroCoin: s.encoding.ZeroCoin,
			Flip:     s.encoding.Flip,
		},
		Weighings: [3]Weighing{},
		Coins:     make([]int, len(s.Coins)),
		Weights:   make([]Weight, len(s.Weights)),
		Unique:    s.Unique,
		Triples:   s.Triples,
		Failures:  make([]Failure, len(s.Failures)),
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

// Answer true if the solution is a valid solution. This will be true if it could
// be successfully reversed, false otherwise.
func (s *Solution) IsValid() bool {
	if s.flags&REVERSED == 0 {
		c, err := s.Reverse()
		return err == nil && c.flags&REVERSED != 0
	} else {
		return true
	}
}
