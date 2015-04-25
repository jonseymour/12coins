package lib

import (
	"encoding/json"
	"fmt"
	"sort"
)

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

type Solver struct {
	Weighings [3][2][]int `json:"weighings"`
	Coins     []int       `json:"coins"`
	Weights   []Weight    `json:"weights"`
	ZeroCoin  int         `json:"zero-coin,omitempty"`
	Unique    []int       `json:"unique,omitempty"`
	Pairs     [][2]int    `json:"pairs,omitempty"`
	Triples   []int       `json:"triples,omitempty"`
	Flip      *int        `json:"flip,omitempty"`
	Valid     *bool       `json:"valid,omitempty"`
}

func (s *Solver) decide(scale Scale) (int, Weight, int) {
	scale.SetZeroCoin(s.ZeroCoin)

	results := [3]Weight{}

	results[0] = scale.Weigh(s.Weighings[0][0], s.Weighings[0][1])
	results[1] = scale.Weigh(s.Weighings[1][0], s.Weighings[1][1])
	results[2] = scale.Weigh(s.Weighings[2][0], s.Weighings[2][1])

	if s.Flip != nil {
		results[*s.Flip] = Heavy - results[*s.Flip]
	}

	a := results[0]
	b := results[1]
	c := results[2]

	i := int(a*9 + b*3 + c - 13)
	o := abs(i)
	if len(s.Coins) == 12 {
		if o < 1 || o > 12 {
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

func (s *Solver) Decide(scale Scale) (int, Weight) {
	f, w, _ := s.decide(scale)
	return f, w
}

func (s *Solver) SetZeroCoin(coin int) {
	s.ZeroCoin = coin
}

func (s *Solver) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

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
		Weighings: [3][2][]int{},
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
			p := make([]int, len(s.Weighings[i][j]), len(s.Weighings[i][j]))
			copy(p, s.Weighings[i][j])
			clone.Weighings[i][j] = p
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
			for k, e := range clone.Weighings[i][j] {
				clone.Weighings[i][j][k] = p.Index(e) + clone.ZeroCoin
			}
			sort.Sort(sort.IntSlice(clone.Weighings[i][j]))
		}
	}

	for i, _ := range clone.Coins {
		clone.Coins[i] = i + clone.ZeroCoin
	}

	return clone, nil
}

func (s *Solver) Normalize() *Solver {
	clone := s.Clone()

	for i, _ := range clone.Weighings {
		for j, _ := range []int{0, 1} {
			sort.Sort(sort.IntSlice(clone.Weighings[i][j]))
		}
	}
	return clone
}

func (s *Solver) Reverse() (*Solver, error) {
	clone := s.Clone()

	clone.Coins = make([]int, 27, 27)
	clone.Weights = make([]Weight, 27, 27)

	for i, _ := range clone.Coins {
		clone.Coins[i] = clone.ZeroCoin
		clone.Weights[i] = Equal
	}

	for _, w := range []Weight{Light, Heavy} {
		for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
			o := NewOracle(i, w, 1)
			ri, rw, rx := clone.decide(o)
			if ri != i {
				if clone.Weights[rx] != Equal {
					s.Valid = pbool(false)
					return s, fmt.Errorf("cannot distinguish between (%d, %v) and (%d, %v)", clone.Coins[rx], clone.Weights[rx], i, rw)
				}
				clone.Coins[rx] = i
			}
			clone.Weights[rx] = w
		}
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

func pi(i int) *int {
	return &i
}

func pbool(b bool) *bool {
	return &b
}

func (s *Solver) resetCounts() {
	s.Unique = []int{}
	s.Triples = []int{}
	s.Pairs = [][2]int{}
}

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
			for _, e := range clone.Weighings[i][j] {
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
