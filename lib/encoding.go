package lib

import (
	"encoding/json"
)

type encoding struct {
	Weighings *[3][2][]int `json:"weighings,omitempty"` // the weighings of the solution
	Unique    *[]int       `json:"unique,omitempty"`    // the coins that appear in one weighing
	Pairs     *[3][2]int   `json:"pairs,omitempty"`     // the pairs that appear in exactly two weighings
	Triples   *[]int       `json:"triples,omitempty"`   // the coins that appear in all 3 weighings
	Structure *[3]string   `json:"structure,omitempty"` // An encoding of the structure
	ZeroCoin  *int         `json:"zero-coin,omitempty"` // the zero coin of the weighings. either 0 or 1.
}

// Convert the solution to its JSON representation.
func (s *Solution) String() string {
	s.Encode()
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Solution) Encode() {
	z := s.GetZeroCoin()
	tmp := [3][2][]int{}
	s.encoding.ZeroCoin = pi(z)
	if *s.encoding.ZeroCoin == 1 {
		s.encoding.ZeroCoin = nil
	}
	s.encoding.Weighings = &tmp
	for i, w := range s.Weighings {
		for j, p := range w.Pans() {
			s.encoding.Weighings[i][j] = p.AsCoins(z)
		}
	}
	if s.Unique != nil {
		tmp := s.Unique.AsCoins(z)
		s.encoding.Unique = &tmp
	}
	if s.Triples != nil {
		tmp := s.Triples.AsCoins(z)
		s.encoding.Triples = &tmp
	}
	if s.Unique != nil {
		tmp := [3][2]int{}
		for i, _ := range tmp {
			if s.Pairs[i] != nil {
				copy(tmp[i][0:], s.Pairs[i].AsCoins(z))
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

func (s *Solution) Decode() {
	z := s.GetZeroCoin()
	if s.encoding.Weighings != nil {
		for i, w := range *s.encoding.Weighings {
			s.Weighings[i] = NewWeighing(NewCoinSet(w[0], z), NewCoinSet(w[1], z))
		}
	}
	if s.encoding.Unique != nil {
		s.Unique = NewCoinSet(*s.encoding.Unique, z)
	}
	if s.encoding.Triples != nil {
		s.Triples = NewCoinSet(*s.encoding.Triples, z)
	}
	if s.encoding.Pairs != nil {
		for i, p := range *s.encoding.Pairs {
			s.Pairs[i] = NewCoinSet(p[0:], z)
		}
	}
	if s.encoding.Structure != nil {
		for i, t := range *s.encoding.Structure {
			s.Structure[i], _ = ParseStructure(t, i)
		}
	}
}
