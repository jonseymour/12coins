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
	Permutation []int       `json:"permutation,omitempty"`
	Weighings   [3][2][]int `json:"weighings"`
	Coins       [12]int     `json:"coins"`
	Weights     [12]Weight  `json:"weights"`
	Mirror      bool        `json:"mirror,omitempty"`
	ZeroCoin    int         `json:"zero-coin,omitempty"`
}

func (s *Solver) Decide(scale Scale) (int, Weight) {
	scale.SetZeroCoin(s.ZeroCoin)
	a := scale.Weigh(s.Weighings[0][0], s.Weighings[0][1])
	b := scale.Weigh(s.Weighings[1][0], s.Weighings[1][1])
	c := scale.Weigh(s.Weighings[2][0], s.Weighings[2][1])

	i := int(a*9 + b*3 + c - 13)
	o := abs(i)
	if s.Mirror {
		o = 13 - o
	}

	f := s.Coins[o-1]
	w := s.Weights[o-1]

	if i > 0 {
		w = Heavy - w
	}

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
	clone := Solver{
		Weighings: [3][2][]int{},
		Coins:     [12]int{},
		Weights:   [12]Weight{},
		ZeroCoin:  s.ZeroCoin,
		Mirror:    s.Mirror,
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

func (s *Solver) Relabel() {
	c := make([]int, len(s.Coins), len(s.Coins))
	for i, e := range s.Coins {
		c[i] = e
	}
	p := NewPermutation(c, s.ZeroCoin)

	for i, _ := range s.Weighings {
		for j, _ := range []int{0, 1} {
			for k, e := range s.Weighings[i][j] {
				s.Weighings[i][j][k] = p.Index(e) + s.ZeroCoin
			}
			sort.Sort(sort.IntSlice(s.Weighings[i][j]))
		}
	}

	for i, _ := range s.Coins {
		s.Coins[i] = i + s.ZeroCoin
	}
}

func (s *Solver) Normalize() {
	for i, _ := range s.Weighings {
		for j, _ := range []int{0, 1} {
			sort.Sort(sort.IntSlice(s.Weighings[i][j]))
		}
	}
}

func (s *Solver) Reverse() error {
	seen := [12]bool{}
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
		o := NewOracle(i, Light, 1)
		ri, rw := s.Decide(o)
		if ri != i {
			if seen[ri-1] {
				return fmt.Errorf("cannot distinguish between (%d, %v) and (%d, %v) ", s.Coins[ri-1], s.Weights[ri-1], i, rw)
			} else {
				seen[ri-1] = true
			}
			s.Coins[ri-1] = i
		}
		if rw != Light {
			s.Weights[ri-1] = Heavy - s.Weights[ri-1]
		}
	}
	return nil
}
