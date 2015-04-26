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

type flag uint

const (
	REVERSED flag = 1
	GROUPED
	ANALYSED
	RELABLED
	NORMALISED
	CANONICALISED
)

// Describes a test failure. A test failure is an instance of a coin and weight such that the
// results of the weighings for that coin and weight are indistiguishable from some other coin
// and weight.
type Failure struct {
	Coin   int    `json:"coin"`
	Weight Weight `json:"weight"`
}

type encoding struct {
	Weighings *[3][2][]int `json:"weighings,omitempty"` // the weighings of the solution
	Unique    *[]int       `json:"unique,omitempty"`    // the coins that appear in one weighing
	Pairs     *[3][2]int   `json:"pairs,omitempty"`     // the pairs that appear in exactly two weighings
	Triples   *[]int       `json:"triples,omitempty"`   // the coins that appear in all 3 weighings
	Structure *[3]string   `json:"structure,omitempty"` // An encoding of the structure
}

// Describes a possibly invalid solution to the 12 coins problem.
type Solver struct {
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
func (s *Solver) decide(scale Scale) (int, Weight, int) {
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
		Weighings: [3]Weighing{},
		Coins:     make([]int, len(s.Coins), len(s.Coins)),
		Weights:   make([]Weight, len(s.Weights), len(s.Weights)),
		ZeroCoin:  s.ZeroCoin,
		Unique:    s.Unique,
		Triples:   s.Triples,
		Flip:      tmp,
		Valid:     v,
		flags:     s.flags,
	}

	copy(clone.Pairs[0:], s.Pairs[0:])
	copy(clone.Weighings[0:], s.Weighings[0:])

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

	for i, w := range clone.Weighings {
		coinSet := [2]CoinSet{}
		for j, pan := range w.Pans() {
			coins := pan.AsCoins(clone.ZeroCoin)
			for k, e := range coins {
				coins[k] = p.Index(e) + clone.ZeroCoin
			}
			sort.Sort(sort.IntSlice(coins))
			coinSet[j] = NewCoinSet(coins, clone.ZeroCoin)
		}
		clone.Weighings[i] = NewWeighing(coinSet[0], coinSet[1])
	}

	for i, _ := range clone.Coins {
		clone.Coins[i] = i + clone.ZeroCoin
	}

	clone.flags |= (RELABLED | NORMALISED) &^ (CANONICALISED | GROUPED | ANALYSED)

	return clone, nil
}

// Sort the coins in each weighing in increasing numerical order.
func (s *Solver) Normalize() *Solver {
	clone := s.Clone()

	for i, _ := range clone.Weighings {
		clone.Weighings[i] = NewWeighing(clone.Weighings[i].Pan(0).Sort(), clone.Weighings[i].Pan(1).Sort())
	}
	clone.flags |= NORMALISED &^ (CANONICALISED)
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
	clone.flags |= REVERSED
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
	s.Unique = nil
	s.Triples = nil
	s.Pairs = [3]CoinSet{nil, nil, nil}
}

// Tabulate the singletons, pairs and triples of the solution.
func (s *Solver) Groupings() (*Solver, error) {
	clone := s.Clone()
	clone.resetCounts()

	a := clone.Weighings[0].Both()
	b := clone.Weighings[1].Both()
	c := clone.Weighings[2].Both()

	ab := a.Intersection(b)
	bc := b.Intersection(c)
	ca := c.Intersection(a)

	all := a.Union(b).Union(c)
	triples := ab.Intersection(bc).Intersection(ca)
	singletons := a.Complement(b.Union(c)).Union(b.Complement(a.Union(c))).Union(c.Complement(a.Union(b)))
	pairs := all.Complement(triples).Complement(singletons)

	if triples.Size() != 3 || singletons.Size() != 3 || pairs.Size() != 6 {
		s.Valid = pbool(false)
		return s, fmt.Errorf("invalid grouping sizes: %v, %v, %v", triples, pairs, singletons)
	}

	abp := pairs.Intersection(ab)
	bcp := pairs.Intersection(bc)
	cap := pairs.Intersection(ca)

	clone.Triples = triples
	clone.Unique = singletons
	clone.Pairs = [3]CoinSet{abp, bcp, cap}
	clone.flags |= GROUPED

	return clone, nil
}

// Return a clone of the receiver in which the structure has been populated.
func (s *Solver) AnalyseStructure() (*Solver, error) {
	var r *Solver
	var err error

	if s.Unique == nil || s.Triples == nil || s.Pairs[0] == nil || s.Pairs[1] == nil || s.Pairs[2] == nil {
		r, err = s.Groupings()
	} else {
		r = s.Clone()
	}

	if err != nil {
		return r, err
	}
	r.flags |= ANALYSED

	return r, nil
}

// Return a clone of the receiver in which the weighings have been permuted into the
// the canonical order and all sets are ordered sets.
func (s *Solver) Canonical() (*Solver, error) {
	var r *Solver
	var err error

	if s.Structure[0] == nil || s.Structure[1] == nil || s.Structure[2] == nil {
		r, err = s.AnalyseStructure()
	} else {
		r = s.Clone()
	}

	if err != nil {
		return r, err
	}

	r.flags |= CANONICALISED
	return r, nil
}

// Convert the solution to its JSON representation.
func (s *Solver) String() string {
	s.Encode()
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Solver) Encode() {
	tmp := [3][2][]int{}
	s.encoding.Weighings = &tmp
	for i, w := range s.Weighings {
		for j, p := range w.Pans() {
			s.encoding.Weighings[i][j] = p.AsCoins(s.ZeroCoin)
		}
	}
	if s.Unique != nil {
		tmp := s.Unique.AsCoins(s.ZeroCoin)
		s.encoding.Unique = &tmp
	}
	if s.Triples != nil {
		tmp := s.Triples.AsCoins(s.ZeroCoin)
		s.encoding.Triples = &tmp
	}
	if s.Unique != nil {
		tmp := [3][2]int{}
		for i, _ := range tmp {
			if s.Pairs[i] != nil {
				copy(tmp[i][0:], s.Pairs[i].AsCoins(s.ZeroCoin))
			}
		}
		s.encoding.Pairs = &tmp
	}
	structure := [3]string{}
	count := 0
	for i, _ := range structure {
		if s.Structure[i] != nil {
			structure[i] = s.Structure[i].String()
			count += 1
		}
	}
	if count == 3 {
		s.encoding.Structure = &structure
	}
}

func (s *Solver) Decode() {
	if s.encoding.Weighings != nil {
		for i, w := range *s.encoding.Weighings {
			s.Weighings[i] = NewWeighing(NewCoinSet(w[0], s.ZeroCoin), NewCoinSet(w[1], s.ZeroCoin))
		}
	}
	if s.encoding.Unique != nil {
		s.Unique = NewCoinSet(*s.encoding.Unique, s.ZeroCoin)
	}
	if s.encoding.Triples != nil {
		s.Triples = NewCoinSet(*s.encoding.Triples, s.ZeroCoin)
	}
	if s.encoding.Pairs != nil {
		for i, p := range *s.encoding.Pairs {
			s.Pairs[i] = NewCoinSet(p[0:], s.ZeroCoin)
		}
	}
	if s.encoding.Structure != nil {
		for i, t := range *s.encoding.Structure {
			s.Structure[i], _ = ParseStructure(t, i)
		}
	}
}
