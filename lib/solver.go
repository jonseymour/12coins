package lib

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Calculate the absolute value of the specified integer.
func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

// Describes a test failure. A test failure is an instance of a coin and weight such that the
// results of the weighings for that coin and weight are indistiguishable from some other coin
// and weight.
type Failure struct {
	Coin   int    `json:"coin"`
	Weight Weight `json:"weight"`
}

// Describes a possibly invalid solution to the 12 coins problem.
type Solver struct {
	EncodedWeighings [3][2][]int   `json:"weighings,omitempty"` // the weighings of the solution
	Weighings        [3][2]CoinSet `json:"-",omit`
	Coins            []int         `json:"coins,omitempty"`     // a mapping between abs(27*a+9*b+c-13)-1 and the coin identity
	Weights          []Weight      `json:"weights,omitempty"`   // a mapping between sgn(27*a+9*b+c-13)-1 and the coin weight
	ZeroCoin         int           `json:"zero-coin,omitempty"` // the zero coin of the weighings. either 0 or 1.
	Unique           []int         `json:"unique,omitempty"`    // the coins that appear in one weighing
	Pairs            [][2]int      `json:"pairs,omitempty"`     // the pairs that appear in exactly two weighings
	Triples          []int         `json:"triples,omitempty"`   // the coins that appear in all 3 weighings
	Flip             *int          `json:"flip,omitempty"`      // the weighing which needs to be flipped to guarantee abs(27*a+9*b+c-13)-1 is between 0 and 11
	Valid            *bool         `json:"valid,omitempty"`     // true if the solution is valid
	Failures         []Failure     `json:"failures,omitempty"`  // a list of tests for which the solution is ambiguous
}

// Decide the relative weight of a coin by generating a linear combination of the three weighings and using
// this to index the array.
func (s *Solver) decide(scale Scale) (int, Weight, int) {
	scale.SetZeroCoin(s.ZeroCoin)

	results := [3]Weight{}

	results[0] = scale.Weigh(s.Weighings[0][0].AsCoins(s.ZeroCoin), s.Weighings[0][1].AsCoins(s.ZeroCoin))
	results[1] = scale.Weigh(s.Weighings[1][0].AsCoins(s.ZeroCoin), s.Weighings[1][1].AsCoins(s.ZeroCoin))
	results[2] = scale.Weigh(s.Weighings[2][0].AsCoins(s.ZeroCoin), s.Weighings[2][1].AsCoins(s.ZeroCoin))

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

// Invoke the internal decide method to decide which coin
// is counterfeit and what it's relative weight is.
func (s *Solver) Decide(scale Scale) (int, Weight) {
	f, w, _ := s.decide(scale)
	return f, w
}

// Configure the zero coin of the solution.
func (s *Solver) SetZeroCoin(coin int) {
	s.ZeroCoin = coin
}

// Create a deep clone of the receiver.
func (s *Solver) Clone() *Solver {
	tmp := s.Flip
	if tmp != nil {
		tmp = pi(*tmp)
	}
	v := s.Valid
	if v != nil {
		v = pbool(*v)
	}
	clone := Solver{
		Weighings: [3][2]CoinSet{},
		Coins:     make([]int, len(s.Coins), len(s.Coins)),
		Weights:   make([]Weight, len(s.Weights), len(s.Weights)),
		ZeroCoin:  s.ZeroCoin,
		Unique:    append([]int{}, s.Unique...),
		Triples:   append([]int{}, s.Triples...),
		Pairs:     append([][2]int{}, s.Pairs...),
		Flip:      tmp,
		Valid:     v,
	}

	for j, _ := range []int{0, 1} {
		for i, _ := range s.Weighings {
			clone.Weighings[i][j] = s.Weighings[i][j]
		}
	}
	for i, e := range s.Coins {
		clone.Coins[i] = e
	}

	for i, e := range s.Weights {
		clone.Weights[i] = e
	}

	return &clone
}

// Relabel the coins of the weighing such that the Coins
// slice is numbered in strictly increasing numerical order.
func (s *Solver) Relabel() (*Solver, error) {

	var clone *Solver
	var err error

	if len(s.Coins) != 12 {
		if clone, err = s.Reverse(); err != nil {
			s.Valid = pbool(false)
			return s, err
		}
	} else {
		clone = s.Clone()
	}

	clone.resetCounts()

	c := make([]int, len(clone.Coins), len(clone.Coins))
	for i, e := range clone.Coins {
		c[i] = e
	}
	p := NewPermutation(c, clone.ZeroCoin)

	for i, _ := range clone.Weighings {
		for j, _ := range []int{0, 1} {
			coins := clone.Weighings[i][j].AsCoins(clone.ZeroCoin)
			for k, e := range coins {
				coins[k] = p.Index(e) + clone.ZeroCoin
			}
			sort.Sort(sort.IntSlice(coins))
			clone.Weighings[i][j] = NewCoinSet(coins, clone.ZeroCoin)
		}
	}

	for i, _ := range clone.Coins {
		clone.Coins[i] = i + clone.ZeroCoin
	}

	return clone, nil
}

// Sort the coins in each weighing in increasing numerical order.
func (s *Solver) Normalize() *Solver {
	clone := s.Clone()

	for i, _ := range clone.Weighings {
		for j, _ := range []int{0, 1} {
			clone.Weighings[i][j] = clone.Weighings[i][j].Sort()
		}
	}
	return clone
}

//
// If the receiver is a valid solution to the 12 coins problem,
// return a clone of the receiver in which the Coins and Weights
// slice and the Flip pointer have been populated with values
// required to make Decide(Scale) return the correct values for
// all inputs.
//
// Otherwise, return a pointer to the receiver.
//
// The .Valid value of the returned pointer is always non nil
// and always indicates whether the receiver was a valid solution.
//
// If err is non-nil, then .Valid of the result will point to a false
// value and .Failures of the result will list the tests that
// caused the receiver to be marked invalid.
//
func (s *Solver) Reverse() (*Solver, error) {
	clone := s.Clone()

	clone.Coins = make([]int, 27, 27)
	clone.Weights = make([]Weight, 27, 27)
	clone.Flip = nil

	for i, _ := range clone.Coins {
		clone.Coins[i] = clone.ZeroCoin
		clone.Weights[i] = Equal
	}

	failures := make(map[int]bool)

	fail := func(coin int, weight Weight) {
		if !failures[coin] {
			failures[coin] = true
			s.Failures = append(s.Failures, Failure{
				Coin:   coin,
				Weight: weight,
			})
		}
	}

	for _, w := range []Weight{Light, Heavy} {
		for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
			o := NewOracle(i, w, 1)
			ri, _, rx := clone.decide(o)
			if ri != i {
				if clone.Weights[rx] != Equal {
					fail(clone.Coins[rx], clone.Weights[rx])
					fail(i, w)
					continue
				}
				clone.Coins[rx] = i
			}
			clone.Weights[rx] = w
		}
	}

	if len(s.Failures) != 0 {
		s.Valid = pbool(false)
		return s, fmt.Errorf("not a valid solution because of %d failures", len(s.Failures))
	}

	// exploit symmetry where it exists

	if clone.Weights[0] == Equal {
		clone.Coins = clone.Coins[1:13]
		clone.Weights = clone.Weights[1:13]
	} else {

		//
		// A curious truth is that within the first 13 positions, of the
		// array, the unassigned slots can only ever occur at positions
		// 0,2,6 or 8.
		//
		// If it occurs an 0, then LLL is not a valid weighing. If it occurs
		// at 2, then LLH is not a valid weighing. If it occurs at 6, then LHL
		// is not a valid weighing. If it occurs at 8 then LHH is not a valid
		// weighing.
		//
		// The desired outcome is that the empty slots occur at
		// 0 which allows the sum derived from the other bits to
		// index the counterfeit coin directly.
		//
		// This can be arranged by identifying the weighing that is causing
		// the empty slot to happen at something other than 0 and flipping
		// the contribution of that weighing to the sum.
		//

		if clone.Weights[8] == Equal {
			clone.Flip = pi(0)
		} else if clone.Weights[6] == Equal {
			clone.Flip = pi(1)
		} else {
			clone.Flip = pi(2)
		}
		return clone.Reverse()
	}
	clone.Valid = pbool(true)
	return clone, nil
}

// convert an integer value into a pointer to that value.
func pi(i int) *int {
	return &i
}

// convert a boolean value into a pointer to that value.
func pbool(b bool) *bool {
	return &b
}

func (s *Solver) resetCounts() {
	s.Unique = []int{}
	s.Triples = []int{}
	s.Pairs = [][2]int{}
}

// Tabulate the singletons, pairs and triples of the solution.
func (s *Solver) Groupings() (*Solver, error) {
	clone := s.Clone()
	clone.Unique = []int{}
	clone.Triples = []int{}
	clone.Pairs = [][2]int{[2]int{-1, -1}, [2]int{-1, -1}, [2]int{-1, -1}}
	counts := make(map[int]int)
	sets := make(map[int]int)
	pairs := []int{}
	for i := 0; i < 12; i++ {
		counts[i] = 0
		sets[i] = 0
	}
	for i, _ := range clone.Weighings {
		for j, _ := range []int{0, 1} {
			for _, e := range clone.Weighings[i][j].AsCoins(clone.ZeroCoin) {
				x := e - s.ZeroCoin
				if x < 0 || x > 11 {
					s.Valid = pbool(false)
					return s, fmt.Errorf("invalid coin detected at %d, %d -> %d", i, j, e)
				}
				counts[x] += 1
				sets[x] |= (1 << uint(i))
			}
		}
	}
	for k, v := range counts {
		switch v {
		case 1:
			clone.Unique = append(clone.Unique, k)
		case 2:
			pairs = append(pairs, k)
		case 3:
			clone.Triples = append(clone.Triples, k)
		default:
			s.Valid = pbool(false)
			return s, fmt.Errorf("invalid count detected for coin %d -> %d", k, v)
		}
	}

	sort.Sort(sort.IntSlice(pairs))

	for i, _ := range clone.Weighings {
		for _, k := range pairs {
			mask := 1<<uint(i) | 1<<uint((i+1)%3)
			if sets[k]&mask == mask {
				if clone.Pairs[i][0] < 0 {
					clone.Pairs[i][0] = k
				} else {
					clone.Pairs[i][1] = k
				}
			}
		}
	}

	sort.Sort(sort.IntSlice(clone.Unique))
	sort.Sort(sort.IntSlice(clone.Triples))

	return clone, nil
}

// Convert the solution to its JSON representation.
func (s *Solver) String() string {
	s.Encode()
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Solver) Encode() {
	for i, w := range s.Weighings {
		for j, p := range w {
			s.EncodedWeighings[i][j] = p.AsCoins(s.ZeroCoin)
		}
	}
}

func (s *Solver) Decode() {
	for i, w := range s.EncodedWeighings {
		for j, p := range w {
			s.Weighings[i][j] = NewCoinSet(p, s.ZeroCoin)
		}
	}
}
