package lib

import (
	"fmt"
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
	Encode(s *Solution, i int, p []int)
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
	_type StructureType
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
func (sp *structureP) Encode(s *Solution, i int, p []int) {
	left := s.Weighings[i].Left()
	right := s.Weighings[i].Right()

	ll, rr := splitPair(s.Pairs, left)

	switch i {
	case 0:
		p[A] = left.Intersection(s.Unique).ExactlyOne(0)
		p[D] = ll.ExactlyOne(0)
		p[L] = right.Intersection(s.Triples).ExactlyOne(0)
		p[E] = rr.ExactlyOne(0)
	case 1:
		p[B] = left.Intersection(s.Unique).ExactlyOne(0)
		p[F] = ll.ExactlyOne(0)
		p[J] = right.Intersection(s.Triples).ExactlyOne(0)
		p[G] = rr.ExactlyOne(0)
	case 2:
		p[C] = left.Intersection(s.Unique).ExactlyOne(0)
		p[H] = ll.ExactlyOne(0)
		p[K] = right.Intersection(s.Triples).ExactlyOne(0)
		p[I] = rr.ExactlyOne(0)
	default:
		panic(fmt.Errorf("illegal argument: i: %d", i))
	}
}

type structureQ struct {
	structure
}

func (sp *structureQ) Encode(s *Solution, i int, p []int) {
	left := s.Weighings[i].Left()
	right := s.Weighings[i].Right()

	ll, rr := splitPair(s.Pairs, left)

	p[D] = ll.ExactlyOne(0)
	p[E] = rr.ExactlyOne(0)
	p[A] = right.Intersection(s.Unique).ExactlyOne(0)

	row1right := s.Weighings[1].Pan(1)
	row2right := s.Weighings[2].Pan(1)
	p[L] = left.Intersection(s.Triples).Complement(row1right).Complement(row2right).ExactlyOne(0)
}

type structureR struct {
	structure
}

func (sp *structureR) Encode(s *Solution, i int, p []int) {
	left := s.Weighings[i].Left()
	right := s.Weighings[i].Right()
	allPairs := s.Pairs[0].Union(s.Pairs[1]).Union(s.Pairs[2])

	p[J] = right.Intersection(s.Triples).ExactlyOne(0)
	p[B] = right.Intersection(s.Unique).ExactlyOne(0)

	row0 := s.Weighings[0].Both()
	row2 := s.Weighings[2].Both()
	p[F] = left.Intersection(row0).Intersection(allPairs).ExactlyOne(0)
	p[G] = right.Intersection(row0).Intersection(allPairs).ExactlyOne(0)
	p[H] = left.Intersection(row2).Intersection(allPairs).ExactlyOne(0)
	p[I] = right.Intersection(row2).Intersection(allPairs).ExactlyOne(0)
}

type structureS struct {
	structure
}

func (sp *structureS) Encode(s *Solution, i int, p []int) {
	right := s.Weighings[i].Right()

	p[C] = right.Intersection(s.Unique).ExactlyOne(0)
	p[K] = right.Intersection(s.Triples).ExactlyOne(0)
}

type structureT struct {
	structure
}

func (sp *structureT) Encode(s *Solution, i int, p []int) {
	left := s.Weighings[i].Pan(0)

	row0right := s.Weighings[0].Right()
	row1right := s.Weighings[1].Right()

	p[C] = left.Intersection(s.Unique).ExactlyOne(0)
	p[K] = left.Intersection(s.Triples).Complement(row0right).Complement(row1right).ExactlyOne(0)
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
func NewStructure(t StructureType) Structure {
	s := structure{
		_type: t,
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
	return fmt.Sprintf("%v", s._type)
}

func (s *structure) Type() StructureType {
	return s._type
}

func ParseStructure(r string) (Structure, error) {
	var (
		s *string
	)
	var t StructureType
	switch r {
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
	return NewStructure(t), nil
}

func (s *Solution) deriveOneStructure(nT uint8, l CoinSet) (Structure, error) {

	u := l.Intersection(s.Unique)
	switch nT {
	case 3:
		switch u.Size() {
		case 1:
			return NewStructure(T), nil
		case 0:
			return NewStructure(Q), nil
		default:
			return nil, fmt.Errorf("illegal state: t==3, u > 1")
		}
	case 2:
		switch u.Size() {
		case 1:
			return NewStructure(P), nil
		case 0:
			for _, pair := range s.Pairs {
				match := pair.Intersection(l)
				switch match.Size() {
				case 0:
					continue
				case 1:
					return NewStructure(R), nil
				case 2:
					return NewStructure(S), nil
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

// Derive the structure of the 3 weighings. Return a number F which encodes how the weighings
// must be flipped in order to obtain canonical form.
func (s *Solution) deriveStructure() (Flips, error) {
	var flips Flips
	var err error
	for i, w := range s.Weighings {
		l := w.Left()
		t := l.Intersection(s.Triples)
		if t.Size() < 2 {
			flips[i] = [2]int{1, 0}
			l = w.Right()
			t = l.Intersection(s.Triples)
		} else {
			flips[i] = [2]int{0, 1}
		}
		if s.Structure[i], err = s.deriveOneStructure(t.Size(), l); err != nil {
			return flips, err
		}
	}
	return flips, nil
}

// Derive the permutation that maps the canonical weighing order to the receiver's order.
func (s *Solution) deriveCanonicalOrder() ([3]int, [3]StructureType, error) {
	np := 0
	p := [3]int{0, 1, 2} // a permutation of rows from the canonical form to the current form
	st := [3]StructureType{P, P, P}
	for i, rs := range s.Structure {

		switch rs.Type() {
		case Q:
			st[0] = Q
			p[0] = i
		case R:
			st[1] = R
			p[1] = i
		case S, T:
			st[2] = rs.Type()
			p[2] = i
		case P:
			if st[0] == P && np == 0 {
				// if there is a single P, then we need to move it. but otherwise
				// we leave it in place. required for PRS and PRT cases where P
				// is not in position 0. specifically must not move the second or
				// third P of a PPP or QPP case.
				p[0] = i
			}
			np += 1
		default:
			return p, st, fmt.Errorf("unknown structure: %v", rs)
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
			return p, st, fmt.Errorf("illegal state: p[0] < 0 || p[0] > 2: %d", p[0])
		}
	}
	return p, st, nil
}

func (s *Solution) deriveCanonical() (*Solution, error) {
	var r *Solution
	var err error
	if s.flags&ANALYSED == 0 {
		if s, err = s.AnalyseStructure(); err != nil {
			return s, err
		}
	}

	r = s.Clone()

	r.reset()
	r.Triples = s.Triples
	r.Unique = s.Unique
	r.Pairs = s.Pairs
	st := [3]StructureType{}
	for i, _ := range r.Weighings {
		si := s.order[i]
		sw := s.Weighings[si]
		sf := s.flips[si]
		st[i] = s.Structure[si].Type()
		r.Weighings[i] = NewWeighing(sw.Pan(sf[0]), sw.Pan(sf[1]))
		r.Structure[i] = NewStructure(st[i])
	}
	r.order = [3]int{0, 1, 2}
	r.flips = [3][2]int{}
	r.Coins = []int{}
	r.Weights = []Weight{}
	p := make([]int, 12)
	for i, e := range r.Structure {
		e.Encode(r, i, p)
	}

	S := EncodeStructure(r.order, st)
	N := Number(p[0:])*176 + S

	r.Unique = NewOrderedCoinSet(p[0:3], 0)
	r.Pairs[0] = NewOrderedCoinSet(p[3:5], 0)
	r.Pairs[1] = NewOrderedCoinSet(p[5:7], 0)
	r.Pairs[2] = NewOrderedCoinSet(p[7:9], 0)
	r.Triples = NewOrderedCoinSet(p[9:12], 0)

	r.encoding.P = p
	r.encoding.F = pu(0)
	r.encoding.S = &S
	r.encoding.N = &N

	r.flags |= (CANONICALISED | NUMBERED | GROUPED) &^ REVERSED
	return r, nil
}

// Encode p and s as a number between 0 and 21.
func EncodeStructure(p [3]int, st [3]StructureType) uint {
	s := uint(0)

	switch st[0] {
	case P:
		switch st[2] {
		case T:
			s = Number(p[0:]) + 16
		case S:
			s = Number(p[0:]) + 10
		case P:
			s = 0
		}
	case Q:
		if st[1] == P {
			s = 1 + Number(p[0:])/2
		} else {
			s = Number(p[0:]) + 4
		}
	default:
		panic(fmt.Errorf("illegal state: st[0] not in (P,Q)"))
	}
	return s
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

	if flips, err := r.deriveStructure(); err != nil {
		s.markInvalid()
		return s, err
	} else if o, st, err := r.deriveCanonicalOrder(); err != nil {
		s.markInvalid()
		return s, err
	} else {

		// st now contains the canonical structure - one of qrs, prt, prs, qpp or ppp.
		// p now contains the mapping from the canonical structure to the actual structure
		// F now contain the flips required to arrange each weighing in canonical order.

		// sS now encodes sT as a single number between 0 and 21

		F := flips.Encode()
		S := EncodeStructure(o, st)

		r.order = o
		r.flips = flips
		r.flags |= ANALYSED

		var canonical *Solution
		if canonical, err = r.deriveCanonical(); err != nil {
			return s, err
		}

		N := (*canonical.encoding.N/176)*176 + 22*F + S

		r.Triples = canonical.Triples
		r.Unique = canonical.Unique
		r.Pairs = canonical.Pairs

		r.encoding.P = canonical.P
		r.encoding.S = &S
		r.encoding.F = &F
		r.encoding.N = &N

		// r.encoding.P now contains the permutation to be applied to 1,12

		r.flags |= NUMBERED
		return r, nil
	}
}

// Return a clone of the receiver in which the weighings have been permuted into the
// the canonical order and all sets are ordered sets.
func (s *Solution) Canonical() (*Solution, error) {
	var r *Solution
	var err error

	if s.flags&CANONICALISED == 0 {
		r, err = s.deriveCanonical()
	} else {
		r = s.Clone()
	}

	return r, err
}
