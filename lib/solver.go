package lib

import (
	"encoding/json"
	"sort"
)

type Solver struct {
	Permutation []int       `json:"permutation,omitempty"`
	Weighings   [3][2][]int `json:"weighings"`
	Coins       [12]int     `json:"coins"`
	Weights     [12]Weight  `json:"weights"`
}

func (s *Solver) Decide(scale Scale) (int, Weight) {
	a := scale.Weigh(s.Weighings[0][0], s.Weighings[0][1])
	b := scale.Weigh(s.Weighings[1][0], s.Weighings[1][1])
	c := scale.Weigh(s.Weighings[2][0], s.Weighings[2][1])

	i := a*9 + b*3 + c
	o := i

	if i > 12 {
		o = 26 - i
	}

	f := s.Coins[int(o-1)]
	w := s.Weights[int(o-1)]

	if i > 12 {
		w = Heavy - w
	}

	return f, w
}

func (s *Solver) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Solver) Clone() *Solver {
	clone := Solver{
		Weighings: [3][2][]int{},
		Coins:     [12]int{},
		Weights:   [12]Weight{},
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

	return s
}

func (s *Solver) Relabel() {
	c := make([]int, len(s.Coins), len(s.Coins))
	for i, e := range s.Coins {
		c[i] = e
	}
	p := NewPermutation(c)

	for i, _ := range s.Weighings {
		for j, _ := range []int{0, 1} {
			for k, e := range s.Weighings[i][j] {
				s.Weighings[i][j][k] = p.Index(e)
			}
			sort.Sort(sort.IntSlice(s.Weighings[i][j]))
		}
	}

	for i, _ := range s.Coins {
		s.Coins[i] = i
	}
}
