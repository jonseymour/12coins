package lib

import (
	"fmt"
)

type StructureType uint8

const (
	P StructureType = iota
	Q
	R
	S
	T
)

type Structure interface {
	String() string
}

// Encodes the structure of a weighing.
type structure struct {
	_type       StructureType
	permutation [2]int
	index       int
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
func NewStructure(t StructureType, p [2]int, i int) Structure {
	if p[0] == p[1] {
		panic(fmt.Errorf("illegal argument: p[0] == p[1]: %d=%d", p[0], p[1]))
	}
	if i < 0 || i > 2 {
		panic(fmt.Errorf("illegal argument: i < 0 || i > 2: %d", i))
	}
	return &structure{
		_type:       t,
		permutation: p,
		index:       i,
	}
}

func (s *structure) String() string {
	return fmt.Sprintf("%v[%d, %d]", s._type, s.permutation[0], s.permutation[1])
}

func ParseStructure(r string, i int) (Structure, error) {
	var (
		s  *string
		p0 *int
		p1 *int
	)
	if n, err := fmt.Sscanf(r, "%s[%d, %d]", &s, &p0, &p1); n == 3 && err == nil {
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
		return NewStructure(t, [2]int{*p0, *p1}, i), nil
	} else {
		return nil, err
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

	for i, w := range r.Weighings {
		pi := [2]int{0, 1}
		tl := w.Left().Intersection(r.Triples)
		tr := w.Right().Intersection(r.Triples)
		t := tl
		l := w.Left()
		if tr.Size() > tl.Size() {
			pi = [2]int{1, 0}
			t = tr
			l = w.Right()
		}
		u := l.Intersection(r.Unique)
		switch t.Size() {
		case 3:
			switch u.Size() {
			case 1:
				r.Structure[i] = NewStructure(T, pi, i)
			case 0:
				r.Structure[i] = NewStructure(Q, pi, i)
			default:
				s.flags = INVALID
				return s, fmt.Errorf("illegal state: t==3, u > 1")
			}
		case 2:
			switch u.Size() {
			case 1:
				r.Structure[i] = NewStructure(P, pi, i)
			case 0:
				for _, pair := range r.Pairs {
					match := pair.Intersection(l)
					switch match.Size() {
					case 0:
						continue
					case 1:
						r.Structure[i] = NewStructure(R, pi, i)
					case 2:
						r.Structure[i] = NewStructure(S, pi, i)
					default:
						s.flags = INVALID
						return s, fmt.Errorf("illegal state: t==2 && u==0 && j==0 && l == 0")
					}
					break
				}
			default:
				s.flags = INVALID
				return s, fmt.Errorf("illegal state: t==2, u > 1")
			}
		default:
			s.flags = INVALID
			return s, fmt.Errorf("illegal state: t < 2 || t > 3")
		}
	}

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
