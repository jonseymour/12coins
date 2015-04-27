package lib

import (
	"fmt"
	"os"
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
	Type() StructureType
	String() string
}

// Encodes the structure of a weighing.
type structure struct {
	_type       StructureType
	permutation [2]int
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
	return &structure{
		_type:       t,
		permutation: p,
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

	p := [3]int{0, 1, 2}
	st := [3]StructureType{P, P, P}
	F := 0

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
			F |= (1 << uint(i))
		}
		u := l.Intersection(r.Unique)
		switch t.Size() {
		case 3:
			switch u.Size() {
			case 1:
				r.Structure[i] = NewStructure(T, pi)
			case 0:
				r.Structure[i] = NewStructure(Q, pi)
			default:
				s.markInvalid()
				return s, fmt.Errorf("illegal state: t==3, u > 1")
			}
		case 2:
			switch u.Size() {
			case 1:
				r.Structure[i] = NewStructure(P, pi)
			case 0:
				for _, pair := range r.Pairs {
					match := pair.Intersection(l)
					switch match.Size() {
					case 0:
						continue
					case 1:
						r.Structure[i] = NewStructure(R, pi)
					case 2:
						r.Structure[i] = NewStructure(S, pi)
					default:
						s.markInvalid()
						return s, fmt.Errorf("illegal state: t==2 && u==0 && j==0 && l == 0")
					}
					break
				}
			default:
				s.markInvalid()
				return s, fmt.Errorf("illegal state: t==2, u > 1")
			}
		default:
			s.markInvalid()
			return s, fmt.Errorf("illegal state: t < 2 || t > 3")
		}

		switch r.Structure[i].Type() {
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
			if st[0] == P {
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
		if st[1] != P || st[2] != P {
			panic(fmt.Errorf("illegal state: st[1] != P || st[2] != P: %v", st))
		}
		switch p[0] {
		case 0:
			p[0] = 0
			p[1] = 1
			p[2] = 2
		case 1:
			p[0] = 1
			p[1] = 2
			p[2] = 0
		case 2:
			p[0] = 2
			p[1] = 0
			p[2] = 1
		default:
			panic(fmt.Errorf("illegal state: p[0] < 0 || p[0] > 2: %d", p[0]))
		}
	}

	sS := uint(0)

	switch st[0] {
	case P:
		if st[1] == P {
			sS = 21
		} else {
			switch st[2] {
			case T:
				sS = Number(p[0:]) + 12
			case S:
				sS = Number(p[0:]) + 6
			}
		}
	case Q:
		if st[1] == P {
			sS = 18 + Number(p[0:])/2
		} else {
			sS = Number(p[0:])
		}
	default:
		panic(fmt.Errorf("illegal state: st[1] != P"))
	}

	fmt.Fprintf(os.Stderr, "%d %v\n", sS, p)

	r.encoding.S = pi(int(sS))
	r.encoding.F = pi(F)

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
