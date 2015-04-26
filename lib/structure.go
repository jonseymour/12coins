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
