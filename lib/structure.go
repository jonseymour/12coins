package lib

import (
	"fmt"
	//	"os"
)

type StructureType uint8

type Flips [3][2]int
const (
	P StructureType = iota
	Q
	R
	S
	T
)

const (
	A uint = iota
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
)

type Structure interface {
	Type() StructureType
	String() string

	// Extracts the identity of the coin.
	// s is the current Solution
	// i is the row index in the canonical solution
	// r is the mapping between the canonical solution and s
	// p is the permutation from {0,11} to the current Solution
	Encode(s *Solution, i int, r [3]int, p *[12]int)
	Pan(i int) int
}

func (f Flips) Encode() uint {
	F := uint(0)
	for i, p := range f {
		if p[0] == 0 {
			F |= (1 << uint(i))
		}
	}
	return F
}

// Encodes the structure of a weighing.
type structure struct {
	_type       StructureType
	permutation [2]int
}

func (s *structure) Pan(i int) int {
	return s.permutation[i]
}

type structureP struct {
	structure
}

// Return two sets containing the left and right members of the
// pair which intersects a set of left coins.
func splitPair(pairs [3]CoinSet, left CoinSet) (CoinSet, CoinSet) {
	var ll CoinSet
	var rr CoinSet
	for _, pp := range pairs {
		ll = pp.Intersection(left)
		if ll.Size() == 1 {
			rr = pp.Complement(ll)
			return ll, rr
		}
	}
	panic(fmt.Errorf("illegal state: could not find an expected split pair: %v, %v", pairs, left))
}

// (3T,1L,1U) (1T,2J,1R)
func (sp *structureP) Encode(s *Solution, i int, r [3]int, p *[12]int) {
	left := s.Weighings[r[i]].Pan(sp.permutation[0])
	right := s.Weighings[r[i]].Pan(sp.permutation[1])

	ll, rr := splitPair(s.Pairs, left)

	switch i {
	case 0:
		(*p)[A] = left.Intersection(s.Unique).ExactlyOne(0)
		(*p)[D] = ll.ExactlyOne(0)
		(*p)[L] = right.Intersection(s.Triples).ExactlyOne(0)
		(*p)[E] = rr.ExactlyOne(0)
	case 1:
		(*p)[B] = left.Intersection(s.Unique).ExactlyOne(0)
		(*p)[F] = ll.ExactlyOne(0)
		(*p)[J] = right.Intersection(s.Triples).ExactlyOne(0)
		(*p)[G] = rr.ExactlyOne(0)
	case 2:
		(*p)[C] = left.Intersection(s.Unique).ExactlyOne(0)
		(*p)[H] = ll.ExactlyOne(0)
		(*p)[K] = right.Intersection(s.Triples).ExactlyOne(0)
		(*p)[I] = rr.ExactlyOne(0)
	default:
		panic(fmt.Errorf("illegal argument: i: %d", i))
	}
}

type structureQ struct {
	structure
}

func (sp *structureQ) Encode(s *Solution, i int, r [3]int, p *[12]int) {
	left := s.Weighings[r[i]].Pan(sp.permutation[0])
	right := s.Weighings[r[i]].Pan(sp.permutation[1])

	ll, rr := splitPair(s.Pairs, left)

	(*p)[D] = ll.ExactlyOne(0)
	(*p)[E] = rr.ExactlyOne(0)
	(*p)[A] = right.Intersection(s.Unique).ExactlyOne(0)

	row1right := s.Weighings[r[1]].Pan(s.Structure[r[1]].Pan(1))
	row2right := s.Weighings[r[2]].Pan(s.Structure[r[2]].Pan(1))
	(*p)[L] = left.Intersection(s.Triples).Complement(row1right).Complement(row2right).ExactlyOne(0)
}

type structureR struct {
	structure
}

func (sp *structureR) Encode(s *Solution, i int, r [3]int, p *[12]int) {
	left := s.Weighings[r[1]].Pan(sp.permutation[0])
	right := s.Weighings[r[1]].Pan(sp.permutation[1])
	allPairs := s.Pairs[0].Union(s.Pairs[1]).Union(s.Pairs[2])

	(*p)[J] = right.Intersection(s.Triples).ExactlyOne(0)
	(*p)[B] = right.Intersection(s.Unique).ExactlyOne(0)

	row0 := s.Weighings[r[0]].Both()
	row2 := s.Weighings[r[2]].Both()
	(*p)[F] = left.Intersection(row0).Intersection(allPairs).ExactlyOne(0)
	(*p)[G] = right.Intersection(row0).Intersection(allPairs).ExactlyOne(0)
	(*p)[H] = left.Intersection(row2).Intersection(allPairs).ExactlyOne(0)
	(*p)[I] = right.Intersection(row2).Intersection(allPairs).ExactlyOne(0)
}

type structureS struct {
	structure
}

func (sp *structureS) Encode(s *Solution, i int, r [3]int, p *[12]int) {
	right := s.Weighings[r[i]].Pan(sp.permutation[1])

	(*p)[C] = right.Intersection(s.Unique).ExactlyOne(0)
	(*p)[K] = right.Intersection(s.Triples).ExactlyOne(0)
}

type structureT struct {
	structure
}

func (sp *structureT) Encode(s *Solution, i int, r [3]int, p *[12]int) {
	left := s.Weighings[r[i]].Pan(sp.permutation[0])

	row0right := s.Weighings[r[0]].Pan(s.Structure[r[0]].Pan(1))
	row1right := s.Weighings[r[1]].Pan(s.Structure[r[1]].Pan(1))

	(*p)[C] = left.Intersection(s.Unique).ExactlyOne(0)
	(*p)[K] = left.Intersection(s.Triples).Complement(row0right).Complement(row1right).ExactlyOne(0)
}

func (t StructureType) String() string {
	switch t {
	case P:
		return "p"
	case Q:
		return "q"
	case R:
		return "r"
	case S:
		return "s"
	case T:
		return "t"
	default:
		panic(fmt.Errorf("unhandled case: %d", t))
	}
}

// Returns a new structure of the specified type.
func NewStructure(t StructureType, p [2]int) Structure {
	if p[0] == p[1] {
		panic(fmt.Errorf("illegal argument: p[0] == p[1]: %d=%d", p[0], p[1]))
	}
	s := structure{
		_type:       t,
		permutation: p,
	}
	switch t {
	case P:
		return &structureP{
			structure: s,
		}
	case Q:
		return &structureQ{
			structure: s,
		}
	case R:
		return &structureR{
			structure: s,
		}
	case S:
		return &structureS{
			structure: s,
		}
	case T:
		return &structureT{
			structure: s,
		}
	default:
		panic(fmt.Errorf("illegal argument: t: %v", t))
	}
}

func (s *structure) String() string {
	return fmt.Sprintf("%v[%d,%d]", s._type, s.permutation[0], s.permutation[1])
}

func (s *structure) Type() StructureType {
	return s._type
}

func ParseStructure(r string) (Structure, error) {
	var (
		s  *string
		p0 *int
		p1 *int
	)
	if n, err := fmt.Sscanf(r, "%s[%d,%d]", &s, &p0, &p1); n == 3 && err == nil {
		var t StructureType
		if s == nil || p0 == nil || p1 == nil {
			return nil, fmt.Errorf("scanning failed: %s", r)
		}
		switch *s {
		case "p":
			t = P
		case "q":
			t = Q
		case "r":
			t = R
		case "s":
			t = S
		case "t":
			t = T
		default:
			return nil, fmt.Errorf("failed to parse structure type: %s", s)
		}
		return NewStructure(t, [2]int{*p0, *p1}), nil
	} else {
		return nil, err
	}
}

func (s *Solution) deriveOneStructure(nT uint8, l CoinSet, pi [2]int) (Structure, error) {

	u := l.Intersection(s.Unique)
	switch nT {
	case 3:
		switch u.Size() {
		case 1:
			return NewStructure(T, pi), nil
		case 0:
			return NewStructure(Q, pi), nil
		default:
			return nil, fmt.Errorf("illegal state: t==3, u > 1")
		}
	case 2:
		switch u.Size() {
		case 1:
			return NewStructure(P, pi), nil
		case 0:
			for _, pair := range s.Pairs {
				match := pair.Intersection(l)
				switch match.Size() {
				case 0:
					continue
				case 1:
					return NewStructure(R, pi), nil
				case 2:
					return NewStructure(S, pi), nil
				default:
					return nil, fmt.Errorf("illegal state: t==2 && u==0: wrong size")
				}
				break
			}
			return nil, fmt.Errorf("illegal state: t==2, u == 0: end of loop")
		default:
			return nil, fmt.Errorf("illegal state: t==2, u > 1")
		}
	default:
		return nil, fmt.Errorf("illegal state: t < 2 || t > 3")
	}

}

// Return a clone of the receiver in which the structure has been populated.
func (s *Solution) AnalyseStructure() (*Solution, error) {
	var r *Solution
	var err error

	if s.flags&GROUPED == 0 {
		r, err = s.Groupings()
	} else {
		r = s.Clone()
	}

	if err != nil {
		return r, err
	}

	F := 0

	for i, w := range r.Weighings {
		pi := [2]int{0, 1}
		l := w.Left()
		t := l.Intersection(r.Triples)
		if t.Size() < 2 {
			pi = [2]int{1, 0}
			l = w.Right()
			t = l.Intersection(r.Triples)
			F |= (1 << uint(i))
		}
		if r.Structure[i], err = r.deriveOneStructure(t.Size(), l, pi); err != nil {
			s.markInvalid()
			return s, err
		}
	}

	p := [3]int{0, 1, 2} // a permutation of rows from the canonical form to the current form
	st := [3]StructureType{P, P, P}
	lock := false
	for i, rs := range r.Structure {

		switch rs.Type() {
		case Q:
			p[0] = i
			st[0] = Q
		case R:
			st[1] = R
			p[1] = i
		case S:
			st[2] = S
			p[2] = i
		case P:
			if st[0] == P && !lock {
				// if there is a single P, then we need to move it. but otherwise
				// we leave it in place. required for PRS and PRT cases where P
				// is not in position 0. specifically must not move the second or
				// third P of a PPP or QPP case.
				lock = true
				p[0] = i
			}
		case T:
			st[2] = T
			p[2] = i
		default:
			panic(fmt.Errorf("unknown structure: %v", r.Structure[i]))
		}

	}

	if (st[0] == Q || st[0] == P) && st[1] != R {

		// The P and Q structures are ambiguous until this point.

		if st[1] != P || st[2] != P {
			panic(fmt.Errorf("illegal state: st[1] != P || st[2] != P: %v", st))
		}
		switch p[0] {
		case 0: // PPP or QPP
			p[0] = 0
			p[1] = 1
			p[2] = 2
		case 1: // PQP
			p[0] = 1
			p[1] = 2
			p[2] = 0
		case 2: // PPQ
			p[0] = 2
			p[1] = 0
			p[2] = 1
		default:
			panic(fmt.Errorf("illegal state: p[0] < 0 || p[0] > 2: %d", p[0]))
		}
	}

	// st now contains the canonical structure - one of qrs, prt, prs, qpp or ppp.
	// p now contains the mapping from the canonical structure to the actual structure
	// F now contain the flips required to arrange each weighing in canonical order.

	sS := uint(0)

	switch st[0] {
	case P:
		switch st[2] {
		case T:
			sS = Number(p[0:]) + 12
		case S:
			sS = Number(p[0:]) + 6
		case P:
			sS = 21
		}
	case Q:
		if st[1] == P {
			sS = 18 + Number(p[0:])/2
		} else {
			sS = Number(p[0:])
		}
	default:
		panic(fmt.Errorf("illegal state: st[0] not in (P,Q)"))
	}

	// sS now encodes sT as a single number between 0 and 21

	Pa := [12]int{}
	r.encoding.P = &Pa
	r.encoding.S = pi(int(sS))
	r.encoding.F = pi(F)

	for i, e := range p {
		r.Structure[e].Encode(r, i, p, r.encoding.P)
	}

	// r.encoding.P now contains the permutation to be applied to 1,12

	n := Number(Pa[0:])*176 + 8*sS + uint(F)
	r.encoding.N = &n

	r.flags |= ANALYSED
	return r, nil
}

// Return a clone of the receiver in which the weighings have been permuted into the
// the canonical order and all sets are ordered sets.
func (s *Solution) Canonical() (*Solution, error) {
	var r *Solution
	var err error

	if s.flags&ANALYSED == 0 {
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
