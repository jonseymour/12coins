package lib

import (
	"fmt"
)

// Tabulate the singletons, pairs and triples of the solution.
func (s *Solution) Groupings() (*Solution, error) {

	var clone *Solution

	if s.flags&REVERSED == 0 {
		clone = s.Clone()
	} else {
		var err error
		if clone, err = s.Reverse(); err != nil {
			return clone, err
		}
	}

	clone.reset()

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
		s.flags = INVALID
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
